/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainedbft

import (
	"bytes"

	"errors"
	"fmt"

	"chainmaker.org/chainmaker-go/common/msgbus"
	timeservice "chainmaker.org/chainmaker-go/consensus/chainedbft/time_service"
	"chainmaker.org/chainmaker-go/consensus/chainedbft/utils"
	"chainmaker.org/chainmaker-go/consensus/governance"
	"chainmaker.org/chainmaker-go/pb/protogo/common"
	"chainmaker.org/chainmaker-go/pb/protogo/consensus"
	chainedbftpb "chainmaker.org/chainmaker-go/pb/protogo/consensus/chainedbft"
	"chainmaker.org/chainmaker-go/protocol"
	chainUtils "chainmaker.org/chainmaker-go/utils"
	"github.com/gogo/protobuf/proto"
)

var (
	InvalidPeerErr        = errors.New("invalid peer")
	ValidateSignErr       = errors.New("validate sign error")
	VerifySignerFailedErr = errors.New("verify signer failed")
)

//processNewHeight If the local node is one of the validators in current epoch, update SMR state to ConsStateType_NewLevel
//and prepare to generate a new block if local node is proposer in the current level
func (cbi *ConsensusChainedBftImpl) processNewHeight(height uint64, level uint64) {
	if cbi.smr.getHeight() != height || cbi.smr.getCurrentLevel() > level ||
		(cbi.smr.getCurrentLevel() == level && cbi.smr.state != chainedbftpb.ConsStateType_NewHeight) {
		cbi.logger.Debugf("service selfIndexInEpoch [%v] processNewHeight: invalid input[%v:%v], smr height [%v] level [%v] state %v",
			cbi.selfIndexInEpoch, height, level, cbi.smr.getHeight(), cbi.smr.getCurrentLevel(), cbi.smr.state)
		return
	}

	cbi.logger.Debugf("service selfIndexInEpoch [%v] processNewHeight at height [%v] level [%v], smr height [%v] level [%v] state %v epoch %v",
		cbi.selfIndexInEpoch, height, level, cbi.smr.getHeight(), cbi.smr.getCurrentLevel(), cbi.smr.state, cbi.smr.getEpochId())
	begin := chainUtils.CurrentTimeMillisSeconds()
	if !cbi.smr.isValidIdx(cbi.selfIndexInEpoch) {
		cbi.logger.Infof("self selfIndexInEpoch [%v] is not in current consensus epoch", cbi.selfIndexInEpoch)
		return
	}

	cbi.smr.updateState(chainedbftpb.ConsStateType_NewLevel)
	cbi.logger.Debugf("service selfIndexInEpoch [%v] processNewRound at height [%v] level [%v], smr height [%v] level [%v] state %v",
		cbi.selfIndexInEpoch, height, level, cbi.smr.getHeight(), cbi.smr.getCurrentLevel(), cbi.smr.state)

	cbi.processNewLevel(height, level)
	cbi.logger.Debugf("processNewHeight total used time: %d", chainUtils.CurrentTimeMillisSeconds()-begin)
}

//processNewLevel update state to ConsStateType_Propose and prepare to generate a new block if local node is proposer in the current level
func (cbi *ConsensusChainedBftImpl) processNewLevel(height uint64, level uint64) {
	if cbi.smr.getHeight() != height || cbi.smr.getCurrentLevel() > level ||
		(cbi.smr.getCurrentLevel() == level && cbi.smr.state >= chainedbftpb.ConsStateType_Propose) {
		cbi.logger.Debugf("service selfIndexInEpoch [%v] processNewLevel: invalid input [%v:%v], smr height [%v] level [%v] state %v",
			cbi.selfIndexInEpoch, height, level, cbi.smr.getHeight(), cbi.smr.getCurrentLevel(), cbi.smr.state)
		return
	}

	cbi.logger.Debugf("service selfIndexInEpoch [%v] processNewLevel at height [%v] level [%v], smr height [%v] level [%v] state %v",
		cbi.selfIndexInEpoch, height, level, cbi.smr.getHeight(), cbi.smr.getCurrentLevel(), cbi.smr.state)
	start := chainUtils.CurrentTimeMillisSeconds()
	hqcBlock := cbi.chainStore.getCurrentCertifiedBlock()
	hqcLevel, err := utils.GetLevelFromBlock(hqcBlock)
	if err != nil {
		cbi.logger.Errorf("get level from block failed, error %v, height [%v]", err, hqcBlock.Header.BlockHeight)
		return
	}
	if hqcLevel >= level {
		cbi.logger.Errorf("given level [%v] too low than certified %v", level, hqcLevel)
		return
	}

	nextProposerIndex := cbi.getProposer(level)
	if cbi.isValidProposer(level, cbi.selfIndexInEpoch) {
		event := &timeservice.TimerEvent{
			Level:    level,
			Height:   height,
			State:    cbi.smr.state,
			Index:    cbi.selfIndexInEpoch,
			EpochId:  cbi.smr.getEpochId(),
			Duration: timeservice.GetEventTimeout(timeservice.PROPOSAL_BLOCK_TIMEOUT, 0),
		}
		cbi.startTimer(event)
		cbi.logger.Infof("service selfIndexInEpoch [%v], build proposal, height: [%v], level [%v]", cbi.selfIndexInEpoch, height, level)
		go cbi.msgbus.Publish(msgbus.BuildProposal, &chainedbftpb.BuildProposal{
			Height:     height,
			IsProposer: true,
			PreHash:    hqcBlock.Header.BlockHash,
		})
	}

	cbi.logger.Infof("service selfIndexInEpoch [%v], waiting proposal, "+
		"height: [%v], level [%v], nextProposerIndex [%d]", cbi.selfIndexInEpoch, height, level, nextProposerIndex)
	cbi.smr.updateState(chainedbftpb.ConsStateType_Propose)
	cbi.logger.Debugf("processNewLevel total used time: %d", chainUtils.CurrentTimeMillisSeconds()-start)
}

//processProposedBlock receive proposed block form core module, then go to new level
func (cbi *ConsensusChainedBftImpl) processProposedBlock(block *common.Block) {
	height := cbi.smr.getHeight()
	level := cbi.smr.getCurrentLevel()
	cbi.logger.Debugf(`processProposedBlock start, block height: [%v], level: [%v]`, block.Header.BlockHeight, level)

	if !cbi.isValidProposer(level, cbi.selfIndexInEpoch) {
		return
	}
	cbi.logger.Infof("receive proposed block by self at height: %d, level: %d, proposer: %s", block.Header.BlockHeight, level, block.Header.Proposer)
	if int64(height) != block.Header.BlockHeight {
		cbi.logger.Errorf(`service id [%v] selfIndexInEpoch [%v] recieve proposed block height [%v]
		 not equal to smr.height [%v]`, cbi.id, cbi.selfIndexInEpoch, block.Header.BlockHeight,
			height)
		return
	}

	proposal := cbi.constructProposal(block, height, level, cbi.smr.getEpochId())
	cbi.signAndBroadcast(proposal)
}

