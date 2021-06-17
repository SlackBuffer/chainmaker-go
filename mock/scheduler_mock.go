// Code generated by MockGen. DO NOT EDIT.
// Source: scheduler_interface.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	accesscontrol "chainmaker.org/chainmaker-go/pb/protogo/accesscontrol"
	common "chainmaker.org/chainmaker-go/pb/protogo/common"
	protocol "chainmaker.org/chainmaker-go/protocol"
	gomock "github.com/golang/mock/gomock"
)

// MockTxScheduler is a mock of TxScheduler interface.
type MockTxScheduler struct {
	ctrl     *gomock.Controller
	recorder *MockTxSchedulerMockRecorder
}

// MockTxSchedulerMockRecorder is the mock recorder for MockTxScheduler.
type MockTxSchedulerMockRecorder struct {
	mock *MockTxScheduler
}

// NewMockTxScheduler creates a new mock instance.
func NewMockTxScheduler(ctrl *gomock.Controller) *MockTxScheduler {
	mock := &MockTxScheduler{ctrl: ctrl}
	mock.recorder = &MockTxSchedulerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTxScheduler) EXPECT() *MockTxSchedulerMockRecorder {
	return m.recorder
}

// Halt mocks base method.
func (m *MockTxScheduler) Halt() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Halt")
}

// Halt indicates an expected call of Halt.
func (mr *MockTxSchedulerMockRecorder) Halt() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Halt", reflect.TypeOf((*MockTxScheduler)(nil).Halt))
}

// Schedule mocks base method.
func (m *MockTxScheduler) Schedule(block *common.Block, txBatch []*common.Transaction, snapshot protocol.Snapshot) (map[string]*common.TxRWSet, map[string][]*common.ContractEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Schedule", block, txBatch, snapshot)
	ret0, _ := ret[0].(map[string]*common.TxRWSet)
	ret1, _ := ret[1].(map[string][]*common.ContractEvent)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Schedule indicates an expected call of Schedule.
func (mr *MockTxSchedulerMockRecorder) Schedule(block, txBatch, snapshot interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Schedule", reflect.TypeOf((*MockTxScheduler)(nil).Schedule), block, txBatch, snapshot)
}

// SimulateWithDag mocks base method.
func (m *MockTxScheduler) SimulateWithDag(block *common.Block, snapshot protocol.Snapshot) (map[string]*common.TxRWSet, map[string]*common.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SimulateWithDag", block, snapshot)
	ret0, _ := ret[0].(map[string]*common.TxRWSet)
	ret1, _ := ret[1].(map[string]*common.Result)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SimulateWithDag indicates an expected call of SimulateWithDag.
func (mr *MockTxSchedulerMockRecorder) SimulateWithDag(block, snapshot interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SimulateWithDag", reflect.TypeOf((*MockTxScheduler)(nil).SimulateWithDag), block, snapshot)
}

// MockTxSimContext is a mock of TxSimContext interface.
type MockTxSimContext struct {
	ctrl     *gomock.Controller
	recorder *MockTxSimContextMockRecorder
}

// MockTxSimContextMockRecorder is the mock recorder for MockTxSimContext.
type MockTxSimContextMockRecorder struct {
	mock *MockTxSimContext
}

// NewMockTxSimContext creates a new mock instance.
func NewMockTxSimContext(ctrl *gomock.Controller) *MockTxSimContext {
	mock := &MockTxSimContext{ctrl: ctrl}
	mock.recorder = &MockTxSimContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTxSimContext) EXPECT() *MockTxSimContextMockRecorder {
	return m.recorder
}

// CallContract mocks base method.
func (m *MockTxSimContext) CallContract(contractId *common.ContractId, method string, byteCode []byte, parameter map[string]string, gasUsed uint64, refTxType common.TxType) (*common.ContractResult, common.TxStatusCode) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CallContract", contractId, method, byteCode, parameter, gasUsed, refTxType)
	ret0, _ := ret[0].(*common.ContractResult)
	ret1, _ := ret[1].(common.TxStatusCode)
	return ret0, ret1
}

// CallContract indicates an expected call of CallContract.
func (mr *MockTxSimContextMockRecorder) CallContract(contractId, method, byteCode, parameter, gasUsed, refTxType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CallContract", reflect.TypeOf((*MockTxSimContext)(nil).CallContract), contractId, method, byteCode, parameter, gasUsed, refTxType)
}

// Del mocks base method.
func (m *MockTxSimContext) Del(name string, key []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Del", name, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Del indicates an expected call of Del.
func (mr *MockTxSimContextMockRecorder) Del(name, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Del", reflect.TypeOf((*MockTxSimContext)(nil).Del), name, key)
}

// Get mocks base method.
func (m *MockTxSimContext) Get(name string, key []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", name, key)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockTxSimContextMockRecorder) Get(name, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockTxSimContext)(nil).Get), name, key)
}

// GetAccessControl mocks base method.
func (m *MockTxSimContext) GetAccessControl() (protocol.AccessControlProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccessControl")
	ret0, _ := ret[0].(protocol.AccessControlProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccessControl indicates an expected call of GetAccessControl.
func (mr *MockTxSimContextMockRecorder) GetAccessControl() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccessControl", reflect.TypeOf((*MockTxSimContext)(nil).GetAccessControl))
}

// GetBlockHeight mocks base method.
func (m *MockTxSimContext) GetBlockHeight() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockHeight")
	ret0, _ := ret[0].(int64)
	return ret0
}

// GetBlockHeight indicates an expected call of GetBlockHeight.
func (mr *MockTxSimContextMockRecorder) GetBlockHeight() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockHeight", reflect.TypeOf((*MockTxSimContext)(nil).GetBlockHeight))
}

// GetBlockProposer mocks base method.
func (m *MockTxSimContext) GetBlockProposer() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockProposer")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// GetBlockProposer indicates an expected call of GetBlockProposer.
func (mr *MockTxSimContextMockRecorder) GetBlockProposer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockProposer", reflect.TypeOf((*MockTxSimContext)(nil).GetBlockProposer))
}

