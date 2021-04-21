// Code generated by MockGen. DO NOT EDIT.
// Source: cache_interface.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	common "chainmaker.org/chainmaker-go/pb/protogo/common"
	gomock "github.com/golang/mock/gomock"
)

// MockProposalCache is a mock of ProposalCache interface.
type MockProposalCache struct {
	ctrl     *gomock.Controller
	recorder *MockProposalCacheMockRecorder
}

// MockProposalCacheMockRecorder is the mock recorder for MockProposalCache.
type MockProposalCacheMockRecorder struct {
	mock *MockProposalCache
}

// NewMockProposalCache creates a new mock instance.
func NewMockProposalCache(ctrl *gomock.Controller) *MockProposalCache {
	mock := &MockProposalCache{ctrl: ctrl}
	mock.recorder = &MockProposalCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProposalCache) EXPECT() *MockProposalCacheMockRecorder {
	return m.recorder
}

// ClearProposedBlockAt mocks base method.
func (m *MockProposalCache) ClearProposedBlockAt(height int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ClearProposedBlockAt", height)
}

// ClearProposedBlockAt indicates an expected call of ClearProposedBlockAt.
func (mr *MockProposalCacheMockRecorder) ClearProposedBlockAt(height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearProposedBlockAt", reflect.TypeOf((*MockProposalCache)(nil).ClearProposedBlockAt), height)
}

// GetProposedBlock mocks base method.
func (m *MockProposalCache) GetProposedBlock(b *common.Block) (*common.Block, map[string]*common.TxRWSet, map[string][]*common.ContractEvent) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProposedBlock", b)
	ret0, _ := ret[0].(*common.Block)
	ret1, _ := ret[1].(map[string]*common.TxRWSet)
	ret2, _ := ret[2].(map[string][]*common.ContractEvent)
	return ret0, ret1, ret2
}

// GetProposedBlock indicates an expected call of GetProposedBlock.
func (mr *MockProposalCacheMockRecorder) GetProposedBlock(b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProposedBlock", reflect.TypeOf((*MockProposalCache)(nil).GetProposedBlock), b)
}

// GetProposedBlockByHashAndHeight mocks base method.
func (m *MockProposalCache) GetProposedBlockByHashAndHeight(hash []byte, height int64) (*common.Block, map[string]*common.TxRWSet) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProposedBlockByHashAndHeight", hash, height)
	ret0, _ := ret[0].(*common.Block)
	ret1, _ := ret[1].(map[string]*common.TxRWSet)
	return ret0, ret1
}

// GetProposedBlockByHashAndHeight indicates an expected call of GetProposedBlockByHashAndHeight.
func (mr *MockProposalCacheMockRecorder) GetProposedBlockByHashAndHeight(hash, height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProposedBlockByHashAndHeight", reflect.TypeOf((*MockProposalCache)(nil).GetProposedBlockByHashAndHeight), hash, height)
}

// GetProposedBlocksAt mocks base method.
func (m *MockProposalCache) GetProposedBlocksAt(height int64) []*common.Block {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProposedBlocksAt", height)
	ret0, _ := ret[0].([]*common.Block)
	return ret0
}

// GetProposedBlocksAt indicates an expected call of GetProposedBlocksAt.
func (mr *MockProposalCacheMockRecorder) GetProposedBlocksAt(height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProposedBlocksAt", reflect.TypeOf((*MockProposalCache)(nil).GetProposedBlocksAt), height)
}

// GetSelfProposedBlockAt mocks base method.
func (m *MockProposalCache) GetSelfProposedBlockAt(height int64) *common.Block {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSelfProposedBlockAt", height)
	ret0, _ := ret[0].(*common.Block)
	return ret0
}

// GetSelfProposedBlockAt indicates an expected call of GetSelfProposedBlockAt.
func (mr *MockProposalCacheMockRecorder) GetSelfProposedBlockAt(height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSelfProposedBlockAt", reflect.TypeOf((*MockProposalCache)(nil).GetSelfProposedBlockAt), height)
}

// HasProposedBlockAt mocks base method.
func (m *MockProposalCache) HasProposedBlockAt(height int64) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasProposedBlockAt", height)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasProposedBlockAt indicates an expected call of HasProposedBlockAt.
func (mr *MockProposalCacheMockRecorder) HasProposedBlockAt(height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasProposedBlockAt", reflect.TypeOf((*MockProposalCache)(nil).HasProposedBlockAt), height)
}

