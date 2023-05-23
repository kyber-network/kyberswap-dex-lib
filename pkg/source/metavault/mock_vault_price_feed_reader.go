// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/KyberNetwork/kyberswap-dex-lib/pkg/source/metavault (interfaces: IVaultPriceFeedReader)

// Package metavault is a generated GoMock package.
package metavault

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIVaultPriceFeedReader is a mock of IVaultPriceFeedReader interface.
type MockIVaultPriceFeedReader struct {
	ctrl     *gomock.Controller
	recorder *MockIVaultPriceFeedReaderMockRecorder
}

// MockIVaultPriceFeedReaderMockRecorder is the mock recorder for MockIVaultPriceFeedReader.
type MockIVaultPriceFeedReaderMockRecorder struct {
	mock *MockIVaultPriceFeedReader
}

// NewMockIVaultPriceFeedReader creates a new mock instance.
func NewMockIVaultPriceFeedReader(ctrl *gomock.Controller) *MockIVaultPriceFeedReader {
	mock := &MockIVaultPriceFeedReader{ctrl: ctrl}
	mock.recorder = &MockIVaultPriceFeedReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIVaultPriceFeedReader) EXPECT() *MockIVaultPriceFeedReaderMockRecorder {
	return m.recorder
}

// Read mocks base method.
func (m *MockIVaultPriceFeedReader) Read(arg0 context.Context, arg1 string, arg2 []string) (*VaultPriceFeed, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", arg0, arg1, arg2)
	ret0, _ := ret[0].(*VaultPriceFeed)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockIVaultPriceFeedReaderMockRecorder) Read(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockIVaultPriceFeedReader)(nil).Read), arg0, arg1, arg2)
}