// GetBlockchainStore mocks base method.
func (m *MockTxSimContext) GetBlockchainStore() protocol.BlockchainStore {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockchainStore")
	ret0, _ := ret[0].(protocol.BlockchainStore)
	return ret0
}

// GetBlockchainStore indicates an expected call of GetBlockchainStore.
func (mr *MockTxSimContextMockRecorder) GetBlockchainStore() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockchainStore", reflect.TypeOf((*MockTxSimContext)(nil).GetBlockchainStore))
}

// GetChainNodesInfoProvider mocks base method.
func (m *MockTxSimContext) GetChainNodesInfoProvider() (protocol.ChainNodesInfoProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChainNodesInfoProvider")
	ret0, _ := ret[0].(protocol.ChainNodesInfoProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChainNodesInfoProvider indicates an expected call of GetChainNodesInfoProvider.
func (mr *MockTxSimContextMockRecorder) GetChainNodesInfoProvider() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChainNodesInfoProvider", reflect.TypeOf((*MockTxSimContext)(nil).GetChainNodesInfoProvider))
}

// GetCreator mocks base method.
func (m *MockTxSimContext) GetCreator(namespace string) *accesscontrol.SerializedMember {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCreator", namespace)
	ret0, _ := ret[0].(*accesscontrol.SerializedMember)
	return ret0
}

// GetCreator indicates an expected call of GetCreator.
func (mr *MockTxSimContextMockRecorder) GetCreator(namespace interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCreator", reflect.TypeOf((*MockTxSimContext)(nil).GetCreator), namespace)
}

// GetCurrentResult mocks base method.
func (m *MockTxSimContext) GetCurrentResult() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentResult")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// GetCurrentResult indicates an expected call of GetCurrentResult.
func (mr *MockTxSimContextMockRecorder) GetCurrentResult() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentResult", reflect.TypeOf((*MockTxSimContext)(nil).GetCurrentResult))
}

// GetDepth mocks base method.
func (m *MockTxSimContext) GetDepth() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDepth")
	ret0, _ := ret[0].(int)
	return ret0
}

// GetDepth indicates an expected call of GetDepth.
func (mr *MockTxSimContextMockRecorder) GetDepth() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDepth", reflect.TypeOf((*MockTxSimContext)(nil).GetDepth))
}

// GetSender mocks base method.
func (m *MockTxSimContext) GetSender() *accesscontrol.SerializedMember {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSender")
	ret0, _ := ret[0].(*accesscontrol.SerializedMember)
	return ret0
}

// GetSender indicates an expected call of GetSender.
func (mr *MockTxSimContextMockRecorder) GetSender() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSender", reflect.TypeOf((*MockTxSimContext)(nil).GetSender))
}

// GetStateKvHandle mocks base method.
func (m *MockTxSimContext) GetStateKvHandle(arg0 int32) (protocol.StateIterator, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStateKvHandle", arg0)
	ret0, _ := ret[0].(protocol.StateIterator)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetStateKvHandle indicates an expected call of GetStateKvHandle.
func (mr *MockTxSimContextMockRecorder) GetStateKvHandle(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStateKvHandle", reflect.TypeOf((*MockTxSimContext)(nil).GetStateKvHandle), arg0)
}

// GetStateSqlHandle mocks base method.
func (m *MockTxSimContext) GetStateSqlHandle(arg0 int32) (protocol.SqlRows, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStateSqlHandle", arg0)
	ret0, _ := ret[0].(protocol.SqlRows)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetStateSqlHandle indicates an expected call of GetStateSqlHandle.
func (mr *MockTxSimContextMockRecorder) GetStateSqlHandle(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStateSqlHandle", reflect.TypeOf((*MockTxSimContext)(nil).GetStateSqlHandle), arg0)
}

// GetTx mocks base method.
func (m *MockTxSimContext) GetTx() *common.Transaction {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTx")
	ret0, _ := ret[0].(*common.Transaction)
	return ret0
}

// GetTx indicates an expected call of GetTx.
func (mr *MockTxSimContextMockRecorder) GetTx() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTx", reflect.TypeOf((*MockTxSimContext)(nil).GetTx))
}

// GetTxExecSeq mocks base method.
func (m *MockTxSimContext) GetTxExecSeq() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxExecSeq")
	ret0, _ := ret[0].(int)
	return ret0
}

