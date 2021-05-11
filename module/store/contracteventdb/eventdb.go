/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/
package contracteventdb

import (
	"chainmaker.org/chainmaker-go/store/serialization"
)

// ContractEventDB provides handle to contract event
// This implementation provides a mysql based data model
type ContractEventDB interface {

	// CommitBlock commits the event in an atomic operation
	CommitBlock(blockInfo *serialization.BlockWithSerializedInfo) error

	//CreateTable create table
	CreateTable(ddl string) error

	// GetLastSavepoint returns the last block height
	GetLastSavepoint() (uint64, error)

	// Close is used to close database
	Close()
}