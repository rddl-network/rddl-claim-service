// Code generated by MockGen. DO NOT EDIT.
// Source: ./service/planetmint_client.go

// Package testutil is a generated GoMock package.
package testutil

import (
	reflect "reflect"

	types "github.com/cosmos/cosmos-sdk/types"
	gomock "github.com/golang/mock/gomock"
)

// MockIPlanetmintClient is a mock of IPlanetmintClient interface.
type MockIPlanetmintClient struct {
	ctrl     *gomock.Controller
	recorder *MockIPlanetmintClientMockRecorder
}

// MockIPlanetmintClientMockRecorder is the mock recorder for MockIPlanetmintClient.
type MockIPlanetmintClientMockRecorder struct {
	mock *MockIPlanetmintClient
}

// NewMockIPlanetmintClient creates a new mock instance.
func NewMockIPlanetmintClient(ctrl *gomock.Controller) *MockIPlanetmintClient {
	mock := &MockIPlanetmintClient{ctrl: ctrl}
	mock.recorder = &MockIPlanetmintClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPlanetmintClient) EXPECT() *MockIPlanetmintClientMockRecorder {
	return m.recorder
}

// SendConfirmation mocks base method.
func (m *MockIPlanetmintClient) SendConfirmation(claimID int, beneficiary string) (types.TxResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendConfirmation", claimID, beneficiary)
	ret0, _ := ret[0].(types.TxResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendConfirmation indicates an expected call of SendConfirmation.
func (mr *MockIPlanetmintClientMockRecorder) SendConfirmation(claimID, beneficiary interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendConfirmation", reflect.TypeOf((*MockIPlanetmintClient)(nil).SendConfirmation), claimID, beneficiary)
}
