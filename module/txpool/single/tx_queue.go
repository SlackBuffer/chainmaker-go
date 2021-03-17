/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package single

import (
	commonPb "chainmaker.org/chainmaker-go/pb/protogo/common"
	"fmt"

	"chainmaker.org/chainmaker-go/common/linkedhashmap"
	"chainmaker.org/chainmaker-go/logger"
	"chainmaker.org/chainmaker-go/protocol"
	"chainmaker.org/chainmaker-go/utils"
)

type txValidateFunc func(tx *commonPb.Transaction, source protocol.TxSource) error

type txQueue struct {
	log      *logger.CMLogger
	validate txValidateFunc

	pendingCache  *linkedhashmap.LinkedHashMap
	commonTxQueue *txList // common transaction queue
	configTxQueue *txList // config transaction queue
}

func newQueue(blockStore protocol.BlockchainStore, log *logger.CMLogger, validate txValidateFunc) *txQueue {
	pendingCache := linkedhashmap.NewLinkedHashMap()
	queue := txQueue{
		log:           log,
		validate:      validate,
		pendingCache:  pendingCache,
		commonTxQueue: newTxList(log, pendingCache, blockStore),
		configTxQueue: newTxList(log, pendingCache, blockStore),
	}
	return &queue
}

func (queue *txQueue) addTxsToConfigQueue(memTxs *mempoolTxs) {
	queue.configTxQueue.Put(memTxs.txs, memTxs.source, queue.validate)
}

func (queue *txQueue) addTxsToCommonQueue(memTxs *mempoolTxs) {
	queue.commonTxQueue.Put(memTxs.txs, memTxs.source, queue.validate)
}

func (queue *txQueue) deleteTxsInPending(txIds []*commonPb.Transaction) {
	for _, tx := range txIds {
		queue.pendingCache.Remove(tx.Header.TxId)
	}
}

func (queue *txQueue) get(txId string) (tx *commonPb.Transaction, inBlockHeight int64) {
	if tx, inBlockHeight := queue.commonTxQueue.Get(txId); tx != nil {
		return tx, inBlockHeight
	}
	if tx, inBlockHeight := queue.configTxQueue.Get(txId); tx != nil {
		return tx, inBlockHeight
	}
	return nil, -1
}

func (queue *txQueue) configTxsCount() int {
	return queue.configTxQueue.Size()
}

func (queue *txQueue) commonTxsCount() int {
	return queue.commonTxQueue.Size()
}

func (queue *txQueue) deleteConfigTxs(txIds []string) {
	queue.configTxQueue.Delete(txIds)
}

func (queue *txQueue) deleteCommonTxs(txIds []string) {
	queue.commonTxQueue.Delete(txIds)
}

func (queue *txQueue) fetch(expectedCount int, blockHeight int64, validateTxTime func(tx *commonPb.Transaction) error) []*commonPb.Transaction {
	// 1. fetch the config transaction
	if configQueueLen := queue.configTxsCount(); configQueueLen > 0 {
		if txs, txIds := queue.configTxQueue.Fetch(1, validateTxTime, blockHeight); len(txs) > 0 {
			queue.log.Debugw("FetchTxBatch get config txs", "txCount", 1, "configQueueLen", configQueueLen, "txsLen", len(txs), "txIds", txIds)
			return txs
		}
	}

	// 2. fetch the common transaction
	if txQueueLen := queue.commonTxsCount(); txQueueLen > 0 {
		if txs, txIds := queue.commonTxQueue.Fetch(expectedCount, validateTxTime, blockHeight); len(txs) > 0 {
			queue.log.Debugw("FetchTxBatch get common txs", "txCount", expectedCount, "txQueueLen", txQueueLen, "txsLen", len(txs), "txIds", txIds)
			return txs
		}
	}
	return nil
}

func (queue *txQueue) appendTxsToPendingCache(txs []*commonPb.Transaction, blockHeight int64) {
	if utils.IsConfigTx(txs[0]) && len(txs) == 1 {
		queue.configTxQueue.appendTxsToPendingCache(txs, blockHeight)
	}
	if !utils.IsConfigTx(txs[0]) {
		queue.commonTxQueue.appendTxsToPendingCache(txs, blockHeight)
	}
}

func (queue *txQueue) has(tx *commonPb.Transaction, checkPending bool) bool {
	if queue.commonTxQueue.Has(tx.Header.TxId, checkPending) {
		return true
	}
	return queue.configTxQueue.Has(tx.Header.TxId, checkPending)
}

func (queue *txQueue) status() string {
	return fmt.Sprintf("common txs len: %d, config txs len: %d", queue.commonTxQueue.Size(), queue.configTxQueue.Size())
}