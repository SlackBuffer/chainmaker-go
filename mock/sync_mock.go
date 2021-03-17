// Code generated by MockGen. DO NOT EDIT.
// Source: sync_interface.go

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockSyncService is a mock of SyncService interface
type MockSyncService struct {
	ctrl     *gomock.Controller
	recorder *MockSyncServiceMockRecorder
}

// MockSyncServiceMockRecorder is the mock recorder for MockSyncService
type MockSyncServiceMockRecorder struct {
	mock *MockSyncService
}

// NewMockSyncService creates a new mock instance
func NewMockSyncService(ctrl *gomock.Controller) *MockSyncService {
	mock := &MockSyncService{ctrl: ctrl}
	mock.recorder = &MockSyncServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSyncService) EXPECT() *MockSyncServiceMockRecorder {
	return m.recorder
}

// Start mocks base method
func (m *MockSyncService) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
func (mr *MockSyncServiceMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockSyncService)(nil).Start))
}

// Stop mocks base method
func (m *MockSyncService) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop
func (mr *MockSyncServiceMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockSyncService)(nil).Stop))
}