func (cbi *ConsensusChainedBftImpl) processLocalTimeout(height uint64, level uint64) {
	if !cbi.smr.processLocalTimeout(level) {
		return
	}
	var vote *chainedbftpb.ConsensusPayload
	if lastVotedLevel, lastVote := cbi.smr.getLastVote(); lastVotedLevel == level {
		// retry send last vote
		vote = cbi.retryVote(lastVote)
	} else {
		vote = cbi.constructVote(height, level, cbi.smr.getEpochId(), nil)
	}
	cbi.logger.Debugf("service selfIndexInEpoch [%v] processLocalTimeout: broadcasts timeout "+
		"vote [%v:%v] to other validators", cbi.selfIndexInEpoch, height, level)
	cbi.smr.setLastVote(vote, level)
	cbi.signAndBroadcast(vote)
}

func (cbi *ConsensusChainedBftImpl) retryVote(lastVote *chainedbftpb.ConsensusPayload) *chainedbftpb.ConsensusPayload {
	var (
		err      error
		data     []byte
		sign     []byte
		voteMsg  = lastVote.GetVoteMsg()
		voteData = voteMsg.VoteData
	)
	cbi.logger.Debugf("service index [%v] processLocalTimeout: "+
		"get last vote [%v:%v] to other validators, blockId [%x]", cbi.selfIndexInEpoch, voteData.Height, voteData.Level, voteData.BlockID)
	tempVoteData := &chainedbftpb.VoteData{
		NewView:   true,
		Level:     voteData.Level,
		Author:    voteData.Author,
		Height:    voteData.Height,
		BlockID:   voteData.BlockID,
		EpochId:   cbi.smr.getEpochId(),
		AuthorIdx: voteData.AuthorIdx,
	}
	if data, err = proto.Marshal(tempVoteData); err != nil {
		cbi.logger.Errorf("marshal vote failed: %s", err)
		return nil
	}
	if sign, err = cbi.singer.Sign(cbi.chainConf.ChainConfig().Crypto.Hash, data); err != nil {
		cbi.logger.Errorf("failed to sign data failed, err %v data %v", err, data)
		return nil
	}
	serializeMember, err := cbi.singer.GetSerializedMember(true)
	if err != nil {
		cbi.logger.Errorf("failed to get signer serializeMember failed, err %v", err)
		return nil
	}

	tempVoteData.Signature = &common.EndorsementEntry{
		Signer:    serializeMember,
		Signature: sign,
	}
	return &chainedbftpb.ConsensusPayload{
		Type: chainedbftpb.MessageType_VoteMessage,
		Data: &chainedbftpb.ConsensusPayload_VoteMsg{&chainedbftpb.VoteMsg{
			VoteData: tempVoteData,
			SyncInfo: voteMsg.SyncInfo,
		}},
	}
}

func (cbi *ConsensusChainedBftImpl) verifyJustifyQC(qc *chainedbftpb.QuorumCert) error {
	if !qc.NewView && qc.BlockID == nil {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] validate qc failed, nil block id", cbi.selfIndexInEpoch)
		return fmt.Errorf(fmt.Sprintf("nil block id in qc"))
	}

	if qc.NewView && qc.BlockID != nil {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] validate qc failed, invalid block id", cbi.selfIndexInEpoch)
		return fmt.Errorf(fmt.Sprintf("invalid block id in qc"))
	}
	if cbi.smr.getEpochId() == qc.EpochId+1 {
		return nil
	}
	if qc.EpochId != cbi.smr.getEpochId() && (cbi.smr.getEpochId() != qc.EpochId+1) {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] validate qc failed, invalid "+
			"epoch id [%v],need [%v]", cbi.selfIndexInEpoch, qc.EpochId, cbi.smr.getEpochId())
		return fmt.Errorf(fmt.Sprintf("invalid epoch id in qc"))
	}

	newViewNum, votedBlockNum, err := cbi.countNumFromVotes(qc)
	if err != nil {
		return err
	}

	if qc.Level > 0 && qc.NewView && newViewNum < cbi.smr.min() {
		return fmt.Errorf(fmt.Sprintf("vote new view num [%v] less than expected [%v]",
			newViewNum, cbi.smr.min()))
	}
	if qc.Level > 0 && !qc.NewView && votedBlockNum < cbi.smr.min() {
		return fmt.Errorf(fmt.Sprintf("vote block num [%v] less than expected [%v]",
			votedBlockNum, cbi.smr.min()))
	}
	return nil
}

func (cbi *ConsensusChainedBftImpl) needFetch(syncInfo *chainedbftpb.SyncInfo) (bool, error) {
	var (
		err       error
		rootLevel uint64
		qc        = syncInfo.HighestQC
	)
	if rootLevel, err = cbi.chainStore.getRootLevel(); err != nil {
		return false, fmt.Errorf("get root level fail")
	}
	if qc.Level < rootLevel {
		cbi.logger.Debugf("service selfIndexInEpoch [%v] needFetch: syncinfo has an older qc [%v:%v] than root level [%v]",
			cbi.selfIndexInEpoch, qc.Height, qc.Level, rootLevel)
		return false, fmt.Errorf("sync info has a highest quorum certificate with level older than root level")
	}
	if exist, _ := cbi.chainStore.getQC(string(qc.BlockID), qc.Height); exist != nil {
		cbi.logger.Debugf("service selfIndexInEpoch [%v] needFetch: local already has a qc [%v:%v %x]",
			cbi.selfIndexInEpoch, qc.Height, qc.Level, qc.BlockID)
		return false, nil
	}
	if exist, _ := cbi.chainStore.getBlock(string(qc.BlockID), qc.Height); exist != nil {
		cbi.logger.Debugf("service selfIndexInEpoch [%v] needFetch: local already has a qc block [%v:%v %x]",
			cbi.selfIndexInEpoch, qc.Height, qc.Level, qc.BlockID)
		return false, nil
	}

	var (
		currentQC      = cbi.chainStore.getCurrentQC()
		currentTCLevel = cbi.smr.getHighestTCLevel()
		level          = currentQC.Level
	)
	if currentQC.Level < currentTCLevel {
		level = currentTCLevel
	}
	if qc.Level >= level+3 {
		cbi.logger.Debugf("service selfIndexInEpoch [%v] needFetch: local received sync info from a future level [%v], local qc [%v] tc [%v]",
			cbi.selfIndexInEpoch, qc.Level, currentQC.Level, currentTCLevel)
		return false, fmt.Errorf("received sync info from a future level")
	}
	return true, nil
}

