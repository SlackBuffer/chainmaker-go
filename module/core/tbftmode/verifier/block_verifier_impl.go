/*
Copyright (C) BABEC. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package verifier

import (
	"chainmaker.org/chainmaker-go/core/common"
	"chainmaker.org/chainmaker-go/store/statedb/statesqldb"
	"encoding/hex"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"sync"

	commonErrors "chainmaker.org/chainmaker-go/common/errors"
	"chainmaker.org/chainmaker-go/common/msgbus"
	"chainmaker.org/chainmaker-go/consensus"
	"chainmaker.org/chainmaker-go/localconf"
	"chainmaker.org/chainmaker-go/monitor"
	commonpb "chainmaker.org/chainmaker-go/pb/protogo/common"
	consensuspb "chainmaker.org/chainmaker-go/pb/protogo/consensus"
	"chainmaker.org/chainmaker-go/protocol"
	"chainmaker.org/chainmaker-go/utils"
	"github.com/prometheus/client_golang/prometheus"
)

const LOCKED = "LOCKED" // LOCKED mark

// BlockVerifierImpl implements BlockVerifier interface.
// Verify block and transactions.
type BlockVerifierImpl struct {
	chainId         string                   // chain id, to identity this chain
	msgBus          msgbus.MessageBus        // message bus
	txScheduler     protocol.TxScheduler     // scheduler orders tx batch into DAG form and returns a block
	snapshotManager protocol.SnapshotManager // snapshot manager
	ledgerCache     protocol.LedgerCache     // ledger cache
	blockchainStore protocol.BlockchainStore // blockchain store

	reentrantLocks *reentrantLocks                // reentrant lock for avoid concurrent verify block
	proposalCache  protocol.ProposalCache         // proposal cache
	chainConf      protocol.ChainConf             // chain config
	ac             protocol.AccessControlProvider // access control manager
	log            protocol.Logger                // logger
	txPool         protocol.TxPool                // tx pool to check if tx is duplicate
	mu             sync.Mutex                     // to avoid concurrent map modify
	verifierBlock  *common.VerifierBlock

	metricBlockVerifyTime *prometheus.HistogramVec // metrics monitor
}

type BlockVerifierConfig struct {
	ChainId         string
	MsgBus          msgbus.MessageBus
	SnapshotManager protocol.SnapshotManager
	BlockchainStore protocol.BlockchainStore
	LedgerCache     protocol.LedgerCache
	TxScheduler     protocol.TxScheduler
	ProposedCache   protocol.ProposalCache
	ChainConf       protocol.ChainConf
	AC              protocol.AccessControlProvider
	TxPool          protocol.TxPool
	VmMgr           protocol.VmManager
}

func NewBlockVerifier(config BlockVerifierConfig, log protocol.Logger) (protocol.BlockVerifier, error) {
	v := &BlockVerifierImpl{
		chainId:         config.ChainId,
		msgBus:          config.MsgBus,
		txScheduler:     config.TxScheduler,
		snapshotManager: config.SnapshotManager,
		ledgerCache:     config.LedgerCache,
		blockchainStore: config.BlockchainStore,
		reentrantLocks: &reentrantLocks{
			reentrantLocks: make(map[string]interface{}),
		},
		proposalCache: config.ProposedCache,
		chainConf:     config.ChainConf,
		ac:            config.AC,
		log:           log,
		txPool:        config.TxPool,
	}

	conf := &common.VerifierBlockConf{
		ChainConf:       config.ChainConf,
		Log:             v.log,
		LedgerCache:     config.LedgerCache,
		Ac:              config.AC,
		SnapshotManager: config.SnapshotManager,
		VmMgr:           config.VmMgr,
		TxPool:          config.TxPool,
		BlockchainStore: config.BlockchainStore,
	}
	v.verifierBlock = common.NewVerifierBlock(conf)

	if localconf.ChainMakerConfig.MonitorConfig.Enabled {
		v.metricBlockVerifyTime = monitor.NewHistogramVec(monitor.SUBSYSTEM_CORE_VERIFIER, "metric_block_verify_time",
			"block verify time metric", []float64{0.005, 0.01, 0.015, 0.05, 0.1, 1, 10}, "chainId")
	}

	return v, nil
}

// VerifyBlock, to check if block is valid
func (v *BlockVerifierImpl) VerifyBlock(block *commonpb.Block, mode protocol.VerifyMode) (err error) {

	startTick := utils.CurrentTimeMillisSeconds()
	if err = utils.IsEmptyBlock(block); err != nil {
		v.log.Error(err)
		return err
	}

	v.log.Debugf("verify receive [%d](%x,%d,%d), from sync %d",
		block.Header.BlockHeight, block.Header.BlockHash, block.Header.TxCount, len(block.Txs), mode)
	// avoid concurrent verify, only one block hash can be verified at the same time
	if !v.reentrantLocks.lock(string(block.Header.BlockHash)) {
		v.log.Warnf("block(%d,%x) concurrent verify, yield", block.Header.BlockHeight, block.Header.BlockHash)
		return commonErrors.ErrConcurrentVerify
	}
	defer v.reentrantLocks.unlock(string(block.Header.BlockHash))

	var isValid bool
	var contractEventMap map[string][]*commonpb.ContractEvent
	// to check if the block has verified before
	b, txRwSet, EventMap := v.proposalCache.GetProposedBlock(block)
	contractEventMap = EventMap

	if b != nil {
		isSqlDb := v.chainConf.ChainConfig().Contract.EnableSqlSupport
		notSolo := consensuspb.ConsensusType_SOLO != v.chainConf.ChainConfig().Consensus.Type
		if notSolo || isSqlDb {
			// the block has verified befo
			// the block has verified before
			v.log.Infof("verify success repeat [%d](%x)", block.Header.BlockHeight, block.Header.BlockHash)
			isValid = true
			if protocol.CONSENSUS_VERIFY == mode {
				// consensus mode, publish verify result to message bus
				v.msgBus.Publish(msgbus.VerifyResult, parseVerifyResult(block, isValid))
			}
			lastBlock, _ := v.proposalCache.GetProposedBlockByHashAndHeight(block.Header.PreBlockHash, block.Header.BlockHeight-1)
			if lastBlock == nil {
				v.log.Debugf("no pre-block be found, preHeight:%d, preBlockHash:%x", block.Header.BlockHeight-1, block.Header.PreBlockHash)
				return nil
			}
			cutBlocks := v.proposalCache.KeepProposedBlock(lastBlock.Header.BlockHash, lastBlock.Header.BlockHeight)
			if len(cutBlocks) > 0 {
				v.log.Infof("cut block block hash: %s, height: %v", hex.EncodeToString(lastBlock.Header.BlockHash), lastBlock.Header.BlockHeight)
				v.cutBlocks(cutBlocks, lastBlock)
			}
			err := v.proposalCache.SetProposedBlock(block, txRwSet, EventMap, v.proposalCache.IsProposedAt(block.Header.BlockHeight))
			return err
		}
	}

	txRWSetMap, contractEventMap, timeLasts, err := v.validateBlock(block)
	if err != nil {
		v.log.Warnf("verify failed [%d](%x),preBlockHash:%x, %s",
			block.Header.BlockHeight, block.Header.BlockHash, block.Header.PreBlockHash, err.Error())
		if protocol.CONSENSUS_VERIFY == mode {
			v.msgBus.Publish(msgbus.VerifyResult, parseVerifyResult(block, isValid))
		}

		// rollback sql
		if v.chainConf.ChainConfig().Contract.EnableSqlSupport {
			_ = v.blockchainStore.RollbackDbTransaction(block.GetTxKey())
			// drop database if create contract fail
			if len(block.Txs) == 0 && utils.IsManageContractAsConfigTx(block.Txs[0], true) {
				var payload commonpb.ContractMgmtPayload
				if err := proto.Unmarshal(block.Txs[0].RequestPayload, &payload); err == nil {
					if payload.ContractId != nil {
						dbName := statesqldb.GetContractDbName(v.chainId, payload.ContractId.ContractName)
						v.blockchainStore.ExecDdlSql(payload.ContractId.ContractName, "drop database "+dbName)
					}
				}
			}
		}
		return err
	}

	// sync mode, need to verify consensus vote signature
	if protocol.SYNC_VERIFY == mode {
		if err = v.verifyVoteSig(block); err != nil {
			v.log.Warnf("verify failed [%d](%x), votesig %s",
				block.Header.BlockHeight, block.Header.BlockHash, err.Error())
			return err
		}
	}

	// verify success, cache block and read write set
	v.log.Debugf("set proposed block(%d,%x)", block.Header.BlockHeight, block.Header.BlockHash)
	if err = v.proposalCache.SetProposedBlock(block, txRWSetMap, contractEventMap, false); err != nil {
		return err
	}

	// mark transactions in block as pending status in txpool
	v.txPool.AddTxsToPendingCache(block.Txs, block.Header.BlockHeight)

	isValid = true
	if protocol.CONSENSUS_VERIFY == mode {
		v.msgBus.Publish(msgbus.VerifyResult, parseVerifyResult(block, isValid))
	}
	elapsed := utils.CurrentTimeMillisSeconds() - startTick
	v.log.Infof("verify success [%d,%x](%v,%d)", block.Header.BlockHeight, block.Header.BlockHash,
		timeLasts, elapsed)
	if localconf.ChainMakerConfig.MonitorConfig.Enabled {
		v.metricBlockVerifyTime.WithLabelValues(v.chainId).Observe(float64(elapsed) / 1000)
	}
	return nil
}

func (v *BlockVerifierImpl) validateBlock(block *commonpb.Block) (map[string]*commonpb.TxRWSet, map[string][]*commonpb.ContractEvent, []int64, error) {
	hashType := v.chainConf.ChainConfig().Crypto.Hash
	timeLasts := make([]int64, 0)
	var err error
	var lastBlock *commonpb.Block
	txCapacity := int64(v.chainConf.ChainConfig().Block.BlockTxCapacity)
	if block.Header.TxCount > txCapacity {
		return nil, nil, timeLasts, fmt.Errorf("txcapacity expect <= %d, got %d)", txCapacity, block.Header.TxCount)
	}

	if err = common.IsTxCountValid(block); err != nil {
		return nil, nil, timeLasts, err
	}

	lastBlock, err = v.verifierBlock.FetchLastBlock(block, lastBlock)
	if err != nil {
		return nil, nil, timeLasts, err
	}
	// proposed height == proposing height - 1
	proposedHeight := lastBlock.Header.BlockHeight
	// check if this block height is 1 bigger than last block height
	lastBlockHash := lastBlock.Header.BlockHash
	err = common.CheckPreBlock(block, lastBlock, err, lastBlockHash, proposedHeight)
	if err != nil {
		return nil, nil, timeLasts, err
	}

	return v.verifierBlock.ValidateBlock(block, lastBlock, hashType, timeLasts)
}

func (v *BlockVerifierImpl) verifyVoteSig(block *commonpb.Block) error {
	return consensus.VerifyBlockSignatures(v.chainConf, v.ac, v.blockchainStore, block, v.ledgerCache)
}

func parseVerifyResult(block *commonpb.Block, isValid bool) *consensuspb.VerifyResult {
	verifyResult := &consensuspb.VerifyResult{
		VerifiedBlock: block,
	}
	if isValid {
		verifyResult.Code = consensuspb.VerifyResult_SUCCESS
		verifyResult.Msg = "OK"
	} else {
		verifyResult.Msg = "FAIL"
		verifyResult.Code = consensuspb.VerifyResult_FAIL
	}
	return verifyResult
}

func (v *BlockVerifierImpl) cutBlocks(blocksToCut []*commonpb.Block, blockToKeep *commonpb.Block) {
	cutTxs := make([]*commonpb.Transaction, 0)
	txMap := make(map[string]interface{})
	for _, tx := range blockToKeep.Txs {
		txMap[tx.Header.TxId] = struct{}{}
	}
	for _, blockToCut := range blocksToCut {
		v.log.Infof("cut block block hash: %s, height: %v", blockToCut.Header.BlockHash, blockToCut.Header.BlockHeight)
		for _, txToCut := range blockToCut.Txs {
			if _, ok := txMap[txToCut.Header.TxId]; ok {
				// this transaction is kept, do NOT cut it.
				continue
			}
			v.log.Debugf("cut tx hash: %s", txToCut.Header.TxId)
			cutTxs = append(cutTxs, txToCut)
		}
	}
	if len(cutTxs) > 0 {
		v.txPool.RetryAndRemoveTxs(cutTxs, nil)
	}
}
