// Code generated by MockGen. DO NOT EDIT.
// Source: service/shamir_client.go

// Package testutil is a generated GoMock package.
package testutil

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIShamirClient is a mock of IShamirClient interface.
type MockIShamirClient struct {
	ctrl     *gomock.Controller
	recorder *MockIShamirClientMockRecorder
}

// MockIShamirClientMockRecorder is the mock recorder for MockIShamirClient.
type MockIShamirClientMockRecorder struct {
	mock *MockIShamirClient
}

// NewMockIShamirClient creates a new mock instance.
func NewMockIShamirClient(ctrl *gomock.Controller) *MockIShamirClient {
	mock := &MockIShamirClient{ctrl: ctrl}
	mock.recorder = &MockIShamirClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIShamirClient) EXPECT() *MockIShamirClientMockRecorder {
	return m.recorder
}

// IssueTransaction mocks base method.
func (m *MockIShamirClient) IssueTransaction(amount, address string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IssueTransaction", amount, address)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IssueTransaction indicates an expected call of IssueTransaction.
func (mr *MockIShamirClientMockRecorder) IssueTransaction(amount, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IssueTransaction", reflect.TypeOf((*MockIShamirClient)(nil).IssueTransaction), amount, address)
}