func (cbi *ConsensusChainedBftImpl) validateProposalMsg(msg *chainedbftpb.ConsensusMsg) error {
	proposal := msg.Payload.GetProposalMsg().ProposalData
	if proposal.Level < cbi.smr.getCurrentLevel() {
		return fmt.Errorf("old proposal, ignore it")
	}
	if proposal.EpochId != cbi.smr.getEpochId() {
		return fmt.Errorf("err epochId, ignore it")
	}
	if !cbi.validateProposer(msg) {
		return fmt.Errorf("invalid proposer")
	}

	if err := cbi.verifyJustifyQC(proposal.JustifyQC); err != nil {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] validateProposal: block [%v:%v] "+
			"verifyJustifyQC failed, err %v", cbi.selfIndexInEpoch, proposal.Height, proposal.Level, err)
		return fmt.Errorf("failed to verify JustifyQC")
	}
	if !bytes.Equal(proposal.JustifyQC.BlockID, proposal.Block.GetHeader().PreBlockHash) {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] validateProposal: mismatch pre hash [%x] in block, justifyQC %x",
			cbi.selfIndexInEpoch, proposal.Block.GetHeader().PreBlockHash, proposal.JustifyQC.BlockID)
		return fmt.Errorf("mismatch pre hash in block header and justify qc")
	}
	return nil
}

func (cbi *ConsensusChainedBftImpl) validateProposer(msg *chainedbftpb.ConsensusMsg) bool {
	proposal := msg.Payload.GetProposalMsg().ProposalData
	if !cbi.smr.isValidIdx(proposal.ProposerIdx) || !cbi.isValidProposer(proposal.Level, proposal.ProposerIdx) {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] validateProposal: received a proposal "+
			"at height [%v] level [%v] from invalid selfIndexInEpoch [%v] addr [%v]",
			cbi.selfIndexInEpoch, proposal.Height, proposal.Level, proposal.ProposerIdx, proposal.Proposer)
		return false
	}
	if err := cbi.validateSignerAndSignature(msg, cbi.smr.getPeerByIndex(proposal.ProposerIdx)); err != nil {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] validateProposer failed, err %v"+
			" proposal %v, err %v", cbi.selfIndexInEpoch, proposal, err)
		return false
	}
	return true
}

func (cbi *ConsensusChainedBftImpl) processProposal(msg *chainedbftpb.ConsensusMsg) {
	var (
		err         error
		isFetch     bool
		proposalMsg = msg.Payload.GetProposalMsg()
		proposal    = proposalMsg.ProposalData
	)

	cbi.logger.Infof("service selfIndexInEpoch [%v] processProposal step0. proposal.ProposerIdx [%v] ,proposal.Height[%v],"+
		" proposal.Level[%v],proposal.EpochId [%v],expected [%v:%v:%v]", cbi.selfIndexInEpoch, proposal.ProposerIdx, proposal.Height,
		proposal.Level, proposal.EpochId, cbi.smr.getHeight()+1, cbi.smr.getCurrentLevel()+1, cbi.smr.getEpochId())
	validateMsgStart := chainUtils.CurrentTimeMillisSeconds()
	//step0: validate proposal
	if err := cbi.validateProposalMsg(msg); err != nil {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] onReceivedProposal validate proposal failed, err %v",
			cbi.selfIndexInEpoch, err)
		return
	}
	validateMsgEnd := chainUtils.CurrentTimeMillisSeconds()
	usedValidateMsgTime := validateMsgEnd - validateMsgStart
	cbi.logger.Debugf("validate proposal msg success [%d:%d:%d]", proposal.ProposerIdx, proposal.Height, proposal.Level)

	if isFetch, err = cbi.needFetch(proposalMsg.SyncInfo); err != nil {
		cbi.logger.Errorf("needFetch err %v", err)
		return
	} else if isFetch {
		cbi.fetchData(proposal)
	}
	needFetchEnd := chainUtils.CurrentTimeMillisSeconds()
	usedNeedFetchTime := needFetchEnd - validateMsgEnd

	//step1: validate and process new qc from proposal
	cbi.logger.Debugf("service selfIndexInEpoch [%v] processProposal step1 process qc start", cbi.selfIndexInEpoch)
	ok, timeUsed := cbi.processQC(msg)
	if !ok {
		return
	}
	processQCEnd := chainUtils.CurrentTimeMillisSeconds()
	usedProcessQCTime := processQCEnd - needFetchEnd

	//step2: validate new block from proposal
	cbi.logger.Debugf("service selfIndexInEpoch [%v] processProposal step2 validate new block from proposal start", cbi.selfIndexInEpoch)
	if ok := cbi.validateBlock(proposal); !ok {
		return
	}
	validateProposalBlockEnd := chainUtils.CurrentTimeMillisSeconds()
	usedValidateProposalTime := validateProposalBlockEnd - processQCEnd

	//step3: validate consensus args
	cbi.logger.Debugf("service selfIndexInEpoch [%v] processProposal step3 validate consensus args", cbi.selfIndexInEpoch)
	cbi.validateConsensusArg(proposal)
	validateConsensusArgEnd := chainUtils.CurrentTimeMillisSeconds()
	usedValidateConsensusArgTime := validateProposalBlockEnd - validateProposalBlockEnd

	//step4: add proposal to msg pool and add block to chainStore
	cbi.logger.Debugf("service selfIndexInEpoch [%v] processProposal step4 add proposal to msg pool and "+
		"add proposal block to chainStore start", cbi.selfIndexInEpoch)
	if inserted, err := cbi.msgPool.InsertProposal(proposal.Height, proposal.Level, msg); err != nil || !inserted {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] processProposal insert "+
			"a proposal failed, err %v, insert %v", cbi.selfIndexInEpoch, err, inserted)
		return
	}
	insertProposalMsgEnd := chainUtils.CurrentTimeMillisSeconds()
	usedInsertProposalMsgTime := insertProposalMsgEnd - validateConsensusArgEnd

	if executorErr := cbi.chainStore.insertBlock(proposal.GetBlock()); executorErr != nil {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] processProposal add proposal block %v to chainStore failed, err: %s",
			cbi.selfIndexInEpoch, proposal.GetBlock().GetHeader().BlockHeight, executorErr)
		return
	}
	insertProposalBlockEnd := chainUtils.CurrentTimeMillisSeconds()
	usedInsertProposalBlockTime := insertProposalBlockEnd - insertProposalMsgEnd
	//step5: vote it and send vote to next proposer in the epoch
	cbi.generateVoteAndSend(proposal)
	voteProposalEnd := chainUtils.CurrentTimeMillisSeconds()
	usedVoteProposalTime := voteProposalEnd - insertProposalBlockEnd
	totalProcessProposalTime := voteProposalEnd - validateMsgStart
	cbi.logger.Debugf("time costs in processProposal: validateMsgTime: %d, needFetchTime: %d, processQCTime: %d details:[%v], validateProposalBlockTime: %d,"+
		" validateBlockConsensusTime: %d, insertProposalMsgTime: %d, insertProposalBlockTime: %d, voteProposalTime: %d, totalUsedTime: %d ",
		usedValidateMsgTime, usedNeedFetchTime, usedProcessQCTime, timeUsed, usedValidateProposalTime, usedValidateConsensusArgTime,
		usedInsertProposalMsgTime, usedInsertProposalBlockTime, usedVoteProposalTime, totalProcessProposalTime)
}

