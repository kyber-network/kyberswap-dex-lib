// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/KyberNetwork/kyberswap-dex-lib/pkg/source/madmex (interfaces: IChainlinkFlagsReader)

// Package madmex is a generated GoMock package.
package madmex

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIChainlinkFlagsReader is a mock of IChainlinkFlagsReader interface.
type MockIChainlinkFlagsReader struct {
	ctrl     *gomock.Controller
	recorder *MockIChainlinkFlagsReaderMockRecorder
}

// MockIChainlinkFlagsReaderMockRecorder is the mock recorder for MockIChainlinkFlagsReader.
type MockIChainlinkFlagsReaderMockRecorder struct {
	mock *MockIChainlinkFlagsReader
}

// NewMockIChainlinkFlagsReader creates a new mock instance.
func NewMockIChainlinkFlagsReader(ctrl *gomock.Controller) *MockIChainlinkFlagsReader {
	mock := &MockIChainlinkFlagsReader{ctrl: ctrl}
	mock.recorder = &MockIChainlinkFlagsReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIChainlinkFlagsReader) EXPECT() *MockIChainlinkFlagsReaderMockRecorder {
	return m.recorder
}

// Read mocks base method.
func (m *MockIChainlinkFlagsReader) Read(arg0 context.Context, arg1 string) (*ChainlinkFlags, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", arg0, arg1)
	ret0, _ := ret[0].(*ChainlinkFlags)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockIChainlinkFlagsReaderMockRecorder) Read(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockIChainlinkFlagsReader)(nil).Read), arg0, arg1)
}
