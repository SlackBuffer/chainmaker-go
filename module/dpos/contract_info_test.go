/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package dpos

import (
	"fmt"
	"testing"

	"chainmaker.org/chainmaker-go/vm/native"
	"github.com/golang/mock/gomock"
)

func initTestImpl(t *testing.T) (*DPoSImpl, func()) {
	ctrl := gomock.NewController(t)
	mockStore := newMockBlockChainStore(ctrl)
	mockConf := newMockChainConf(ctrl)
	impl := NewDPoSImpl(mockConf, mockStore)
	return impl, func() { ctrl.Finish() }
}

func TestGetStakeAddr(t *testing.T) {
	fmt.Println(native.StakeContractAddr())
}