func (cbi *ConsensusChainedBftImpl) generateVoteAndSend(proposal *chainedbftpb.ProposalData) bool {
	cbi.smr.updateState(chainedbftpb.ConsStateType_Vote)
	cbi.logger.Debugf("service selfIndexInEpoch [%v] processProposal step5 construct vote and "+
		"send vote to next proposer start", cbi.selfIndexInEpoch)
	level, err := utils.GetLevelFromBlock(proposal.GetBlock())
	if err != nil {
		cbi.logger.Errorf("GetLevelFromBlock failed, reason: %s, block: %d:%x", err, proposal.Height, proposal.GetBlock().Header.BlockHash)
		return false
	}
	vote := cbi.constructVote(uint64(proposal.GetBlock().GetHeader().GetBlockHeight()), level, cbi.smr.getEpochId(), proposal.GetBlock())
	cbi.logger.Debugf("service selfIndexInEpoch [%v] processProposal step6 send vote msg to next leader start", cbi.selfIndexInEpoch)
	cbi.sendVote2Next(proposal, vote)
	return true
}

func (cbi *ConsensusChainedBftImpl) fetchData(proposal *chainedbftpb.ProposalData) {
	cbi.logger.Infof("service selfIndexInEpoch [%v] validateProposal need sync up to [%v:%v]",
		cbi.selfIndexInEpoch, proposal.JustifyQC.Height, proposal.JustifyQC.Level)

	//fetch block and qc from proposer
	req := &blockSyncReq{
		targetPeer:  proposal.ProposerIdx,
		blockID:     proposal.JustifyQC.BlockID,
		height:      proposal.JustifyQC.Height,
		startLevel:  cbi.chainStore.getCurrentQC().Level + 1,
		targetLevel: proposal.JustifyQC.Level,
	}

	//note: WaitGroup is used here to provide the blocking function, waiting for SyncManager to synchronize the data,
	//and when SyncManager does not synchronize to the block in the timeout, WaitGroup.Done is called to resolve the blocking
	cbi.syncer.blockSyncReqC <- req
	<-cbi.syncer.reqDone
	cbi.logger.Infof("service selfIndexInEpoch [%v] onReceivedProposal finish sync startlevel [%v] targetlevel [%v]",
		cbi.selfIndexInEpoch, req.startLevel, req.targetLevel)
}

// processQC insert qc and process qc from proposal msg
func (cbi *ConsensusChainedBftImpl) processQC(msg *chainedbftpb.ConsensusMsg) (bool, []int64) {
	proposal := msg.Payload.GetProposalMsg().ProposalData
	syncInfo := msg.Payload.GetProposalMsg().SyncInfo
	cbi.logger.Debugf("processQC start. height: [%d], level: [%d], blockHash: [%x], "+
		"JustifyQC.NewView: [%v]", proposal.Height, proposal.Level, proposal.Block.Header.BlockHash, proposal.JustifyQC.NewView)
	usedTimes := make([]int64, 0, 3)
	validateJustifyQCStart := chainUtils.CurrentTimeMillisSeconds()
	//step1: validate and process qc
	if !cbi.smr.voteRules(proposal.Level, proposal.JustifyQC) {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] validateProposal block [%v:%v] pass "+
			"safety rules check failed", cbi.selfIndexInEpoch, proposal.Height, proposal.Level)
		return false, usedTimes
	}
	validateJustifyQCEnd := chainUtils.CurrentTimeMillisSeconds()
	usedValidateJustifyQC := validateJustifyQCEnd - validateJustifyQCStart
	usedTimes = append(usedTimes, usedValidateJustifyQC)

	if !proposal.JustifyQC.NewView {
		if err := cbi.chainStore.insertQC(proposal.JustifyQC); err != nil {
			cbi.logger.Errorf("insert qc to chainStore failed: %s, qc info: %s", err, proposal.JustifyQC.String())
			return false, usedTimes
		}
	}
	insertJustifyQCEnd := chainUtils.CurrentTimeMillisSeconds()
	usedInsertJustifyQCTime := insertJustifyQCEnd - validateJustifyQCEnd
	usedTimes = append(usedTimes, usedInsertJustifyQCTime)

	if proposal.ProposerIdx != cbi.selfIndexInEpoch {
		//local already handle it when aggregating qc
		cbi.processCertificates(proposal.JustifyQC, syncInfo.HighestTC)
	}
	usedProcessCertificatesEnd := chainUtils.CurrentTimeMillisSeconds() - insertJustifyQCEnd
	usedTimes = append(usedTimes, usedProcessCertificatesEnd)
	cbi.logger.Debugf("check proposal.Level[%d], currentLevel[%d]", proposal.Level, cbi.smr.getCurrentLevel())
	if proposal.Level != cbi.smr.getCurrentLevel() {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] processProposal proposal [%v:%v] does not match the "+
			"smr level [%v:%v], ignore proposal", cbi.selfIndexInEpoch, proposal.Height, proposal.Level,
			cbi.smr.getHeight(), cbi.smr.getCurrentLevel())
		return false, usedTimes
	}
	return true, usedTimes
}

func (cbi *ConsensusChainedBftImpl) validateBlock(proposal *chainedbftpb.ProposalData) bool {
	var (
		err      error
		preBlock *common.Block
	)
	if preBlock, err = cbi.chainStore.getBlock(string(proposal.Block.Header.PreBlockHash), proposal.Height-1); err != nil || preBlock == nil {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] validateProposal failed to get preBlock [%v], err %v",
			cbi.selfIndexInEpoch, proposal.Height-1, err)
		return false
	}
	if !bytes.Equal(preBlock.Header.BlockHash, proposal.JustifyQC.BlockID) {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] validateProposal failed, qc'block not equal to block's prehash",
			cbi.selfIndexInEpoch)
		return false
	}

	if err = cbi.blockVerifier.VerifyBlock(proposal.Block, protocol.CONSENSUS_VERIFY); err != nil {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] processProposal invalid proposal block "+
			"from proposer %v at height [%v] level [%v], err %v", cbi.selfIndexInEpoch,
			proposal.ProposerIdx, proposal.Height, proposal.Level, err)
		return false
	}
	return true
}

func (cbi *ConsensusChainedBftImpl) validateConsensusArg(proposal *chainedbftpb.ProposalData) bool {
	var (
		err           error
		txRWSet       *common.TxRWSet
		consensusArgs *consensus.BlockHeaderConsensusArgs
	)

	if consensusArgs, err = utils.GetConsensusArgsFromBlock(proposal.Block); err != nil {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] processProposal: GetConsensusArgsFromBlock err from proposer %v"+
			" at height [%v] level [%v], err %v", cbi.selfIndexInEpoch, proposal.ProposerIdx, proposal.Height, proposal.Level, err)
		return false
	}
	if txRWSet, err = governance.CheckAndCreateGovernmentArgs(proposal.Block, cbi.store, cbi.proposalCache); err != nil {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] processProposal: CheckAndCreateGovernmentArgs err from proposer"+
			" %v at height [%v] level [%v], err %v", cbi.selfIndexInEpoch, proposal.ProposerIdx, proposal.Height, proposal.Level, err)
		return false
	}

	txRWSetBytes, _ := proto.Marshal(txRWSet)
	ConsensusDataBytes, _ := proto.Marshal(consensusArgs.ConsensusData)
	if !bytes.Equal(txRWSetBytes, ConsensusDataBytes) {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] processProposal: invalid Consensus Args "+
			"from proposer %v at height [%v] level [%v], proposal data:[%v] local data:[%v]", cbi.selfIndexInEpoch,
			proposal.ProposerIdx, proposal.Height, proposal.Level, txRWSet, consensusArgs.ConsensusData)
		return false
	}
	return true
}

