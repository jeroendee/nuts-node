// Code generated by MockGen. DO NOT EDIT.
// Source: vcr/verifier/interface.go

// Package verifier is a generated GoMock package.
package verifier

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	ssi "github.com/nuts-foundation/go-did"
	vc "github.com/nuts-foundation/go-did/vc"
	credential "github.com/nuts-foundation/nuts-node/vcr/credential"
)

// MockVerifier is a mock of Verifier interface.
type MockVerifier struct {
	ctrl     *gomock.Controller
	recorder *MockVerifierMockRecorder
}

// MockVerifierMockRecorder is the mock recorder for MockVerifier.
type MockVerifierMockRecorder struct {
	mock *MockVerifier
}

// NewMockVerifier creates a new mock instance.
func NewMockVerifier(ctrl *gomock.Controller) *MockVerifier {
	mock := &MockVerifier{ctrl: ctrl}
	mock.recorder = &MockVerifierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVerifier) EXPECT() *MockVerifierMockRecorder {
	return m.recorder
}

// GetRevocation mocks base method.
func (m *MockVerifier) GetRevocation(id ssi.URI) (*credential.Revocation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRevocation", id)
	ret0, _ := ret[0].(*credential.Revocation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRevocation indicates an expected call of GetRevocation.
func (mr *MockVerifierMockRecorder) GetRevocation(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRevocation", reflect.TypeOf((*MockVerifier)(nil).GetRevocation), id)
}

// IsRevoked mocks base method.
func (m *MockVerifier) IsRevoked(credentialID ssi.URI) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRevoked", credentialID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsRevoked indicates an expected call of IsRevoked.
func (mr *MockVerifierMockRecorder) IsRevoked(credentialID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRevoked", reflect.TypeOf((*MockVerifier)(nil).IsRevoked), credentialID)
}

// RegisterRevocation mocks base method.
func (m *MockVerifier) RegisterRevocation(revocation credential.Revocation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterRevocation", revocation)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterRevocation indicates an expected call of RegisterRevocation.
func (mr *MockVerifierMockRecorder) RegisterRevocation(revocation interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterRevocation", reflect.TypeOf((*MockVerifier)(nil).RegisterRevocation), revocation)
}

// Validate mocks base method.
func (m *MockVerifier) Validate(credentialToVerify vc.VerifiableCredential, at *time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", credentialToVerify, at)
	ret0, _ := ret[0].(error)
	return ret0
}

// Validate indicates an expected call of Validate.
func (mr *MockVerifierMockRecorder) Validate(credentialToVerify, at interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockVerifier)(nil).Validate), credentialToVerify, at)
}

// Verify mocks base method.
func (m *MockVerifier) Verify(credential vc.VerifiableCredential, allowUntrusted, checkSignature bool, validAt *time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify", credential, allowUntrusted, checkSignature, validAt)
	ret0, _ := ret[0].(error)
	return ret0
}

// Verify indicates an expected call of Verify.
func (mr *MockVerifierMockRecorder) Verify(credential, allowUntrusted, checkSignature, validAt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*MockVerifier)(nil).Verify), credential, allowUntrusted, checkSignature, validAt)
}

// VerifyVP mocks base method.
func (m *MockVerifier) VerifyVP(presentation vc.VerifiablePresentation, verifyVCs bool, validAt *time.Time) ([]vc.VerifiableCredential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyVP", presentation, verifyVCs, validAt)
	ret0, _ := ret[0].([]vc.VerifiableCredential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyVP indicates an expected call of VerifyVP.
func (mr *MockVerifierMockRecorder) VerifyVP(presentation, verifyVCs, validAt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyVP", reflect.TypeOf((*MockVerifier)(nil).VerifyVP), presentation, verifyVCs, validAt)
}

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockStore) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStoreMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStore)(nil).Close))
}

// GetRevocations mocks base method.
func (m *MockStore) GetRevocations(id ssi.URI) ([]*credential.Revocation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRevocations", id)
	ret0, _ := ret[0].([]*credential.Revocation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRevocations indicates an expected call of GetRevocations.
func (mr *MockStoreMockRecorder) GetRevocations(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRevocations", reflect.TypeOf((*MockStore)(nil).GetRevocations), id)
}

// StoreRevocation mocks base method.
func (m *MockStore) StoreRevocation(r credential.Revocation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreRevocation", r)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreRevocation indicates an expected call of StoreRevocation.
func (mr *MockStoreMockRecorder) StoreRevocation(r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreRevocation", reflect.TypeOf((*MockStore)(nil).StoreRevocation), r)
}
