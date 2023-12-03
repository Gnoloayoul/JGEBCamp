// Code generated by MockGen. DO NOT EDIT.
// Source: ./interactive.go

// Package daomocks is a generated GoMock package.
package daomocks

import (
	context "context"
	reflect "reflect"

	dao "github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/dao"
	gomock "go.uber.org/mock/gomock"
)

// MockInteractiveDAO is a mock of InteractiveDAO interface.
type MockInteractiveDAO struct {
	ctrl     *gomock.Controller
	recorder *MockInteractiveDAOMockRecorder
}

// MockInteractiveDAOMockRecorder is the mock recorder for MockInteractiveDAO.
type MockInteractiveDAOMockRecorder struct {
	mock *MockInteractiveDAO
}

// NewMockInteractiveDAO creates a new mock instance.
func NewMockInteractiveDAO(ctrl *gomock.Controller) *MockInteractiveDAO {
	mock := &MockInteractiveDAO{ctrl: ctrl}
	mock.recorder = &MockInteractiveDAOMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInteractiveDAO) EXPECT() *MockInteractiveDAOMockRecorder {
	return m.recorder
}

// BatchIncrReadCnt mocks base method.
func (m *MockInteractiveDAO) BatchIncrReadCnt(ctx context.Context, bizs []string, ids []int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchIncrReadCnt", ctx, bizs, ids)
	ret0, _ := ret[0].(error)
	return ret0
}

// BatchIncrReadCnt indicates an expected call of BatchIncrReadCnt.
func (mr *MockInteractiveDAOMockRecorder) BatchIncrReadCnt(ctx, bizs, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchIncrReadCnt", reflect.TypeOf((*MockInteractiveDAO)(nil).BatchIncrReadCnt), ctx, bizs, ids)
}

// DeleteLikeInfo mocks base method.
func (m *MockInteractiveDAO) DeleteLikeInfo(ctx context.Context, biz string, bizId, uid int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLikeInfo", ctx, biz, bizId, uid)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLikeInfo indicates an expected call of DeleteLikeInfo.
func (mr *MockInteractiveDAOMockRecorder) DeleteLikeInfo(ctx, biz, bizId, uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLikeInfo", reflect.TypeOf((*MockInteractiveDAO)(nil).DeleteLikeInfo), ctx, biz, bizId, uid)
}

// Get mocks base method.
func (m *MockInteractiveDAO) Get(ctx context.Context, biz string, bizId int64) (dao.Interactive, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, biz, bizId)
	ret0, _ := ret[0].(dao.Interactive)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockInteractiveDAOMockRecorder) Get(ctx, biz, bizId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockInteractiveDAO)(nil).Get), ctx, biz, bizId)
}

// GetCollectionInfo mocks base method.
func (m *MockInteractiveDAO) GetCollectionInfo(ctx context.Context, biz string, bizId, uid int64) (dao.UserCollectionBiz, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCollectionInfo", ctx, biz, bizId, uid)
	ret0, _ := ret[0].(dao.UserCollectionBiz)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCollectionInfo indicates an expected call of GetCollectionInfo.
func (mr *MockInteractiveDAOMockRecorder) GetCollectionInfo(ctx, biz, bizId, uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCollectionInfo", reflect.TypeOf((*MockInteractiveDAO)(nil).GetCollectionInfo), ctx, biz, bizId, uid)
}

// GetLikeInfo mocks base method.
func (m *MockInteractiveDAO) GetLikeInfo(ctx context.Context, biz string, bizId, uid int64) (dao.UserLikeBiz, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLikeInfo", ctx, biz, bizId, uid)
	ret0, _ := ret[0].(dao.UserLikeBiz)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLikeInfo indicates an expected call of GetLikeInfo.
func (mr *MockInteractiveDAOMockRecorder) GetLikeInfo(ctx, biz, bizId, uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLikeInfo", reflect.TypeOf((*MockInteractiveDAO)(nil).GetLikeInfo), ctx, biz, bizId, uid)
}

// IncrReadCnt mocks base method.
func (m *MockInteractiveDAO) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrReadCnt", ctx, biz, bizId)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrReadCnt indicates an expected call of IncrReadCnt.
func (mr *MockInteractiveDAOMockRecorder) IncrReadCnt(ctx, biz, bizId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrReadCnt", reflect.TypeOf((*MockInteractiveDAO)(nil).IncrReadCnt), ctx, biz, bizId)
}

// InsertCollectionBiz mocks base method.
func (m *MockInteractiveDAO) InsertCollectionBiz(ctx context.Context, cb dao.UserCollectionBiz) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertCollectionBiz", ctx, cb)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertCollectionBiz indicates an expected call of InsertCollectionBiz.
func (mr *MockInteractiveDAOMockRecorder) InsertCollectionBiz(ctx, cb interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertCollectionBiz", reflect.TypeOf((*MockInteractiveDAO)(nil).InsertCollectionBiz), ctx, cb)
}

// InsertLikeInfo mocks base method.
func (m *MockInteractiveDAO) InsertLikeInfo(ctx context.Context, biz string, bizId, uid int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertLikeInfo", ctx, biz, bizId, uid)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertLikeInfo indicates an expected call of InsertLikeInfo.
func (mr *MockInteractiveDAOMockRecorder) InsertLikeInfo(ctx, biz, bizId, uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertLikeInfo", reflect.TypeOf((*MockInteractiveDAO)(nil).InsertLikeInfo), ctx, biz, bizId, uid)
}
