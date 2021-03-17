/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

// description: chainmaker-go
//
// @author: xwc1125
// @date: 2020/11/9
package single

import (
	commonPb "chainmaker.org/chainmaker-go/pb/protogo/common"
	"fmt"

	"chainmaker.org/chainmaker-go/common/linkedhashmap"
	"chainmaker.org/chainmaker-go/localconf"
	"chainmaker.org/chainmaker-go/logger"
	"chainmaker.org/chainmaker-go/monitor"
	"chainmaker.org/chainmaker-go/protocol"
	"chainmaker.org/chainmaker-go/utils"

	"github.com/prometheus/client_golang/prometheus"
)

// txList Structure of store transactions in memory
type txList struct {
	log              *logger.CMLogger
	blockchainStore  protocol.BlockchainStore
	metricTxPoolSize *prometheus.GaugeVec

	queue        *linkedhashmap.LinkedHashMap // common transaction will be stored
	pendingCache *linkedhashmap.LinkedHashMap // A place where transactions are stored after Fetch
}

func newTxList(log *logger.CMLogger, pendingCache *linkedhashmap.LinkedHashMap, blockchainStore protocol.BlockchainStore) *txList {
	list := &txList{
		log:             log,
		blockchainStore: blockchainStore,

		queue:        linkedhashmap.NewLinkedHashMap(),
		pendingCache: pendingCache,
	}
	if localconf.ChainMakerConfig.MonitorConfig.Enabled {
		list.metricTxPoolSize = monitor.NewGaugeVec(monitor.SUBSYSTEM_TXPOOL, "metric_tx_pool_size", "tx pool size", "chainId", "poolType")
	}
	return list
}

// Put Add transaction to the txList
func (l *txList) Put(txs []*commonPb.Transaction, source protocol.TxSource, validate txValidateFunc) {
	if len(txs) == 0 {
		return
	}
	for _, tx := range txs {
		if validate == nil || validate(tx, source) == nil {
			l.queue.Add(tx.Header.TxId, tx)
		}
	}

	if localconf.ChainMakerConfig.MonitorConfig.Enabled {
		if utils.IsConfigTx(txs[0]) {
			go l.metricTxPoolSize.WithLabelValues(txs[0].Header.ChainId, "config").Set(float64(l.queue.Size()))
		} else {
			go l.metricTxPoolSize.WithLabelValues(txs[0].Header.ChainId, "normal").Set(float64(l.queue.Size()))
		}
	}
}

// Delete Delete transactions from TXList by the txIds
func (l *txList) Delete(txIds []string) {
	l.log.Debugf("remove txIds", "ids", txIds)
	for _, txId := range txIds {
		l.queue.Remove(txId)
		l.pendingCache.Remove(txId)
	}
}

// Fetch Gets a list of stored transactions
func (l *txList) Fetch(count int, validate func(tx *commonPb.Transaction) error, blockHeight int64) ([]*commonPb.Transaction, []string) {
	queueLen := l.queue.Size()
	if queueLen < count {
		count = queueLen
	}

	var (
		txs      []*commonPb.Transaction
		txIds    []string
		errKeys  []string
		cacheKVs []*valInPendingCache
	)
	defer func() {
		if len(txs) > 0 {
			l.monitor(txs[0], l.Size())
		}
		begin := utils.CurrentTimeMillisSeconds()
		for _, txId := range errKeys {
			l.queue.Remove(txId)
		}
		for _, val := range cacheKVs {
			l.queue.Remove(val.tx.Header.TxId)
			l.pendingCache.Add(val.tx.Header.TxId, val)
		}
		l.log.Debugf("eliminate data, elapse time: %d", utils.CurrentTimeMillisSeconds()-begin)
	}()

	l.log.Debugw("txList Fetch", "count", count, "queueLen", queueLen)
	if queueLen > 0 {
		cacheKVs, txs, txIds, errKeys = l.getTxsFromQueue(count, blockHeight, validate)
		l.log.Debugw("txList Fetch txsNormal", "count", count, "queueLen", queueLen, "txsLen", len(txs), "errKeys", len(errKeys), "cacheKeys", len(cacheKVs))
	}
	return txs, txIds
}

func (l *txList) getTxsFromQueue(count int, blockHeight int64, validate func(tx *commonPb.Transaction) error) (
	cacheKVs []*valInPendingCache, txs []*commonPb.Transaction, txIds []string, errKeys []string) {

	txs = make([]*commonPb.Transaction, 0, count)
	txIds = make([]string, 0, count)
	errKeys = make([]string, 0, count)
	cacheKVs = make([]*valInPendingCache, 0, count)
	node := l.queue.GetLinkList().Front()
	for node != nil && count > 0 {
		txId := node.Value.(string)
		tx := l.queue.Get(txId).(*commonPb.Transaction)
		if validate != nil && validate(tx) != nil {
			errKeys = append(errKeys, txId)
		} else {
			txs = append(txs, tx)
			txIds = append(txIds, txId)
			cacheKVs = append(cacheKVs, &valInPendingCache{
				tx:            tx,
				inBlockHeight: blockHeight,
			})
			if val := l.pendingCache.Get(txId); val != nil {
				l.log.Errorf("tx:%s duplicate to package block, txInPoolHeight: %d", txId, val.(*valInPendingCache).inBlockHeight)
			}
		}
		count--
		node = node.Next()
	}
	return
}

func (l *txList) monitor(tx *commonPb.Transaction, len int) {
	chainId := tx.Header.ChainId
	isConfigTx := utils.IsConfigTx(tx)

	if localconf.ChainMakerConfig.MonitorConfig.Enabled && chainId != "" {
		if isConfigTx {
			go l.metricTxPoolSize.WithLabelValues(chainId, "config").Set(float64(len))
		} else {
			go l.metricTxPoolSize.WithLabelValues(chainId, "normal").Set(float64(len))
		}
	}
}

// Has Determine if the transaction exists in the txList
func (l *txList) Has(txId string, checkPending bool) (exist bool) {
	if checkPending {
		if val := l.pendingCache.Get(txId); val != nil {
			return true
		}
	}
	return l.queue.Get(txId) != nil
}

// Get Retrieve the transaction from txList by the txId
// inBlockHeight: return -1 when the transaction does not exist,
// return 0 when the transaction is in the queue to wait to be generate block,
// return positive integer, indicating that the tx is in an unchained block.
func (l *txList) Get(txId string) (tx *commonPb.Transaction, inBlockHeight int64) {
	if pendingVal := l.pendingCache.Get(txId); pendingVal != nil {
		l.log.Debugw(fmt.Sprintf("txList Get Transaction by txId = %s in pendingCache", txId), "exist", true)
		val := pendingVal.(*valInPendingCache)
		return val.tx, val.inBlockHeight
	}
	if val := l.queue.Get(txId); val != nil {
		l.log.Debugw(fmt.Sprintf("txList Get Transaction by txId = %s in queue", txId), "exist", true)
		return val.(*commonPb.Transaction), 0
	}
	l.log.Debugw(fmt.Sprintf("txList Get Transaction by txId = %s", txId), "exist", false)
	return nil, -1
}

func (l *txList) appendTxsToPendingCache(txs []*commonPb.Transaction, blockHeight int64) {
	for _, tx := range txs {
		l.pendingCache.Add(tx.Header.TxId, &valInPendingCache{
			tx:            tx,
			inBlockHeight: blockHeight,
		})
		l.queue.Remove(tx.Header.TxId)
	}
}

// Size Gets the number of transactions stored in the txList
func (l *txList) Size() int {
	return l.queue.Size()
}