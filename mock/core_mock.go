// Code generated by MockGen. DO NOT EDIT.
// Source: core_interface.go

// Package mock is a generated GoMock package.
package mock

import (
	common "chainmaker.org/chainmaker-go/pb/protogo/common"
	chainedbft "chainmaker.org/chainmaker-go/pb/protogo/consensus/chainedbft"
	txpool "chainmaker.org/chainmaker-go/pb/protogo/txpool"
	protocol "chainmaker.org/chainmaker-go/protocol"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockBlockCommitter is a mock of BlockCommitter interface
type MockBlockCommitter struct {
	ctrl     *gomock.Controller
	recorder *MockBlockCommitterMockRecorder
}

// MockBlockCommitterMockRecorder is the mock recorder for MockBlockCommitter
type MockBlockCommitterMockRecorder struct {
	mock *MockBlockCommitter
}

// NewMockBlockCommitter creates a new mock instance
func NewMockBlockCommitter(ctrl *gomock.Controller) *MockBlockCommitter {
	mock := &MockBlockCommitter{ctrl: ctrl}
	mock.recorder = &MockBlockCommitterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBlockCommitter) EXPECT() *MockBlockCommitterMockRecorder {
	return m.recorder
}

// AddBlock mocks base method
func (m *MockBlockCommitter) AddBlock(blk *common.Block) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBlock", blk)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddBlock indicates an expected call of AddBlock
func (mr *MockBlockCommitterMockRecorder) AddBlock(blk interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBlock", reflect.TypeOf((*MockBlockCommitter)(nil).AddBlock), blk)
}

// MockBlockProposer is a mock of BlockProposer interface
type MockBlockProposer struct {
	ctrl     *gomock.Controller
	recorder *MockBlockProposerMockRecorder
}

// MockBlockProposerMockRecorder is the mock recorder for MockBlockProposer
type MockBlockProposerMockRecorder struct {
	mock *MockBlockProposer
}

// NewMockBlockProposer creates a new mock instance
func NewMockBlockProposer(ctrl *gomock.Controller) *MockBlockProposer {
	mock := &MockBlockProposer{ctrl: ctrl}
	mock.recorder = &MockBlockProposerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBlockProposer) EXPECT() *MockBlockProposerMockRecorder {
	return m.recorder
}

// Start mocks base method
func (m *MockBlockProposer) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
func (mr *MockBlockProposerMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockBlockProposer)(nil).Start))
}

// Stop mocks base method
func (m *MockBlockProposer) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop
func (mr *MockBlockProposerMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockBlockProposer)(nil).Stop))
}

// OnReceiveTxPoolSignal mocks base method
func (m *MockBlockProposer) OnReceiveTxPoolSignal(proposeSignal *txpool.TxPoolSignal) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnReceiveTxPoolSignal", proposeSignal)
}

// OnReceiveTxPoolSignal indicates an expected call of OnReceiveTxPoolSignal
func (mr *MockBlockProposerMockRecorder) OnReceiveTxPoolSignal(proposeSignal interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnReceiveTxPoolSignal", reflect.TypeOf((*MockBlockProposer)(nil).OnReceiveTxPoolSignal), proposeSignal)
}

// OnReceiveProposeStatusChange mocks base method
func (m *MockBlockProposer) OnReceiveProposeStatusChange(proposeStatus bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnReceiveProposeStatusChange", proposeStatus)
}

// OnReceiveProposeStatusChange indicates an expected call of OnReceiveProposeStatusChange
func (mr *MockBlockProposerMockRecorder) OnReceiveProposeStatusChange(proposeStatus interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnReceiveProposeStatusChange", reflect.TypeOf((*MockBlockProposer)(nil).OnReceiveProposeStatusChange), proposeStatus)
}

// OnReceiveChainedBFTProposal mocks base method
func (m *MockBlockProposer) OnReceiveChainedBFTProposal(proposal *chainedbft.BuildProposal) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnReceiveChainedBFTProposal", proposal)
}

// OnReceiveChainedBFTProposal indicates an expected call of OnReceiveChainedBFTProposal
func (mr *MockBlockProposerMockRecorder) OnReceiveChainedBFTProposal(proposal interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnReceiveChainedBFTProposal", reflect.TypeOf((*MockBlockProposer)(nil).OnReceiveChainedBFTProposal), proposal)
}

// MockBlockVerifier is a mock of BlockVerifier interface
type MockBlockVerifier struct {
	ctrl     *gomock.Controller
	recorder *MockBlockVerifierMockRecorder
}

// MockBlockVerifierMockRecorder is the mock recorder for MockBlockVerifier
type MockBlockVerifierMockRecorder struct {
	mock *MockBlockVerifier
}

// NewMockBlockVerifier creates a new mock instance
func NewMockBlockVerifier(ctrl *gomock.Controller) *MockBlockVerifier {
	mock := &MockBlockVerifier{ctrl: ctrl}
	mock.recorder = &MockBlockVerifierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBlockVerifier) EXPECT() *MockBlockVerifierMockRecorder {
	return m.recorder
}

// VerifyBlock mocks base method
func (m *MockBlockVerifier) VerifyBlock(block *common.Block, mode protocol.VerifyMode) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyBlock", block, mode)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyBlock indicates an expected call of VerifyBlock
func (mr *MockBlockVerifierMockRecorder) VerifyBlock(block, mode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyBlock", reflect.TypeOf((*MockBlockVerifier)(nil).VerifyBlock), block, mode)
}