func (cbi *ConsensusChainedBftImpl) sendVote2Next(proposal *chainedbftpb.ProposalData, vote *chainedbftpb.ConsensusPayload) {
	nextLeaderIndex := cbi.getProposer(proposal.Level + 1)

	cbi.logger.Debugf("service selfIndexInEpoch [%v] processProposal send vote to next leader [%v]",
		cbi.selfIndexInEpoch, nextLeaderIndex)

	cbi.smr.setLastVote(vote, proposal.Level)
	if nextLeaderIndex == cbi.selfIndexInEpoch {
		consensusMessage := &chainedbftpb.ConsensusMsg{Payload: vote}
		if err := utils.SignConsensusMsg(consensusMessage, cbi.chainConf.ChainConfig().Crypto.Hash, cbi.singer); err != nil {
			cbi.logger.Errorf("sign consensus message failed, err %v", err)
			return
		}
		cbi.logger.Debugf("send vote msg to self[%d], voteHeight:[%d], voteLevel:[%d], voteBlockID:[%x]", cbi.selfIndexInEpoch,
			proposal.Height, proposal.Level, proposal.Block.Header.BlockHash)
		cbi.internalMsgCh <- consensusMessage
	} else {
		cbi.logger.Debugf("send vote msg to other peer [%d], voteHeight:[%d], voteLevel:[%d], voteBlockID:[%x]", nextLeaderIndex,
			proposal.Height, proposal.Level, proposal.Block.Header.BlockHash)
		cbi.signAndSendToPeer(vote, nextLeaderIndex)
	}
}

func (cbi *ConsensusChainedBftImpl) validateVoteData(voteData *chainedbftpb.VoteData) error {
	var (
		err       error
		data      []byte
		author    = voteData.GetAuthor()
		authorIdx = voteData.GetAuthorIdx()
	)
	if author == nil {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] validateVoteData received a "+
			"vote data with nil author", cbi.selfIndexInEpoch)
		return fmt.Errorf("nil author")
	}

	if peer := cbi.smr.getPeerByIndex(authorIdx); peer == nil || peer.id != string(author) {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] validateVoteData received a "+
			"vote data from invalid peer,vote authorIdx [%v]", cbi.selfIndexInEpoch, authorIdx)
		return InvalidPeerErr
	}

	cbi.logger.Debugf("service selfIndexInEpoch [%v] validateVoteData, voteData %v", cbi.selfIndexInEpoch, voteData)
	sign := voteData.Signature
	voteData.Signature = nil
	defer func() {
		voteData.Signature = sign
	}()
	if data, err = proto.Marshal(voteData); err != nil {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] validateVoteData "+
			"marshal vote failed, data %v , err %v", cbi.selfIndexInEpoch, voteData, err)
		return fmt.Errorf("failed to marshal payload")
	}
	if err = utils.VerifyDataSign(data, sign, cbi.accessControlProvider); err != nil {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] validateVoteData "+
			"verify vote failed, data signature, err %v", cbi.selfIndexInEpoch, err)
		return fmt.Errorf("failed to verify voteData signature")
	}
	return nil
}

func (cbi *ConsensusChainedBftImpl) validateVoteMsg(msg *chainedbftpb.ConsensusMsg) error {
	var (
		peer      *peer
		voteMsg   = msg.Payload.GetVoteMsg()
		author    = voteMsg.VoteData.GetAuthor()
		authorIdx = voteMsg.VoteData.GetAuthorIdx()
	)
	if author == nil {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] validateVoteMsg: received a "+
			"vote msg with nil author", cbi.selfIndexInEpoch)
		return fmt.Errorf("nil author")
	}

	if peer = cbi.smr.getPeerByIndex(authorIdx); peer == nil || peer.id != string(author) {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] validateVoteMsg: received a vote msg from invalid peer", cbi.selfIndexInEpoch)
		return InvalidPeerErr
	}
	if err := cbi.validateSignerAndSignature(msg, peer); err != nil {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] validateVoteMsg failed, vote %v, err %v", cbi.selfIndexInEpoch, voteMsg, err)
		return ValidateSignErr
	}

	vote := voteMsg.VoteData
	vote.Signature.Signer = msg.SignEntry.Signer
	if err := cbi.validateVoteData(vote); err != nil {
		return fmt.Errorf("verify vote data failed, err %v", err)
	}
	return nil
}

