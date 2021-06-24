// Code generated by MockGen. DO NOT EDIT.
// Source: access_control_interface.go

// Package mock is a generated GoMock package.
package mock

import (
	pkix "crypto/x509/pkix"
	reflect "reflect"

	x509 "chainmaker.org/chainmaker-go/common/crypto/x509"
	accesscontrol "chainmaker.org/chainmaker-go/pb/protogo/accesscontrol"
	common "chainmaker.org/chainmaker-go/pb/protogo/common"
	config "chainmaker.org/chainmaker-go/pb/protogo/config"
	protocol "chainmaker.org/chainmaker-go/protocol"
	gomock "github.com/golang/mock/gomock"
)

// MockPrincipal is a mock of Principal interface.
type MockPrincipal struct {
	ctrl     *gomock.Controller
	recorder *MockPrincipalMockRecorder
}

// MockPrincipalMockRecorder is the mock recorder for MockPrincipal.
type MockPrincipalMockRecorder struct {
	mock *MockPrincipal
}

// NewMockPrincipal creates a new mock instance.
func NewMockPrincipal(ctrl *gomock.Controller) *MockPrincipal {
	mock := &MockPrincipal{ctrl: ctrl}
	mock.recorder = &MockPrincipalMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPrincipal) EXPECT() *MockPrincipalMockRecorder {
	return m.recorder
}

// GetEndorsement mocks base method.
func (m *MockPrincipal) GetEndorsement() []*common.EndorsementEntry {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEndorsement")
	ret0, _ := ret[0].([]*common.EndorsementEntry)
	return ret0
}

// GetEndorsement indicates an expected call of GetEndorsement.
func (mr *MockPrincipalMockRecorder) GetEndorsement() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEndorsement", reflect.TypeOf((*MockPrincipal)(nil).GetEndorsement))
}

// GetMessage mocks base method.
func (m *MockPrincipal) GetMessage() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessage")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// GetMessage indicates an expected call of GetMessage.
func (mr *MockPrincipalMockRecorder) GetMessage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessage", reflect.TypeOf((*MockPrincipal)(nil).GetMessage))
}

// GetResourceName mocks base method.
func (m *MockPrincipal) GetResourceName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResourceName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetResourceName indicates an expected call of GetResourceName.
func (mr *MockPrincipalMockRecorder) GetResourceName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResourceName", reflect.TypeOf((*MockPrincipal)(nil).GetResourceName))
}

// GetTargetOrgId mocks base method.
func (m *MockPrincipal) GetTargetOrgId() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTargetOrgId")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetTargetOrgId indicates an expected call of GetTargetOrgId.
func (mr *MockPrincipalMockRecorder) GetTargetOrgId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTargetOrgId", reflect.TypeOf((*MockPrincipal)(nil).GetTargetOrgId))
}

// MockAccessControlProvider is a mock of AccessControlProvider interface.
type MockAccessControlProvider struct {
	ctrl     *gomock.Controller
	recorder *MockAccessControlProviderMockRecorder
}

// MockAccessControlProviderMockRecorder is the mock recorder for MockAccessControlProvider.
type MockAccessControlProviderMockRecorder struct {
	mock *MockAccessControlProvider
}

// NewMockAccessControlProvider creates a new mock instance.
func NewMockAccessControlProvider(ctrl *gomock.Controller) *MockAccessControlProvider {
	mock := &MockAccessControlProvider{ctrl: ctrl}
	mock.recorder = &MockAccessControlProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccessControlProvider) EXPECT() *MockAccessControlProviderMockRecorder {
	return m.recorder
}