// IsProposedAt mocks base method.
func (m *MockProposalCache) IsProposedAt(height int64) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsProposedAt", height)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsProposedAt indicates an expected call of IsProposedAt.
func (mr *MockProposalCacheMockRecorder) IsProposedAt(height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsProposedAt", reflect.TypeOf((*MockProposalCache)(nil).IsProposedAt), height)
}

// KeepProposedBlock mocks base method.
func (m *MockProposalCache) KeepProposedBlock(hash []byte, height int64) []*common.Block {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KeepProposedBlock", hash, height)
	ret0, _ := ret[0].([]*common.Block)
	return ret0
}

// KeepProposedBlock indicates an expected call of KeepProposedBlock.
func (mr *MockProposalCacheMockRecorder) KeepProposedBlock(hash, height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KeepProposedBlock", reflect.TypeOf((*MockProposalCache)(nil).KeepProposedBlock), hash, height)
}

// ResetProposedAt mocks base method.
func (m *MockProposalCache) ResetProposedAt(height int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ResetProposedAt", height)
}

// ResetProposedAt indicates an expected call of ResetProposedAt.
func (mr *MockProposalCacheMockRecorder) ResetProposedAt(height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetProposedAt", reflect.TypeOf((*MockProposalCache)(nil).ResetProposedAt), height)
}

// SetProposedAt mocks base method.
func (m *MockProposalCache) SetProposedAt(height int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetProposedAt", height)
}

// SetProposedAt indicates an expected call of SetProposedAt.
func (mr *MockProposalCacheMockRecorder) SetProposedAt(height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetProposedAt", reflect.TypeOf((*MockProposalCache)(nil).SetProposedAt), height)
}

// SetProposedBlock mocks base method.
func (m *MockProposalCache) SetProposedBlock(b *common.Block, rwSetMap map[string]*common.TxRWSet, contractEventMap map[string][]*common.ContractEvent, selfProposed bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetProposedBlock", b, rwSetMap, contractEventMap, selfProposed)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetProposedBlock indicates an expected call of SetProposedBlock.
func (mr *MockProposalCacheMockRecorder) SetProposedBlock(b, rwSetMap, contractEventMap, selfProposed interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetProposedBlock", reflect.TypeOf((*MockProposalCache)(nil).SetProposedBlock), b, rwSetMap, contractEventMap, selfProposed)
}

// MockLedgerCache is a mock of LedgerCache interface.
type MockLedgerCache struct {
	ctrl     *gomock.Controller
	recorder *MockLedgerCacheMockRecorder
}

// MockLedgerCacheMockRecorder is the mock recorder for MockLedgerCache.
type MockLedgerCacheMockRecorder struct {
	mock *MockLedgerCache
}

// NewMockLedgerCache creates a new mock instance.
func NewMockLedgerCache(ctrl *gomock.Controller) *MockLedgerCache {
	mock := &MockLedgerCache{ctrl: ctrl}
	mock.recorder = &MockLedgerCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLedgerCache) EXPECT() *MockLedgerCacheMockRecorder {
	return m.recorder
}

// CurrentHeight mocks base method.
func (m *MockLedgerCache) CurrentHeight() (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentHeight")
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CurrentHeight indicates an expected call of CurrentHeight.
func (mr *MockLedgerCacheMockRecorder) CurrentHeight() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentHeight", reflect.TypeOf((*MockLedgerCache)(nil).CurrentHeight))
}

// GetLastCommittedBlock mocks base method.
func (m *MockLedgerCache) GetLastCommittedBlock() *common.Block {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastCommittedBlock")
	ret0, _ := ret[0].(*common.Block)
	return ret0
}

// GetLastCommittedBlock indicates an expected call of GetLastCommittedBlock.
func (mr *MockLedgerCacheMockRecorder) GetLastCommittedBlock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastCommittedBlock", reflect.TypeOf((*MockLedgerCache)(nil).GetLastCommittedBlock))
}

// SetLastCommittedBlock mocks base method.
func (m *MockLedgerCache) SetLastCommittedBlock(b *common.Block) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetLastCommittedBlock", b)
}

// SetLastCommittedBlock indicates an expected call of SetLastCommittedBlock.
func (mr *MockLedgerCacheMockRecorder) SetLastCommittedBlock(b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLastCommittedBlock", reflect.TypeOf((*MockLedgerCache)(nil).SetLastCommittedBlock), b)
}
