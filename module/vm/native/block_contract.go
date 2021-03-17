/*
Copyright (C) BABEC. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package native

import (
	"chainmaker.org/chainmaker-go/localconf"
	"chainmaker.org/chainmaker-go/logger"
	commonPb "chainmaker.org/chainmaker-go/pb/protogo/common"
	discoveryPb "chainmaker.org/chainmaker-go/pb/protogo/discovery"
	"chainmaker.org/chainmaker-go/protocol"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"reflect"
	"strconv"
	"strings"
)

const (
	paramNameBlockHeight = "blockHeight"
	paramNameWithRWSet   = "withRWSet"
	paramNameBlockHash   = "blockHash"
	paramNameTxId        = "txId"
)

var (
	logTemplateMarshalBlockInfoFailed = "marshal block info failed, %s"
	errStoreIsNil                     = fmt.Errorf("store is nil")
)

type BlockContact struct {
	methods map[string]ContractFunc
	log     *logger.CMLogger
}

func newBlockContact(log *logger.CMLogger) *BlockContact {
	return &BlockContact{
		log:     log,
		methods: registerBlockContactMethods(log),
	}
}

func (c *BlockContact) getMethod(methodName string) ContractFunc {
	return c.methods[methodName]
}

func registerBlockContactMethods(log *logger.CMLogger) map[string]ContractFunc {
	queryMethodMap := make(map[string]ContractFunc, 64)
	blockRuntime := &BlockRuntime{log: log}
	queryMethodMap[commonPb.QueryFunction_GET_BLOCK_BY_HEIGHT.String()] = blockRuntime.GetBlockByHeight
	queryMethodMap[commonPb.QueryFunction_GET_BLOCK_WITH_TXRWSETS_BY_HEIGHT.String()] = blockRuntime.GetBlockWithTxRWSetsByHeight
	queryMethodMap[commonPb.QueryFunction_GET_BLOCK_BY_HASH.String()] = blockRuntime.GetBlockByHash
	queryMethodMap[commonPb.QueryFunction_GET_BLOCK_WITH_TXRWSETS_BY_HASH.String()] = blockRuntime.GetBlockWithTxRWSetsByHash
	queryMethodMap[commonPb.QueryFunction_GET_BLOCK_BY_TX_ID.String()] = blockRuntime.GetBlockByTxId
	queryMethodMap[commonPb.QueryFunction_GET_TX_BY_TX_ID.String()] = blockRuntime.GetTxByTxId
	queryMethodMap[commonPb.QueryFunction_GET_LAST_CONFIG_BLOCK.String()] = blockRuntime.GetLastConfigBlock
	queryMethodMap[commonPb.QueryFunction_GET_LAST_BLOCK.String()] = blockRuntime.GetLastBlock
	queryMethodMap[commonPb.QueryFunction_GET_CHAIN_INFO.String()] = blockRuntime.GetChainInfo
	queryMethodMap[commonPb.QueryFunction_GET_NODE_CHAIN_LIST.String()] = blockRuntime.GetNodeChainList
	return queryMethodMap
}

type BlockRuntime struct {
	height    int64
	withRWSet string
	hash      string
	txId      string
	log       *logger.CMLogger
}

// GetNodeChainList return list of chain
func (r *BlockRuntime) GetNodeChainList(txSimContext protocol.TxSimContext, parameters map[string]string) ([]byte, error) {
	var errMsg string
	var err error

	// check params
	if err = r.validateParams(parameters); err != nil {
		return nil, err
	}

	blockChainConfigs := localconf.ChainMakerConfig.GetBlockChains()
	chainIds := make([]string, len(blockChainConfigs))
	for i, blockChainConfig := range blockChainConfigs {
		chainIds[i] = blockChainConfig.ChainId
	}

	chainList := &discoveryPb.ChainList{
		ChainIdList: chainIds,
	}
	chainListBytes, err := proto.Marshal(chainList)
	if err != nil {
		errMsg = fmt.Sprintf("marshal chain list failed, %s", err.Error())
		r.log.Errorf(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	return chainListBytes, nil
}

func (r *BlockRuntime) GetChainInfo(txSimContext protocol.TxSimContext, parameters map[string]string) ([]byte, error) {
	var errMsg string
	var err error

	// check params
	if err = r.validateParams(parameters); err != nil {
		return nil, err
	}

	chainId := txSimContext.GetTx().Header.ChainId

	store := txSimContext.GetBlockchainStore()
	if store == nil {
		return nil, errStoreIsNil
	}

	provider, err := txSimContext.GetChainNodesInfoProvider()
	if err != nil {
		return nil, fmt.Errorf("get ChainNodesInfoProvider error: %s", err)
	}

	var block *commonPb.Block
	var nodes []*discoveryPb.Node

	if block, err = r.getBlockByHeight(store, chainId, -1); err != nil {
		return nil, err
	}

	if nodes, err = r.getChainNodeInfo(provider, chainId); err != nil {
		return nil, err
	}

	chainInfo := &discoveryPb.ChainInfo{
		BlockHeight: block.Header.BlockHeight,
		NodeList:    nodes,
	}

	chainInfoBytes, err := proto.Marshal(chainInfo)
	if err != nil {
		errMsg = fmt.Sprintf("marshal chain info failed, %s", err.Error())
		r.log.Errorf(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	return chainInfoBytes, nil
}

func (r *BlockRuntime) GetBlockByHeight(txSimContext protocol.TxSimContext, parameters map[string]string) ([]byte, error) {
	var errMsg string
	var err error

	// check params
	if err = r.validateParams(parameters, paramNameBlockHeight, paramNameWithRWSet); err != nil {
		return nil, err
	}

	chainId := txSimContext.GetTx().Header.ChainId

	store := txSimContext.GetBlockchainStore()
	if store == nil {
		return nil, errStoreIsNil
	}

	var block *commonPb.Block
	var txRWSets []*commonPb.TxRWSet

	if block, err = r.getBlockByHeight(store, chainId, r.height); err != nil {
		return nil, err
	}

	if strings.ToLower(r.withRWSet) == "true" {
		if txRWSets, err = r.getTxRWSetsByBlock(store, chainId, block); err != nil {
			return nil, err
		}
	}

	blockInfo := &commonPb.BlockInfo{
		Block:     block,
		RwsetList: txRWSets,
	}
	blockInfoBytes, err := proto.Marshal(blockInfo)
	if err != nil {
		errMsg = fmt.Sprintf(logTemplateMarshalBlockInfoFailed, err.Error())
		r.log.Errorf(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	return blockInfoBytes, nil

}

func (r *BlockRuntime) GetBlockWithTxRWSetsByHeight(txSimContext protocol.TxSimContext, parameters map[string]string) ([]byte, error) {
	var errMsg string
	var err error

	// check params
	if err = r.validateParams(parameters, paramNameBlockHeight); err != nil {
		return nil, err
	}

	chainId := txSimContext.GetTx().Header.ChainId

	store := txSimContext.GetBlockchainStore()
	if store == nil {
		return nil, errStoreIsNil
	}

	var block *commonPb.Block
	var txRWSets []*commonPb.TxRWSet

	if block, err = r.getBlockByHeight(store, chainId, r.height); err != nil {
		return nil, err
	}

	if txRWSets, err = r.getTxRWSetsByBlock(store, chainId, block); err != nil {
		return nil, err
	}

	blockInfo := &commonPb.BlockInfo{
		Block:     block,
		RwsetList: txRWSets,
	}
	blockInfoBytes, err := proto.Marshal(blockInfo)
	if err != nil {
		errMsg = fmt.Sprintf(logTemplateMarshalBlockInfoFailed, err.Error())
		r.log.Errorf(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	return blockInfoBytes, nil

}

func (r *BlockRuntime) GetBlockByHash(txSimContext protocol.TxSimContext, parameters map[string]string) ([]byte, error) {
	var errMsg string
	var err error

	// check params
	if err = r.validateParams(parameters, paramNameBlockHash, paramNameWithRWSet); err != nil {
		return nil, err
	}

	chainId := txSimContext.GetTx().Header.ChainId

	store := txSimContext.GetBlockchainStore()
	if store == nil {
		return nil, errStoreIsNil
	}

	var block *commonPb.Block
	var txRWSets []*commonPb.TxRWSet

	if block, err = r.getBlockByHash(store, chainId, r.hash); err != nil {
		return nil, err
	}

	if strings.ToLower(r.withRWSet) == "true" {
		if txRWSets, err = r.getTxRWSetsByBlock(store, chainId, block); err != nil {
			return nil, err
		}
	}

	blockInfo := &commonPb.BlockInfo{
		Block:     block,
		RwsetList: txRWSets,
	}
	blockInfoBytes, err := proto.Marshal(blockInfo)
	if err != nil {
		errMsg = fmt.Sprintf(logTemplateMarshalBlockInfoFailed, err.Error())
		r.log.Errorf(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	return blockInfoBytes, nil

}

func (r *BlockRuntime) GetBlockWithTxRWSetsByHash(txSimContext protocol.TxSimContext, parameters map[string]string) ([]byte, error) {
	var errMsg string
	var err error

	// check params
	if err = r.validateParams(parameters, paramNameBlockHash); err != nil {
		return nil, err
	}

	chainId := txSimContext.GetTx().Header.ChainId

	store := txSimContext.GetBlockchainStore()
	if store == nil {
		return nil, errStoreIsNil
	}

	var block *commonPb.Block
	var txRWSets []*commonPb.TxRWSet

	if block, err = r.getBlockByHash(store, chainId, r.hash); err != nil {
		return nil, err
	}

	if txRWSets, err = r.getTxRWSetsByBlock(store, chainId, block); err != nil {
		return nil, err
	}

	blockInfo := &commonPb.BlockInfo{
		Block:     block,
		RwsetList: txRWSets,
	}
	blockInfoBytes, err := proto.Marshal(blockInfo)
	if err != nil {
		errMsg = fmt.Sprintf(logTemplateMarshalBlockInfoFailed, err.Error())
		r.log.Errorf(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	return blockInfoBytes, nil

}

func (r *BlockRuntime) GetBlockByTxId(txSimContext protocol.TxSimContext, parameters map[string]string) ([]byte, error) {
	var errMsg string
	var err error

	// check params
	if err = r.validateParams(parameters, paramNameTxId, paramNameWithRWSet); err != nil {
		return nil, err
	}

	chainId := txSimContext.GetTx().Header.ChainId

	store := txSimContext.GetBlockchainStore()
	if store == nil {
		return nil, errStoreIsNil
	}

	var block *commonPb.Block
	var txRWSets []*commonPb.TxRWSet

	if block, err = r.getBlockByTxId(store, chainId, r.txId); err != nil {
		return nil, err
	}

	if strings.ToLower(r.withRWSet) == "true" {
		if txRWSets, err = r.getTxRWSetsByBlock(store, chainId, block); err != nil {
			return nil, err
		}
	}

	blockInfo := &commonPb.BlockInfo{
		Block:     block,
		RwsetList: txRWSets,
	}
	blockInfoBytes, err := proto.Marshal(blockInfo)
	if err != nil {
		errMsg = fmt.Sprintf(logTemplateMarshalBlockInfoFailed, err.Error())
		r.log.Errorf(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	return blockInfoBytes, nil

}

func (r *BlockRuntime) GetLastConfigBlock(txSimContext protocol.TxSimContext, parameters map[string]string) ([]byte, error) {
	var errMsg string
	var err error

	// check params
	if err = r.validateParams(parameters, paramNameWithRWSet); err != nil {
		return nil, err
	}

	chainId := txSimContext.GetTx().Header.ChainId

	store := txSimContext.GetBlockchainStore()
	if store == nil {
		return nil, errStoreIsNil
	}

	var block *commonPb.Block
	var txRWSets []*commonPb.TxRWSet

	if block, err = r.getLastConfigBlock(store, chainId); err != nil {
		return nil, err
	}

	if strings.ToLower(r.withRWSet) == "true" {
		if txRWSets, err = r.getTxRWSetsByBlock(store, chainId, block); err != nil {
			return nil, err
		}
	}

	blockInfo := &commonPb.BlockInfo{
		Block:     block,
		RwsetList: txRWSets,
	}
	blockInfoBytes, err := proto.Marshal(blockInfo)
	if err != nil {
		errMsg = fmt.Sprintf(logTemplateMarshalBlockInfoFailed, err.Error())
		r.log.Errorf(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	return blockInfoBytes, nil

}

func (r *BlockRuntime) GetLastBlock(txSimContext protocol.TxSimContext, parameters map[string]string) ([]byte, error) {
	var errMsg string
	var err error

	// check params
	if err = r.validateParams(parameters, paramNameWithRWSet); err != nil {
		return nil, err
	}

	chainId := txSimContext.GetTx().Header.ChainId

	store := txSimContext.GetBlockchainStore()
	if store == nil {
		return nil, errStoreIsNil
	}

	var block *commonPb.Block
	var txRWSets []*commonPb.TxRWSet

	if block, err = r.getBlockByHeight(store, chainId, -1); err != nil {
		return nil, err
	}

	if strings.ToLower(r.withRWSet) == "true" {
		if txRWSets, err = r.getTxRWSetsByBlock(store, chainId, block); err != nil {
			return nil, err
		}
	}

	blockInfo := &commonPb.BlockInfo{
		Block:     block,
		RwsetList: txRWSets,
	}
	blockInfoBytes, err := proto.Marshal(blockInfo)
	if err != nil {
		errMsg = fmt.Sprintf(logTemplateMarshalBlockInfoFailed, err.Error())
		r.log.Errorf(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	return blockInfoBytes, nil

}

func (r *BlockRuntime) GetTxByTxId(txSimContext protocol.TxSimContext, parameters map[string]string) ([]byte, error) {
	var errMsg string
	var err error

	// check params
	if err = r.validateParams(parameters, paramNameTxId); err != nil {
		return nil, err
	}

	chainId := txSimContext.GetTx().Header.ChainId

	store := txSimContext.GetBlockchainStore()
	if store == nil {
		return nil, errStoreIsNil
	}

	var tx *commonPb.Transaction
	var block *commonPb.Block

	if tx, err = r.getTxByTxId(store, chainId, r.txId); err != nil {
		return nil, err
	}

	if block, err = r.getBlockByTxId(store, chainId, r.txId); err != nil {
		return nil, err
	}

	transactionInfo := &commonPb.TransactionInfo{
		Transaction: tx,
		BlockHeight: block.Header.BlockHeight,
	}
	transactionInfoBytes, err := proto.Marshal(transactionInfo)
	if err != nil {
		errMsg = fmt.Sprintf("marshal tx failed, %s", err.Error())
		r.log.Errorf(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	return transactionInfoBytes, nil

}

func (r *BlockRuntime) getChainNodeInfo(provider protocol.ChainNodesInfoProvider, chainId string) ([]*discoveryPb.Node, error) {
	nodeInfos, err := provider.GetChainNodesInfo()
	if err != nil {
		r.log.Errorf("get chain node info failed, [chainId:%s], %s", chainId, err.Error())
		return nil, fmt.Errorf("get chain node info failed failed, %s", err)
	}
	nodes := make([]*discoveryPb.Node, len(nodeInfos))
	for i, nodeInfo := range nodeInfos {
		nodes[i] = &discoveryPb.Node{
			NodeId:      nodeInfo.NodeUid,
			NodeAddress: strings.Join(nodeInfo.NodeAddress, ","),
			NodeTlsCert: nodeInfo.NodeTlsCert,
		}
	}
	return nodes, nil
}

func (r *BlockRuntime) getBlockByHeight(store protocol.BlockchainStore, chainId string, height int64) (*commonPb.Block, error) {
	var (
		block *commonPb.Block
		err   error
	)

	if height == -1 {
		block, err = store.GetLastBlock()
	} else {
		block, err = store.GetBlock(height)
	}
	err = r.handleError(block, err, chainId)
	return block, err
}

func (r *BlockRuntime) getBlockByHash(store protocol.BlockchainStore, chainId string, hash string) (*commonPb.Block, error) {
	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		r.log.Errorf("decode hash failed, [hash:%s], %s", hash, err.Error())
		return nil, fmt.Errorf("decode hash failed, %s", err)
	}
	block, err := store.GetBlockByHash(hashBytes)
	err = r.handleError(block, err, chainId)
	return block, err
}

func (r *BlockRuntime) getBlockByTxId(store protocol.BlockchainStore, chainId string, txId string) (*commonPb.Block, error) {
	block, err := store.GetBlockByTx(txId)
	err = r.handleError(block, err, chainId)
	return block, err
}

func (r *BlockRuntime) getLastConfigBlock(store protocol.BlockchainStore, chainId string) (*commonPb.Block, error) {
	block, err := store.GetLastConfigBlock()
	err = r.handleError(block, err, chainId)
	return block, err
}

func (r *BlockRuntime) getTxByTxId(store protocol.BlockchainStore, chainId string, txId string) (*commonPb.Transaction, error) {
	tx, err := store.GetTx(txId)
	err = r.handleError(tx, err, chainId)
	return tx, err
}

func (r *BlockRuntime) getTxRWSetsByBlock(store protocol.BlockchainStore, chainId string, block *commonPb.Block) ([]*commonPb.TxRWSet, error) {
	var txRWSets []*commonPb.TxRWSet
	for _, tx := range block.Txs {
		txRWSet, err := store.GetTxRWSet(tx.Header.TxId)
		if err != nil {
			r.log.Errorf("get txRWset from store failed, [chainId:%s|txId:%s], %s", chainId, tx.Header.TxId, err.Error())
			return nil, fmt.Errorf("get txRWset failed, %s", err)
		}
		txRWSets = append(txRWSets, txRWSet)
	}
	return txRWSets, nil
}

func (r *BlockRuntime) handleError(value interface{}, err error, chainId string) error {
	typeName := strings.ToLower(strings.Split(fmt.Sprintf("%T", value), ".")[1])
	if err != nil {
		r.log.Errorf("get %s from store failed, [chainId:%s], %s", typeName, chainId, err.Error())
		return fmt.Errorf("get %s failed, %s", typeName, err)
	}
	vi := reflect.ValueOf(value)
	if vi.Kind() == reflect.Ptr && vi.IsNil() {
		errMsg := fmt.Sprintf("no such %s, chainId:%s", typeName, chainId)
		r.log.Errorf(errMsg)
		return errors.New(errMsg)
	}

	return nil
}

func (r *BlockRuntime) validateParams(parameters map[string]string, keyNames ...string) error {
	var (
		errMsg string
		err    error
	)
	if len(parameters) != len(keyNames) {
		errMsg = fmt.Sprintf("invalid params len, need [%s]", strings.Join(keyNames, "|"))
		r.log.Error(errMsg)
		return errors.New(errMsg)
	}

	for _, keyName := range keyNames {
		switch keyName {
		case paramNameBlockHeight:
			value, _ := r.getValue(parameters, paramNameBlockHeight)
			r.height, err = strconv.ParseInt(value, 10, 64)
		case paramNameWithRWSet:
			r.withRWSet, err = r.getValue(parameters, paramNameWithRWSet)
		case paramNameBlockHash:
			r.hash, err = r.getValue(parameters, paramNameBlockHash)
		case paramNameTxId:
			r.txId, err = r.getValue(parameters, paramNameTxId)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *BlockRuntime) getValue(parameters map[string]string, key string) (string, error) {
	value, ok := parameters[key]
	if !ok {
		errMsg := fmt.Sprintf("miss params %s", key)
		r.log.Error(errMsg)
		return "", errors.New(errMsg)
	}
	return value, nil
}