// CreatePrincipal mocks base method.
func (m *MockAccessControlProvider) CreatePrincipal(resourceName string, endorsements []*common.EndorsementEntry, message []byte) (protocol.Principal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePrincipal", resourceName, endorsements, message)
	ret0, _ := ret[0].(protocol.Principal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePrincipal indicates an expected call of CreatePrincipal.
func (mr *MockAccessControlProviderMockRecorder) CreatePrincipal(resourceName, endorsements, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePrincipal", reflect.TypeOf((*MockAccessControlProvider)(nil).CreatePrincipal), resourceName, endorsements, message)
}

// CreatePrincipalForTargetOrg mocks base method.
func (m *MockAccessControlProvider) CreatePrincipalForTargetOrg(resourceName string, endorsements []*common.EndorsementEntry, message []byte, targetOrgId string) (protocol.Principal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePrincipalForTargetOrg", resourceName, endorsements, message, targetOrgId)
	ret0, _ := ret[0].(protocol.Principal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePrincipalForTargetOrg indicates an expected call of CreatePrincipalForTargetOrg.
func (mr *MockAccessControlProviderMockRecorder) CreatePrincipalForTargetOrg(resourceName, endorsements, message, targetOrgId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePrincipalForTargetOrg", reflect.TypeOf((*MockAccessControlProvider)(nil).CreatePrincipalForTargetOrg), resourceName, endorsements, message, targetOrgId)
}

// DeserializeMember mocks base method.
func (m *MockAccessControlProvider) DeserializeMember(serializedMember []byte) (protocol.Member, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeserializeMember", serializedMember)
	ret0, _ := ret[0].(protocol.Member)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeserializeMember indicates an expected call of DeserializeMember.
func (mr *MockAccessControlProviderMockRecorder) DeserializeMember(serializedMember interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeserializeMember", reflect.TypeOf((*MockAccessControlProvider)(nil).DeserializeMember), serializedMember)
}

// FindPolicyByResourceName mocks base method.
func (m *MockAccessControlProvider) FindPolicyByResourceName(resourceName string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindPolicyByResourceName", resourceName)
	ret0, _ := ret[0].(bool)
	return ret0
}

// FindPolicyByResourceName indicates an expected call of FindPolicyByResourceName.
func (mr *MockAccessControlProviderMockRecorder) FindPolicyByResourceName(resourceName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindPolicyByResourceName", reflect.TypeOf((*MockAccessControlProvider)(nil).FindPolicyByResourceName), resourceName)
}

// GetHashAlg mocks base method.
func (m *MockAccessControlProvider) GetHashAlg() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHashAlg")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetHashAlg indicates an expected call of GetHashAlg.
func (mr *MockAccessControlProviderMockRecorder) GetHashAlg() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHashAlg", reflect.TypeOf((*MockAccessControlProvider)(nil).GetHashAlg))
}

// GetLocalOrgId mocks base method.
func (m *MockAccessControlProvider) GetLocalOrgId() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLocalOrgId")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetLocalOrgId indicates an expected call of GetLocalOrgId.
func (mr *MockAccessControlProviderMockRecorder) GetLocalOrgId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLocalOrgId", reflect.TypeOf((*MockAccessControlProvider)(nil).GetLocalOrgId))
}

// GetLocalSigningMember mocks base method.
func (m *MockAccessControlProvider) GetLocalSigningMember() protocol.SigningMember {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLocalSigningMember")
	ret0, _ := ret[0].(protocol.SigningMember)
	return ret0
}

// GetLocalSigningMember indicates an expected call of GetLocalSigningMember.
func (mr *MockAccessControlProviderMockRecorder) GetLocalSigningMember() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLocalSigningMember", reflect.TypeOf((*MockAccessControlProvider)(nil).GetLocalSigningMember))
}

// GetValidEndorsements mocks base method.
func (m *MockAccessControlProvider) GetValidEndorsements(principal protocol.Principal) ([]*common.EndorsementEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidEndorsements", principal)
	ret0, _ := ret[0].([]*common.EndorsementEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValidEndorsements indicates an expected call of GetValidEndorsements.
func (mr *MockAccessControlProviderMockRecorder) GetValidEndorsements(principal interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidEndorsements", reflect.TypeOf((*MockAccessControlProvider)(nil).GetValidEndorsements), principal)
}

// IsCertRevoked mocks base method.
func (m *MockAccessControlProvider) IsCertRevoked(certChain []*x509.Certificate) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsCertRevoked", certChain)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsCertRevoked indicates an expected call of IsCertRevoked.
func (mr *MockAccessControlProviderMockRecorder) IsCertRevoked(certChain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsCertRevoked", reflect.TypeOf((*MockAccessControlProvider)(nil).IsCertRevoked), certChain)
}

// LookUpResourceNameByTxType mocks base method.
func (m *MockAccessControlProvider) LookUpResourceNameByTxType(txType common.TxType) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LookUpResourceNameByTxType", txType)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LookUpResourceNameByTxType indicates an expected call of LookUpResourceNameByTxType.
func (mr *MockAccessControlProviderMockRecorder) LookUpResourceNameByTxType(txType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LookUpResourceNameByTxType", reflect.TypeOf((*MockAccessControlProvider)(nil).LookUpResourceNameByTxType), txType)
}

// NewMemberFromCertPem mocks base method.
func (m *MockAccessControlProvider) NewMemberFromCertPem(orgId, certPEM string) (protocol.Member, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewMemberFromCertPem", orgId, certPEM)
	ret0, _ := ret[0].(protocol.Member)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewMemberFromCertPem indicates an expected call of NewMemberFromCertPem.
func (mr *MockAccessControlProviderMockRecorder) NewMemberFromCertPem(orgId, certPEM interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewMemberFromCertPem", reflect.TypeOf((*MockAccessControlProvider)(nil).NewMemberFromCertPem), orgId, certPEM)
}

// NewMemberFromProto mocks base method.
func (m *MockAccessControlProvider) NewMemberFromProto(serializedMember *accesscontrol.SerializedMember) (protocol.Member, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewMemberFromProto", serializedMember)
	ret0, _ := ret[0].(protocol.Member)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewMemberFromProto indicates an expected call of NewMemberFromProto.
func (mr *MockAccessControlProviderMockRecorder) NewMemberFromProto(serializedMember interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewMemberFromProto", reflect.TypeOf((*MockAccessControlProvider)(nil).NewMemberFromProto), serializedMember)
}

// NewSigningMember mocks base method.
func (m *MockAccessControlProvider) NewSigningMember(member protocol.Member, privateKeyPem, password string) (protocol.SigningMember, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewSigningMember", member, privateKeyPem, password)
	ret0, _ := ret[0].(protocol.SigningMember)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewSigningMember indicates an expected call of NewSigningMember.
func (mr *MockAccessControlProviderMockRecorder) NewSigningMember(member, privateKeyPem, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewSigningMember", reflect.TypeOf((*MockAccessControlProvider)(nil).NewSigningMember), member, privateKeyPem, password)
}

// NewSigningMemberFromCertFile mocks base method.
func (m *MockAccessControlProvider) NewSigningMemberFromCertFile(orgId, prvKeyFile, password, certFile string) (protocol.SigningMember, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewSigningMemberFromCertFile", orgId, prvKeyFile, password, certFile)
	ret0, _ := ret[0].(protocol.SigningMember)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewSigningMemberFromCertFile indicates an expected call of NewSigningMemberFromCertFile.
func (mr *MockAccessControlProviderMockRecorder) NewSigningMemberFromCertFile(orgId, prvKeyFile, password, certFile interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewSigningMemberFromCertFile", reflect.TypeOf((*MockAccessControlProvider)(nil).NewSigningMemberFromCertFile), orgId, prvKeyFile, password, certFile)
}

// ValidateCRL mocks base method.
func (m *MockAccessControlProvider) ValidateCRL(crl []byte) ([]*pkix.CertificateList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateCRL", crl)
	ret0, _ := ret[0].([]*pkix.CertificateList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateCRL indicates an expected call of ValidateCRL.
func (mr *MockAccessControlProviderMockRecorder) ValidateCRL(crl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateCRL", reflect.TypeOf((*MockAccessControlProvider)(nil).ValidateCRL), crl)
}

// ValidateResourcePolicy mocks base method.
func (m *MockAccessControlProvider) ValidateResourcePolicy(resourcePolicy *config.ResourcePolicy) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateResourcePolicy", resourcePolicy)
	ret0, _ := ret[0].(bool)
	return ret0
}

// ValidateResourcePolicy indicates an expected call of ValidateResourcePolicy.
func (mr *MockAccessControlProviderMockRecorder) ValidateResourcePolicy(resourcePolicy interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateResourcePolicy", reflect.TypeOf((*MockAccessControlProvider)(nil).ValidateResourcePolicy), resourcePolicy)
}

// VerifyPrincipal mocks base method.
func (m *MockAccessControlProvider) VerifyPrincipal(principal protocol.Principal) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyPrincipal", principal)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyPrincipal indicates an expected call of VerifyPrincipal.
func (mr *MockAccessControlProviderMockRecorder) VerifyPrincipal(principal interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyPrincipal", reflect.TypeOf((*MockAccessControlProvider)(nil).VerifyPrincipal), principal)
}

// MockMemberDeserializer is a mock of MemberDeserializer interface.
type MockMemberDeserializer struct {
	ctrl     *gomock.Controller
	recorder *MockMemberDeserializerMockRecorder
}

// MockMemberDeserializerMockRecorder is the mock recorder for MockMemberDeserializer.
type MockMemberDeserializerMockRecorder struct {
	mock *MockMemberDeserializer
}

// NewMockMemberDeserializer creates a new mock instance.
func NewMockMemberDeserializer(ctrl *gomock.Controller) *MockMemberDeserializer {
	mock := &MockMemberDeserializer{ctrl: ctrl}
	mock.recorder = &MockMemberDeserializerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMemberDeserializer) EXPECT() *MockMemberDeserializerMockRecorder {
	return m.recorder
}

// DeserializeMember mocks base method.
func (m *MockMemberDeserializer) DeserializeMember(serializedMember []byte) (protocol.Member, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeserializeMember", serializedMember)
	ret0, _ := ret[0].(protocol.Member)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeserializeMember indicates an expected call of DeserializeMember.
func (mr *MockMemberDeserializerMockRecorder) DeserializeMember(serializedMember interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeserializeMember", reflect.TypeOf((*MockMemberDeserializer)(nil).DeserializeMember), serializedMember)
}

// MockMember is a mock of Member interface.
type MockMember struct {
	ctrl     *gomock.Controller
	recorder *MockMemberMockRecorder
}

// MockMemberMockRecorder is the mock recorder for MockMember.
type MockMemberMockRecorder struct {
	mock *MockMember
}

// NewMockMember creates a new mock instance.
func NewMockMember(ctrl *gomock.Controller) *MockMember {
	mock := &MockMember{ctrl: ctrl}
	mock.recorder = &MockMemberMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMember) EXPECT() *MockMemberMockRecorder {
	return m.recorder
}

// GetCertificate mocks base method.
func (m *MockMember) GetCertificate() (*x509.Certificate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCertificate")
	ret0, _ := ret[0].(*x509.Certificate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCertificate indicates an expected call of GetCertificate.
func (mr *MockMemberMockRecorder) GetCertificate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCertificate", reflect.TypeOf((*MockMember)(nil).GetCertificate))
}

// GetMemberId mocks base method.
func (m *MockMember) GetMemberId() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMemberId")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetMemberId indicates an expected call of GetMemberId.
func (mr *MockMemberMockRecorder) GetMemberId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMemberId", reflect.TypeOf((*MockMember)(nil).GetMemberId))
}

// GetOrgId mocks base method.
func (m *MockMember) GetOrgId() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrgId")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetOrgId indicates an expected call of GetOrgId.
func (mr *MockMemberMockRecorder) GetOrgId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrgId", reflect.TypeOf((*MockMember)(nil).GetOrgId))
}

// GetRole mocks base method.
func (m *MockMember) GetRole() []protocol.Role {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRole")
	ret0, _ := ret[0].([]protocol.Role)
	return ret0
}

// GetRole indicates an expected call of GetRole.
func (mr *MockMemberMockRecorder) GetRole() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRole", reflect.TypeOf((*MockMember)(nil).GetRole))
}

// GetSKI mocks base method.
func (m *MockMember) GetSKI() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSKI")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// GetSKI indicates an expected call of GetSKI.
func (mr *MockMemberMockRecorder) GetSKI() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSKI", reflect.TypeOf((*MockMember)(nil).GetSKI))
}

// GetSerializedMember mocks base method.
func (m *MockMember) GetSerializedMember(isFullCert bool) (*accesscontrol.SerializedMember, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSerializedMember", isFullCert)
	ret0, _ := ret[0].(*accesscontrol.SerializedMember)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSerializedMember indicates an expected call of GetSerializedMember.
func (mr *MockMemberMockRecorder) GetSerializedMember(isFullCert interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSerializedMember", reflect.TypeOf((*MockMember)(nil).GetSerializedMember), isFullCert)
}

// Serialize mocks base method.
func (m *MockMember) Serialize(isFullCert bool) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Serialize", isFullCert)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Serialize indicates an expected call of Serialize.
func (mr *MockMemberMockRecorder) Serialize(isFullCert interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Serialize", reflect.TypeOf((*MockMember)(nil).Serialize), isFullCert)
}

// Verify mocks base method.
func (m *MockMember) Verify(hashType string, msg, sig []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify", hashType, msg, sig)
	ret0, _ := ret[0].(error)
	return ret0
}

// Verify indicates an expected call of Verify.
func (mr *MockMemberMockRecorder) Verify(hashType, msg, sig interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*MockMember)(nil).Verify), hashType, msg, sig)
}

// MockSigningMember is a mock of SigningMember interface.
type MockSigningMember struct {
	ctrl     *gomock.Controller
	recorder *MockSigningMemberMockRecorder
}

// MockSigningMemberMockRecorder is the mock recorder for MockSigningMember.
type MockSigningMemberMockRecorder struct {
	mock *MockSigningMember
}

// NewMockSigningMember creates a new mock instance.
func NewMockSigningMember(ctrl *gomock.Controller) *MockSigningMember {
	mock := &MockSigningMember{ctrl: ctrl}
	mock.recorder = &MockSigningMemberMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSigningMember) EXPECT() *MockSigningMemberMockRecorder {
	return m.recorder
}

// GetCertificate mocks base method.
func (m *MockSigningMember) GetCertificate() (*x509.Certificate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCertificate")
	ret0, _ := ret[0].(*x509.Certificate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCertificate indicates an expected call of GetCertificate.
func (mr *MockSigningMemberMockRecorder) GetCertificate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCertificate", reflect.TypeOf((*MockSigningMember)(nil).GetCertificate))
}

// GetMemberId mocks base method.
func (m *MockSigningMember) GetMemberId() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMemberId")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetMemberId indicates an expected call of GetMemberId.
func (mr *MockSigningMemberMockRecorder) GetMemberId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMemberId", reflect.TypeOf((*MockSigningMember)(nil).GetMemberId))
}

// GetOrgId mocks base method.
func (m *MockSigningMember) GetOrgId() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrgId")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetOrgId indicates an expected call of GetOrgId.
func (mr *MockSigningMemberMockRecorder) GetOrgId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrgId", reflect.TypeOf((*MockSigningMember)(nil).GetOrgId))
}

// GetRole mocks base method.
func (m *MockSigningMember) GetRole() []protocol.Role {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRole")
	ret0, _ := ret[0].([]protocol.Role)
	return ret0
}

// GetRole indicates an expected call of GetRole.
func (mr *MockSigningMemberMockRecorder) GetRole() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRole", reflect.TypeOf((*MockSigningMember)(nil).GetRole))
}

// GetSKI mocks base method.
func (m *MockSigningMember) GetSKI() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSKI")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// GetSKI indicates an expected call of GetSKI.
func (mr *MockSigningMemberMockRecorder) GetSKI() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSKI", reflect.TypeOf((*MockSigningMember)(nil).GetSKI))
}

// GetSerializedMember mocks base method.
func (m *MockSigningMember) GetSerializedMember(isFullCert bool) (*accesscontrol.SerializedMember, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSerializedMember", isFullCert)
	ret0, _ := ret[0].(*accesscontrol.SerializedMember)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSerializedMember indicates an expected call of GetSerializedMember.
func (mr *MockSigningMemberMockRecorder) GetSerializedMember(isFullCert interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSerializedMember", reflect.TypeOf((*MockSigningMember)(nil).GetSerializedMember), isFullCert)
}

// Serialize mocks base method.
func (m *MockSigningMember) Serialize(isFullCert bool) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Serialize", isFullCert)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Serialize indicates an expected call of Serialize.
func (mr *MockSigningMemberMockRecorder) Serialize(isFullCert interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Serialize", reflect.TypeOf((*MockSigningMember)(nil).Serialize), isFullCert)
}

// Sign mocks base method.
func (m *MockSigningMember) Sign(hashType string, msg []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sign", hashType, msg)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sign indicates an expected call of Sign.
func (mr *MockSigningMemberMockRecorder) Sign(hashType, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sign", reflect.TypeOf((*MockSigningMember)(nil).Sign), hashType, msg)
}

// Verify mocks base method.
func (m *MockSigningMember) Verify(hashType string, msg, sig []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify", hashType, msg, sig)
	ret0, _ := ret[0].(error)
	return ret0
}

// Verify indicates an expected call of Verify.
func (mr *MockSigningMemberMockRecorder) Verify(hashType, msg, sig interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*MockSigningMember)(nil).Verify), hashType, msg, sig)
}