func (cbi *ConsensusChainedBftImpl) processVote(msg *chainedbftpb.ConsensusMsg) {
	// 1. base check vote msg
	voteMsg := msg.Payload.GetVoteMsg()
	vote := voteMsg.VoteData
	authorIdx := vote.GetAuthorIdx()
	cbi.logger.Debugf("service selfIndexInEpoch [%v] processVote: height [%v] level [%v] epoch [%v] authorIdx [%v],expected [%v:%v:%v]",
		cbi.selfIndexInEpoch, vote.Height, vote.Level, vote.EpochId, authorIdx, cbi.smr.getHeight(), cbi.smr.getCurrentLevel(), cbi.smr.getEpochId())
	if vote.Height < cbi.smr.getHeight() || vote.Level < cbi.smr.getCurrentLevel() || vote.EpochId != cbi.smr.getEpochId() {
		cbi.logger.Infof("service selfIndexInEpoch [%v] processVote: received vote "+
			"at wrong height [%v] level [%v] epoch [%v], expected [%v:%v:%v] authorIdx [%v]", cbi.selfIndexInEpoch, vote.Height,
			vote.Level, vote.EpochId, cbi.smr.getHeight(), cbi.smr.getCurrentLevel(), cbi.smr.getEpochId(), authorIdx)
		return
	}
	cbi.logger.Debugf("process vote step 1 validate base vote info, voteHeight:%d, voteLevel:%d, voteBlockID:%x, authorIndex:%d",
		vote.Height, vote.Level, vote.BlockID, vote.AuthorIdx)
	validateVoteMsgBegin := chainUtils.CurrentTimeMillisSeconds()
	if err := cbi.validateVoteMsg(msg); err != nil {
		return
	}
	validateVoteMsgEnd := chainUtils.CurrentTimeMillisSeconds()
	usedValidateVoteTime := validateVoteMsgEnd - validateVoteMsgBegin

	// 2. only proposer could handle proposal’s vote
	cbi.logger.Debugf("process vote step 2 only proposer will process vote with Proposal type or all peer can process vote with NewView type")
	if !vote.NewView {
		//regular votes are sent to the leaders of the next round only.
		if nextLeaderIndex := cbi.getProposer(vote.Level + 1); nextLeaderIndex != cbi.selfIndexInEpoch {
			cbi.logger.Errorf("service selfIndexInEpoch [%v] processVote: self is not next "+
				"leader[%d] for level [%v]", cbi.selfIndexInEpoch, nextLeaderIndex, vote.Level+1)
			return
		}
	}

	var (
		err      error
		need     bool
		inserted bool
	)
	// 3. Add vote to msgPool
	cbi.logger.Debugf("process vote step 3 check whether need sync from other peer")
	if need, err = cbi.needFetch(voteMsg.SyncInfo); err != nil {
		cbi.logger.Errorf("processVote: needFetch failed, reason: %s", err)
		return
	}
	needFetchEnd := chainUtils.CurrentTimeMillisSeconds()
	usedNeedFetchTime := needFetchEnd - validateVoteMsgEnd
	cbi.logger.Debugf("process vote step 4 inserted the vote ")
	if inserted, err = cbi.msgPool.InsertVote(vote.Height, vote.Level, msg); err != nil {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] processVote: insert vote msg failed, err %v ",
			cbi.selfIndexInEpoch, err)
		return
	} else if !inserted {
		cbi.logger.Warnf("service selfIndexInEpoch [%v] processVote: inserted"+
			" vote msg from validator failed, %v", cbi.selfIndexInEpoch, authorIdx)
		return
	}
	insertVoteMsgEnd := chainUtils.CurrentTimeMillisSeconds()
	usedInsertVoteMsgTime := insertVoteMsgEnd - needFetchEnd

	// 4. fetch data from peers and handle qc in voteInfo
	usedFetchAndHandleQcTime := int64(0)
	cbi.logger.Debugf("process vote step 5 whether need fetch info: [%v] and process", need)
	if need {
		cbi.fetchAndHandleQc(authorIdx, voteMsg)
		usedFetchAndHandleQcTime = chainUtils.CurrentTimeMillisSeconds() - insertVoteMsgEnd
	} else {
		// 5. generate QC if majority are voted and process the new QC if don't need sync data from peers
		cbi.processVotes(vote)
		usedFetchAndHandleQcTime = chainUtils.CurrentTimeMillisSeconds() - insertVoteMsgEnd
	}
	cbi.logger.Debugf("time cost in processVote: validateVoteMsgTime: %d,"+
		" needFetchTime: %d, insertVoteMsgTime: %d, needFetch: %v, handleQCTime: %d.",
		usedValidateVoteTime, usedNeedFetchTime, usedInsertVoteMsgTime, need, usedFetchAndHandleQcTime)
}

//fetchAndHandleQc Fetch the missing block data and the  process the received QC until the data is all fetched.
func (cbi *ConsensusChainedBftImpl) fetchAndHandleQc(authorIdx uint64, voteMsg *chainedbftpb.VoteMsg) {
	cbi.logger.Infof("service selfIndexInEpoch [%v] processVote: need sync up to [%v:%v]",
		cbi.selfIndexInEpoch, voteMsg.SyncInfo.HighestQC.Height, voteMsg.SyncInfo.HighestQC.Level)
	req := &blockSyncReq{
		targetPeer:  authorIdx,
		height:      voteMsg.SyncInfo.HighestQC.Height,
		blockID:     voteMsg.SyncInfo.HighestQC.BlockID,
		targetLevel: voteMsg.SyncInfo.HighestQC.Level,
		startLevel:  cbi.chainStore.getCurrentQC().Level + 1,
	}
	cbi.syncer.blockSyncReqC <- req
	<-cbi.syncer.reqDone
	cbi.logger.Debugf("service selfIndexInEpoch [%v] processVote: finish sync startlevel [%v] targetlevel [%v]",
		cbi.selfIndexInEpoch, req.startLevel, req.targetLevel)
	if cbi.smr.getCurrentLevel() < req.targetLevel {
		cbi.logger.Infof("service index [%v] processVote: sync currentLevel [%v] not catch targetLevel [%v]",
			cbi.selfIndexInEpoch, cbi.smr.getCurrentLevel(), req.targetLevel)
		return
	}
	cbi.processCertificates(voteMsg.SyncInfo.HighestQC, voteMsg.SyncInfo.HighestTC)
}

//processVotes QC is generated if a majority are voted for the special Height and Level.
func (cbi *ConsensusChainedBftImpl) processVotes(vote *chainedbftpb.VoteData) {
	blockID, newView, done := cbi.msgPool.CheckVotesDone(vote.Height, vote.Level)
	if !done {
		return
	}
	//aggregate qc
	cbi.logger.Debugf("service selfIndexInEpoch [%v] processVote: new qc aggregated for height [%v] level [%v]",
		cbi.selfIndexInEpoch, vote.Height, vote.Level)
	votes := cbi.msgPool.GetVotes(vote.Height, vote.Level)
	qc := &chainedbftpb.QuorumCert{
		BlockID: blockID,
		Height:  vote.Height,
		Level:   vote.Level,
		Votes:   votes,
		NewView: newView,
		EpochId: cbi.smr.getEpochId(),
	}
	if blockID != nil {
		if err := cbi.chainStore.insertQC(qc); err != nil {
			cbi.logger.Errorf("service index [%v] processVote: new qc aggregated for height [%v] level [%v] blockId [%x], err=%v",
				cbi.selfIndexInEpoch, vote.Height, vote.Level, blockID, err)
			return
		}
	}
	var tc *chainedbftpb.QuorumCert
	if qc.NewView {
		// If the newly generated QC type is NewView, it means that majority agree on the timeout and assign QC to TC
		tc = qc
	}
	cbi.processCertificates(qc, tc)
}

// processCertificates
// qc When processing a proposalMsg or voteMsg, the tc information is contained in the incoming message;
// in other cases, the parameter is currentQC in local node.
// tc When processing a proposalMsg or voteMsg, the tc information is contained in the incoming message;
// in other cases, the parameter is nil.
func (cbi *ConsensusChainedBftImpl) processCertificates(qc *chainedbftpb.QuorumCert, tc *chainedbftpb.QuorumCert) {
	cbi.logger.Debugf("service selfIndexInEpoch [%v] processCertificates start:, height [%v], level [%v],qc.Height "+
		"[%v] qc.Level [%v], qc.epochID [%d]", cbi.selfIndexInEpoch, cbi.smr.getHeight(), cbi.smr.getCurrentLevel(), qc.Height, qc.Level, qc.EpochId)
	updateLockedQCStart := chainUtils.CurrentTimeMillisSeconds()
	cbi.smr.updateLockedQC(qc)
	updateLockedQCEnd := chainUtils.CurrentTimeMillisSeconds()
	usedLockedQCTime := updateLockedQCEnd - updateLockedQCStart
	cbi.commitBlocksByQC(qc)
	commitBlockEnd := chainUtils.CurrentTimeMillisSeconds()
	usedCommitBlockTime := commitBlockEnd - updateLockedQCEnd
	var (
		tcLevel        = uint64(0)
		currentQC      = qc
		committedLevel = cbi.smr.getLastCommittedLevel()
	)
	if tc != nil {
		tcLevel = tc.Level
		cbi.smr.updateTC(tc)
		currentQC = cbi.chainStore.getCurrentQC()
	}
	if newLevel := cbi.smr.processCertificates(qc.Height, currentQC.Level, tcLevel, committedLevel); newLevel {
		cbi.smr.updateState(chainedbftpb.ConsStateType_NewHeight)
		cbi.processNewHeight(cbi.smr.getHeight(), cbi.smr.getCurrentLevel())
	}
	processQCAndTCEnd := chainUtils.CurrentTimeMillisSeconds()
	usedProcessQCAndTCTime := processQCAndTCEnd - commitBlockEnd
	totalUsedTime := processQCAndTCEnd - updateLockedQCStart
	cbi.logger.Debugf("time costs: lockedQCTime: %d, commitBlockTime: %d, processQCTCTime: %d,"+
		" totalTime: %d", usedLockedQCTime, usedCommitBlockTime, usedProcessQCAndTCTime, totalUsedTime)
}

