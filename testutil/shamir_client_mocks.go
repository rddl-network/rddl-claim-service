// Code generated by MockGen. DO NOT EDIT.
// Source: /home/employee/go/pkg/mod/github.com/rddl-network/shamir-coordinator-service/client@v0.0.4/client.go

// Package testutil is a generated GoMock package.
package testutil

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	service "github.com/rddl-network/shamir-coordinator-service/service"
)

// MockIShamirCoordinatorClient is a mock of IShamirCoordinatorClient interface.
type MockIShamirCoordinatorClient struct {
	ctrl     *gomock.Controller
	recorder *MockIShamirCoordinatorClientMockRecorder
}

// MockIShamirCoordinatorClientMockRecorder is the mock recorder for MockIShamirCoordinatorClient.
type MockIShamirCoordinatorClientMockRecorder struct {
	mock *MockIShamirCoordinatorClient
}

// NewMockIShamirCoordinatorClient creates a new mock instance.
func NewMockIShamirCoordinatorClient(ctrl *gomock.Controller) *MockIShamirCoordinatorClient {
	mock := &MockIShamirCoordinatorClient{ctrl: ctrl}
	mock.recorder = &MockIShamirCoordinatorClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIShamirCoordinatorClient) EXPECT() *MockIShamirCoordinatorClientMockRecorder {
	return m.recorder
}

// GetMnemonics mocks base method.
func (m *MockIShamirCoordinatorClient) GetMnemonics(ctx context.Context) (service.MnemonicsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMnemonics", ctx)
	ret0, _ := ret[0].(service.MnemonicsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMnemonics indicates an expected call of GetMnemonics.
func (mr *MockIShamirCoordinatorClientMockRecorder) GetMnemonics(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMnemonics", reflect.TypeOf((*MockIShamirCoordinatorClient)(nil).GetMnemonics), ctx)
}

// PostMnemonics mocks base method.
func (m *MockIShamirCoordinatorClient) PostMnemonics(ctx context.Context, secret string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostMnemonics", ctx, secret)
	ret0, _ := ret[0].(error)
	return ret0
}

// PostMnemonics indicates an expected call of PostMnemonics.
func (mr *MockIShamirCoordinatorClientMockRecorder) PostMnemonics(ctx, secret interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostMnemonics", reflect.TypeOf((*MockIShamirCoordinatorClient)(nil).PostMnemonics), ctx, secret)
}

// SendTokens mocks base method.
func (m *MockIShamirCoordinatorClient) SendTokens(ctx context.Context, recipient, amount string) (service.SendTokensResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendTokens", ctx, recipient, amount)
	ret0, _ := ret[0].(service.SendTokensResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendTokens indicates an expected call of SendTokens.
func (mr *MockIShamirCoordinatorClientMockRecorder) SendTokens(ctx, recipient, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendTokens", reflect.TypeOf((*MockIShamirCoordinatorClient)(nil).SendTokens), ctx, recipient, amount)
}
