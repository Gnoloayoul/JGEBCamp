// Code generated by MockGen. DO NOT EDIT.
// Source: ./ranking.go

// Package svcmocks is a generated GoMock package.
package svcmocks

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockRankingService is a mock of RankingService interface.
type MockRankingService struct {
	ctrl     *gomock.Controller
	recorder *MockRankingServiceMockRecorder
}

// MockRankingServiceMockRecorder is the mock recorder for MockRankingService.
type MockRankingServiceMockRecorder struct {
	mock *MockRankingService
}

// NewMockRankingService creates a new mock instance.
func NewMockRankingService(ctrl *gomock.Controller) *MockRankingService {
	mock := &MockRankingService{ctrl: ctrl}
	mock.recorder = &MockRankingServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRankingService) EXPECT() *MockRankingServiceMockRecorder {
	return m.recorder
}

// TopN mocks base method.
func (m *MockRankingService) TopN(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TopN", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// TopN indicates an expected call of TopN.
func (mr *MockRankingServiceMockRecorder) TopN(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TopN", reflect.TypeOf((*MockRankingService)(nil).TopN), ctx)
}