func (cbi *ConsensusChainedBftImpl) commitBlocksByQC(qc *chainedbftpb.QuorumCert) {
	commit, block, level := cbi.smr.commitRules(qc)
	if !commit {
		return
	}

	cbi.logger.Debugf("service selfIndexInEpoch [%v] processCertificates: commitRules success, height [%v], level [%v],"+
		" committed level [%v]", cbi.selfIndexInEpoch, block.Header.BlockHeight, level, cbi.chainStore.getCommitLevel())
	if level > cbi.chainStore.getCommitLevel() {
		cbi.logger.Debugf("service selfIndexInEpoch [%v] processCertificates: try committing a block %x on [%d:%v]",
			cbi.selfIndexInEpoch, block.Header.BlockHash, block.Header.BlockHeight, level)
		lastCommitted, err := cbi.chainStore.commitBlock(block)
		if lastCommitted != nil {
			level, err := utils.GetLevelFromBlock(lastCommitted)
			if err != nil {
				cbi.logger.Errorf("GetLevelFromBlock failed, reason: %s", err)
				return
			}
			cbi.smr.setLastCommittedBlock(lastCommitted, level)
			cbi.msgPool.OnBlockSealed(uint64(lastCommitted.Header.BlockHeight))
		}
		if err != nil {
			cbi.logger.Errorf("commit block to the chain failed, reason: %s", err)
		}
	}
}

func (cbi *ConsensusChainedBftImpl) processBlockCommitted(block *common.Block) {
	cbi.logger.Debugf("processBlockCommitted received has committed block, height:%d, hash:%x",
		block.Header.BlockHeight, block.Header.BlockHash)
	// 1. check base commit block info
	if int64(cbi.commitHeight) >= block.Header.BlockHeight {
		cbi.logger.Debugf("service selfIndexInEpoch [%v] persisted block [%v] vs commit height [%v]",
			cbi.selfIndexInEpoch, block.Header.BlockHeight, cbi.commitHeight)
		return
	}
	// 2. insert committed block to chainStore
	updateStateBegin := chainUtils.CurrentTimeMillisSeconds()
	cbi.logger.Debugf("processBlockCommitted step 1 insert complete block")
	if err := cbi.chainStore.insertCompletedBlock(block); err != nil {
		cbi.logger.Errorf("insert block[%d: %x] to chainStore failed", block.Header.BlockHeight, block.Header.BlockHash)
		return
	}
	committedBlock, err := cbi.chainStore.getBlock(string(block.Header.BlockHash), uint64(block.Header.BlockHeight))
	if err != nil {
		cbi.logger.Errorf("get block[%d: %x] from chainStore failed, reason: %s", block.Header.BlockHeight, block.Header.BlockHash, err)
		return
	}
	// 3. update commit info in the consensus
	cbi.logger.Debugf("processBlockCommitted step 2 update the last committed block info")
	cbi.smr.setLastCommittedBlock(committedBlock, cbi.chainStore.getCommitLevel())
	cbi.commitHeight = uint64(block.Header.BlockHeight)
	height := uint64(block.Header.BlockHeight)
	updateStateEnd := chainUtils.CurrentTimeMillisSeconds()
	usedUpdateTime := updateStateEnd - updateStateBegin
	// 4. create next epoch if meet the condition
	cbi.logger.Debugf("processBlockCommitted step 3 create epoch if meet the condition")
	if err := cbi.createNextEpochIfRequired(height); err != nil {
		cbi.logger.Errorf("failed to createNextEpochIfRequired, err %v", err)
		return
	}
	createNextEpochEnd := chainUtils.CurrentTimeMillisSeconds()
	usedCreateNextEpochTime := createNextEpochEnd - updateStateEnd

	// 5. check if need to switch with the epoch
	if cbi.nextEpoch == nil || (cbi.nextEpoch != nil && cbi.nextEpoch.switchHeight > height) {
		cbi.logger.Debugf("processBlockCommitted step 4 no switch epoch and process qc")
		return
	}

	// 6. switch epoch and update field state in consensus
	oldIndex := cbi.selfIndexInEpoch
	cbi.logger.Debugf("processBlockCommitted step 5 switch epoch and process qc")
	if err := cbi.switchNextEpoch(height); err != nil {
		return
	}
	if cbi.smr.isValidIdx(cbi.selfIndexInEpoch) {
		cbi.logger.Infof("service selfIndexInEpoch [%v] start processCertificates,"+
			"height [%v],level [%v]", cbi.selfIndexInEpoch, cbi.smr.getHeight(), cbi.smr.getCurrentLevel())
	} else if oldIndex != cbi.selfIndexInEpoch {
		if oldIndex == utils.InvalidIndex {
			cbi.logger.Infof("service selfIndexInEpoch [%v] got a chance to join consensus group", cbi.selfIndexInEpoch)
		} else {
			cbi.logger.Infof("service old selfIndexInEpoch [%v] next selfIndexInEpoch [%v] leave consensus group",
				oldIndex, cbi.selfIndexInEpoch)
		}
	}
	switchNextEpochEnd := chainUtils.CurrentTimeMillisSeconds()
	usedSwitchNextEpochTime := switchNextEpochEnd - createNextEpochEnd
	// 7. process qc: commitHeight+1
	cbi.processCertificates(cbi.chainStore.getCurrentQC(), nil)
	usedProcessQCTCTime := chainUtils.CurrentTimeMillisSeconds() - switchNextEpochEnd
	cbi.logger.Infof("processBlockCommitted end, blockHeight: [%d], hash: [%x] ...", height, block.Header.BlockHash)
	cbi.logger.Debugf("time costs in processBlockCommitted: updateStateTime: %d, createNextEpochTime: %d, "+
		"switchNextEpochTime: %d, processQCTCTime: %d", usedUpdateTime, usedCreateNextEpochTime, usedSwitchNextEpochTime, usedProcessQCTCTime)
}

