/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package txpool

import (
	"fmt"
	"time"

	"chainmaker.org/chainmaker-go/common/msgbus"
	"chainmaker.org/chainmaker-go/localconf"
	"chainmaker.org/chainmaker-go/logger"
	"chainmaker.org/chainmaker-go/protocol"
	batch "chainmaker.org/chainmaker-go/txpool/batchtxpool"
	"chainmaker.org/chainmaker-go/txpool/single"
)

type PoolType int64

const (
	SINGLE PoolType = iota
	BATCH
)

// TxPoolFactory the factory to create the txPool.
type TxPoolFactory struct {
	chainId string
	nodeId  string

	msgBus          msgbus.MessageBus
	netService      protocol.NetService
	blockchainStore protocol.BlockchainStore
	signer          protocol.SigningMember // The identity of the local node
	chainConf       protocol.ChainConf     // chainConfig
	ac              protocol.AccessControlProvider
}

type Option func(f *TxPoolFactory) error

// NewTxPool Create TXPool by applying the additional configuration.
func (f TxPoolFactory) NewTxPool(txPoolType PoolType, opts ...Option) (protocol.TxPool, error) {
	log := logger.GetLogger(logger.MODULE_TXPOOL)
	tf := &TxPoolFactory{}
	if err := tf.Apply(opts...); err != nil {
		log.Errorw("txPoolFactory apply is error", "err", err)
		return nil, err
	}

	if txPoolType == SINGLE {
		return single.NewTxPoolImpl(tf.chainId, tf.blockchainStore, tf.msgBus, tf.chainConf, tf.ac, tf.netService)
	} else if txPoolType == BATCH {
		batchPool, err := batch.NewBatchTxPool(tf.nodeId, tf.chainId)
		if err != nil {
			return nil, err
		}
		if err := batchPool.Apply(batch.WithMsgBus(tf.msgBus),
			batch.WithPoolSize(int(localconf.ChainMakerConfig.TxPoolConfig.MaxTxPoolSize)),
			batch.WithBatchMaxSize(localconf.ChainMakerConfig.TxPoolConfig.BatchMaxSize),
			batch.WithBatchCreateTimeout(time.Duration(localconf.ChainMakerConfig.TxPoolConfig.BatchCreateTimeout)*time.Millisecond),
		); err != nil {
			return nil, err
		}
	}
	return nil, fmt.Errorf("incorrect transaction pool type")
}

// Apply add the extra configurations to the factory
func (f *TxPoolFactory) Apply(opts ...Option) error {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(f); err != nil {
			return err
		}
	}
	return nil
}

// WithMsgBus config the MessageBus
func WithMsgBus(msgBus msgbus.MessageBus) Option {
	return func(f *TxPoolFactory) error {
		f.msgBus = msgBus
		return nil
	}
}

// WithChainId config the chainID in txPool
func WithChainId(chainId string) Option {
	return func(f *TxPoolFactory) error {
		f.chainId = chainId
		return nil
	}
}

// WithNetService config the NetService in txPool
func WithNetService(netService protocol.NetService) Option {
	return func(f *TxPoolFactory) error {
		f.netService = netService
		return nil
	}
}

// WithBlockchainStore config the BlockchainStore service in txPool
func WithBlockchainStore(blockchainStore protocol.BlockchainStore) Option {
	return func(f *TxPoolFactory) error {
		f.blockchainStore = blockchainStore
		return nil
	}
}

// WithSigner config the signer of the local node
func WithSigner(signer protocol.SigningMember) Option {
	return func(f *TxPoolFactory) error {
		f.signer = signer
		return nil
	}
}

// WithChainConf config the chainConf
func WithChainConf(chainConf protocol.ChainConf) Option {
	return func(f *TxPoolFactory) error {
		f.chainConf = chainConf
		return nil
	}
}

// WithAccessControl config the access strategy
func WithAccessControl(ac protocol.AccessControlProvider) Option {
	return func(f *TxPoolFactory) error {
		f.ac = ac
		return nil
	}
}

func WithNodeId(nodeId string) Option {
	return func(f *TxPoolFactory) error {
		f.nodeId = nodeId
		return nil
	}
}