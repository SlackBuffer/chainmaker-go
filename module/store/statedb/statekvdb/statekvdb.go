/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package statekvdb

import (
	logImpl "chainmaker.org/chainmaker-go/logger"
	storePb "chainmaker.org/chainmaker-go/pb/protogo/store"
	"chainmaker.org/chainmaker-go/protocol"
	"chainmaker.org/chainmaker-go/store/cache"
	"chainmaker.org/chainmaker-go/store/dbprovider"
	"chainmaker.org/chainmaker-go/store/types"
	"chainmaker.org/chainmaker-go/utils"
	"encoding/binary"
	"fmt"
)

const (
	stateDBName            = ""
	contractStoreSeparator = '#'
	stateDBSavepointKey    = "stateDBSavePointKey"
)

// StateKvDB provider a implementation of `statedb.StateDB`
// This implementation provides a key-value based data model
type StateKvDB struct {
	DbProvider dbprovider.Provider
	Cache      *cache.StoreCacheMgr
	Logger     *logImpl.CMLogger
}

// CommitBlock commits the state in an atomic operation
func (s *StateKvDB) CommitBlock(blockWithRWSet *storePb.BlockWithRWSet) error {
	batch := types.NewUpdateBatch()
	// 1. last block height
	block := blockWithRWSet.Block
	lastBlockNumBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(lastBlockNumBytes, uint64(block.Header.BlockHeight))
	batch.Put([]byte(stateDBSavepointKey), lastBlockNumBytes)

	txRWSets := blockWithRWSet.TxRWSets

	//try to add consensusArgs
	consensusArgs, err := utils.GetConsensusArgsFromBlock(block)
	if err == nil && consensusArgs.ConsensusData != nil {
		s.Logger.Debugf("add consensusArgs ConsensusData!")
		txRWSets = append(txRWSets, consensusArgs.ConsensusData)
	}

	for _, txRWSet := range txRWSets {
		for _, txWrite := range txRWSet.TxWrites {
			// 5. state: contractID + stateKey
			txWriteKey := constructStateKey(txWrite.ContractName, txWrite.Key)
			if txWrite.Value == nil {
				batch.Delete(txWriteKey)
			} else {
				batch.Put(txWriteKey, txWrite.Value)
			}
		}
	}

	err = s.writeBatch(block.Header.BlockHeight, batch)
	if err != nil {
		return err
	}
	s.Logger.Debugf("chain[%s]: commit state block[%d]",
		block.Header.ChainId, block.Header.BlockHeight)
	return nil
}

// ReadObject returns the state value for given contract name and key, or returns nil if none exists.
func (s *StateKvDB) ReadObject(contractName string, key []byte) ([]byte, error) {
	objectKey := constructStateKey(contractName, key)
	return s.get(objectKey)
}

// SelectObject returns an iterator that contains all the key-values between given key ranges.
// startKey is included in the results and limit is excluded.
func (s *StateKvDB) SelectObject(contractName string, startKey []byte, limit []byte) protocol.Iterator {
	objectStartKey := constructStateKey(contractName, startKey)
	objectLimitKey := constructStateKey(contractName, limit)
	//todo combine cache and database
	s.Cache.LockForFlush()
	defer s.Cache.UnLockFlush()
	//logger.Debugf("start[%s], limit[%s]", objectStartKey, objectLimitKey)
	return s.getDBHandle().NewIteratorWithRange(objectStartKey, objectLimitKey)
}

// GetLastSavepoint returns the last block height
func (b *StateKvDB) GetLastSavepoint() (uint64, error) {
	bytes, err := b.get([]byte(stateDBSavepointKey))
	if err != nil {
		return 0, err
	} else if bytes == nil {
		return 0, nil
	}
	num := binary.BigEndian.Uint64(bytes)
	return num, nil
}

// Close is used to close database
func (s *StateKvDB) Close() {
	s.DbProvider.Close()
}

func (s *StateKvDB) writeBatch(blockHeight int64, batch protocol.StoreBatcher) error {
	//update cache
	s.Cache.AddBlock(blockHeight, batch)
	go func() {
		err := s.getDBHandle().WriteBatch(batch, false)
		if err != nil {
			panic(fmt.Sprintf("Error writting leveldb: %s", err))
		}
		//db committed, clean cache
		s.Cache.DelBlock(blockHeight)
	}()
	return nil
}

func (s *StateKvDB) get(key []byte) ([]byte, error) {
	//get from cache
	value, exist := s.Cache.Get(string(key))
	if exist {
		return value, nil
	}
	//get from database
	return s.getDBHandle().Get(key)
}

func (s *StateKvDB) has(key []byte) (bool, error) {
	//check has from cache
	isDelete, exist := s.Cache.Has(string(key))
	if exist {
		return !isDelete, nil
	}
	return s.getDBHandle().Has(key)
}

func (s *StateKvDB) getDBHandle() protocol.DBHandle {
	return s.DbProvider.GetDBHandle(stateDBName)
}

func constructStateKey(contractName string, key []byte) []byte {
	return append(append([]byte(contractName), contractStoreSeparator), key...)
}