func (cbi *ConsensusChainedBftImpl) switchNextEpoch(blockHeight uint64) error {
	cbi.logger.Debugf("service [%v] handle block committed: "+
		"start switching to next epoch at height [%v]", cbi.selfIndexInEpoch, blockHeight)
	chainStore, err := openChainStore(cbi.ledgerCache, cbi.blockCommitter, cbi.store, cbi, cbi.logger)
	if err != nil {
		cbi.logger.Errorf("new consensus service failed, err %v", err)
		return err
	}

	cbi.Lock()
	defer cbi.Unlock()
	if cbi.timerService != nil {
		cbi.timerService.Stop()
	}
	if cbi.syncer != nil {
		cbi.syncer.stop()
	}
	cbi.chainStore = chainStore
	cbi.syncer = newSyncManager(cbi)
	cbi.msgPool = cbi.nextEpoch.msgPool
	cbi.timerService = timeservice.NewTimerService()
	cbi.selfIndexInEpoch = cbi.nextEpoch.index
	cbi.smr = newChainedBftSMR(cbi.chainID, cbi.nextEpoch, cbi.chainStore, cbi.timerService)
	cbi.nextEpoch = nil
	cbi.logger.Debugf("create epoch height: [%d], SMR height: [%d]", blockHeight, cbi.smr.getHeight())
	go cbi.timerService.Start()
	go cbi.syncer.start()
	cbi.helper.DiscardAboveHeight(int64(blockHeight))
	return nil
}

func (cbi *ConsensusChainedBftImpl) validateBlockFetch(msg *chainedbftpb.ConsensusMsg) error {
	req := msg.Payload.GetBlockFetchMsg()
	authorIdx := req.GetAuthorIdx()
	peer := cbi.smr.getPeerByIndex(authorIdx)
	if peer == nil {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] validateBlockFetch: received a vote msg from invalid peer", cbi.selfIndexInEpoch)
		return InvalidPeerErr
	}

	err := cbi.validateSignerAndSignature(msg, peer)
	if err != nil {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] validateBlockFetch failed, err %v"+
			" fetch req %v, err %v", cbi.selfIndexInEpoch, req, err)
		return ValidateSignErr
	}

	if req.NumBlocks > MaxSyncBlockNum {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] validateBlockFetch: fetch too many blocks %v",
			cbi.selfIndexInEpoch, req.NumBlocks)
		return fmt.Errorf("fetch too many blocks")
	}
	return nil
}

func (cbi *ConsensusChainedBftImpl) processBlockFetch(msg *chainedbftpb.ConsensusMsg) {
	if err := cbi.validateBlockFetch(msg); err != nil {
		return
	}
	var (
		req    = msg.Payload.GetBlockFetchMsg()
		blocks = make([]*chainedbftpb.BlockPair, 0, req.NumBlocks)

		id        = string(req.BlockID)
		height    = req.Height
		status    = chainedbftpb.BlockFetchStatus_Succeeded
		authorIdx = req.GetAuthorIdx()
	)

	for i := 0; i < int(req.NumBlocks); i++ {
		block, _ := cbi.chainStore.getBlock(id, height)
		qc, _ := cbi.chainStore.getQC(id, height)
		if block == nil || qc == nil {
			status = chainedbftpb.BlockFetchStatus_NotEnoughBlocks
			break
		}
		//clone for marshall
		newBlock := proto.Clone(block).(*common.Block)
		newQc := proto.Clone(qc).(*chainedbftpb.QuorumCert)
		blockPair := &chainedbftpb.BlockPair{
			Block: newBlock,
			QC:    newQc,
		}
		height = height - 1
		id = string(newBlock.Header.PreBlockHash)
		blocks = append(blocks, blockPair)
	}
	if len(blocks) == 0 {
		status = chainedbftpb.BlockFetchStatus_IdNotFound
	}
	rsp := cbi.constructBlockFetchRespMsg(blocks, status)
	cbi.signAndSendToPeer(rsp, authorIdx)
}

func (cbi *ConsensusChainedBftImpl) validateBlockFetchRsp(msg *chainedbftpb.ConsensusMsg) error {
	rsp := msg.Payload.GetBlockFetchRespMsg()
	authorIdx := rsp.GetAuthorIdx()
	peer := cbi.smr.getPeerByIndex(authorIdx)
	if peer == nil {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] validateBlockFetchRsp: received a vote msg from invalid peer", cbi.selfIndexInEpoch)
		return InvalidPeerErr
	}

	if err := cbi.validateSignerAndSignature(msg, peer); err != nil {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] from %v validateBlockFetchRsp failed, err %v"+
			" fetch rsp %v, err %v", cbi.selfIndexInEpoch, rsp.AuthorIdx, rsp, err)
		return ValidateSignErr
	}
	cbi.logger.Infof("service selfIndexInEpoch [%v] from %v validateBlockFetchRsp success %v"+
		" fetch rsp %v ", cbi.selfIndexInEpoch, rsp.AuthorIdx, rsp)
	return nil
}

func (cbi *ConsensusChainedBftImpl) startTimer(event *timeservice.TimerEvent) {
	cbi.timerService.AddEvent(event)
}

//validateSignerAndSignature validate msg signer and signatures
func (cbi *ConsensusChainedBftImpl) validateSignerAndSignature(msg *chainedbftpb.ConsensusMsg, peer *peer) error {
	// check cert id
	// if msg generated by self, check signer, because netService not mapping self nodeId
	if peer.index == cbi.selfIndexInEpoch {
		member, err := cbi.accessControlProvider.NewMemberFromProto(msg.SignEntry.Signer)
		if err != nil {
			cbi.logger.Errorf("service selfIndexInEpoch [%v] validateSignerAndSignature failed,: new member by msg.SignEntry.Signer err %v",
				cbi.selfIndexInEpoch, err)
			return VerifySignerFailedErr
		}
		if cbi.singer.GetOrgId() != member.GetOrgId() {
			cbi.logger.Errorf("service selfIndexInEpoch [%v] validateSignerAndSignature failed,: match msg.signer and msg.payload.selfIndexInEpoch signer",
				cbi.selfIndexInEpoch)
			return VerifySignerFailedErr
		}
	} else { // otherwise, check nodeId by signer.cert
		nodeId, err := utils.GetUidFromProtoSigner(msg.SignEntry.Signer, cbi.netService, cbi.accessControlProvider)
		if err != nil {
			cbi.logger.Errorf("service selfIndexInEpoch [%v] validateSignerAndSignature failed,: get nodeId by msg.SignEntry.Signer err %v",
				cbi.selfIndexInEpoch, err)
			return VerifySignerFailedErr
		}
		if nodeId != peer.id {
			cbi.logger.Errorf("service selfIndexInEpoch [%v] validateSignerAndSignature failed,: match msg.signer and msg.payload.selfIndexInEpoch signer "+
				"nodeId %v, payload nodeId %v", cbi.selfIndexInEpoch, nodeId, peer.id)
			return VerifySignerFailedErr
		}
	}
	//check sign
	if err := utils.VerifyConsensusMsgSign(msg, cbi.accessControlProvider); err != nil {
		cbi.logger.Errorf("service selfIndexInEpoch [%v] validateSignerAndSignature failed,: verify "+
			" msg, err %v", cbi.selfIndexInEpoch, err)
		return fmt.Errorf("verify signature failed")
	}
	return nil
}
