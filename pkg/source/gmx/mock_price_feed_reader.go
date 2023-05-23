// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/KyberNetwork/kyberswap-dex-lib/pkg/source/gmx (interfaces: IPriceFeedReader)

// Package gmx is a generated GoMock package.
package gmx

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIPriceFeedReader is a mock of IPriceFeedReader interface.
type MockIPriceFeedReader struct {
	ctrl     *gomock.Controller
	recorder *MockIPriceFeedReaderMockRecorder
}

// MockIPriceFeedReaderMockRecorder is the mock recorder for MockIPriceFeedReader.
type MockIPriceFeedReaderMockRecorder struct {
	mock *MockIPriceFeedReader
}

// NewMockIPriceFeedReader creates a new mock instance.
func NewMockIPriceFeedReader(ctrl *gomock.Controller) *MockIPriceFeedReader {
	mock := &MockIPriceFeedReader{ctrl: ctrl}
	mock.recorder = &MockIPriceFeedReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPriceFeedReader) EXPECT() *MockIPriceFeedReaderMockRecorder {
	return m.recorder
}

// Read mocks base method.
func (m *MockIPriceFeedReader) Read(arg0 context.Context, arg1 string, arg2 int) (*PriceFeed, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", arg0, arg1, arg2)
	ret0, _ := ret[0].(*PriceFeed)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockIPriceFeedReaderMockRecorder) Read(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockIPriceFeedReader)(nil).Read), arg0, arg1, arg2)
}