// GetTxExecSeq indicates an expected call of GetTxExecSeq.
func (mr *MockTxSimContextMockRecorder) GetTxExecSeq() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxExecSeq", reflect.TypeOf((*MockTxSimContext)(nil).GetTxExecSeq))
}

// GetTxRWSet mocks base method.
func (m *MockTxSimContext) GetTxRWSet(runVmSuccess bool) *common.TxRWSet {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxRWSet")
	ret0, _ := ret[0].(*common.TxRWSet)
	return ret0
}

// GetTxRWSet indicates an expected call of GetTxRWSet.
func (mr *MockTxSimContextMockRecorder) GetTxRWSet(runVmSuccess bool) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxRWSet", reflect.TypeOf((*MockTxSimContext)(nil).GetTxRWSet))
}

// GetTxResult mocks base method.
func (m *MockTxSimContext) GetTxResult() *common.Result {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxResult")
	ret0, _ := ret[0].(*common.Result)
	return ret0
}

// GetTxResult indicates an expected call of GetTxResult.
func (mr *MockTxSimContextMockRecorder) GetTxResult() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxResult", reflect.TypeOf((*MockTxSimContext)(nil).GetTxResult))
}

// Put mocks base method.
func (m *MockTxSimContext) Put(name string, key, value []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", name, key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockTxSimContextMockRecorder) Put(name, key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockTxSimContext)(nil).Put), name, key, value)
}

// PutRecord mocks base method.
func (m *MockTxSimContext) PutRecord(contractName string, value []byte, sqlType protocol.SqlType) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PutRecord", contractName, value)
}

// PutRecord indicates an expected call of PutRecord.
func (mr *MockTxSimContextMockRecorder) PutRecord(contractName, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutRecord", reflect.TypeOf((*MockTxSimContext)(nil).PutRecord), contractName, value)
}

// Select mocks base method.
func (m *MockTxSimContext) Select(name string, startKey, limit []byte) (protocol.StateIterator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Select", name, startKey, limit)
	ret0, _ := ret[0].(protocol.StateIterator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Select indicates an expected call of Select.
func (mr *MockTxSimContextMockRecorder) Select(name, startKey, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Select", reflect.TypeOf((*MockTxSimContext)(nil).Select), name, startKey, limit)
}

// SetStateKvHandle mocks base method.
func (m *MockTxSimContext) SetStateKvHandle(arg0 int32, arg1 protocol.StateIterator) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetStateKvHandle", arg0, arg1)
}

// SetStateKvHandle indicates an expected call of SetStateKvHandle.
func (mr *MockTxSimContextMockRecorder) SetStateKvHandle(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStateKvHandle", reflect.TypeOf((*MockTxSimContext)(nil).SetStateKvHandle), arg0, arg1)
}

// SetStateSqlHandle mocks base method.
func (m *MockTxSimContext) SetStateSqlHandle(arg0 int32, arg1 protocol.SqlRows) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetStateSqlHandle", arg0, arg1)
}

// SetStateSqlHandle indicates an expected call of SetStateSqlHandle.
func (mr *MockTxSimContextMockRecorder) SetStateSqlHandle(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStateSqlHandle", reflect.TypeOf((*MockTxSimContext)(nil).SetStateSqlHandle), arg0, arg1)
}

// SetTxExecSeq mocks base method.
func (m *MockTxSimContext) SetTxExecSeq(arg0 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTxExecSeq", arg0)
}

// SetTxExecSeq indicates an expected call of SetTxExecSeq.
func (mr *MockTxSimContextMockRecorder) SetTxExecSeq(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTxExecSeq", reflect.TypeOf((*MockTxSimContext)(nil).SetTxExecSeq), arg0)
}

// SetTxResult mocks base method.
func (m *MockTxSimContext) SetTxResult(arg0 *common.Result) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTxResult", arg0)
}

// SetTxResult indicates an expected call of SetTxResult.
func (mr *MockTxSimContextMockRecorder) SetTxResult(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTxResult", reflect.TypeOf((*MockTxSimContext)(nil).SetTxResult), arg0)
}
