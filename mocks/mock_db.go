// Code generated by MockGen. DO NOT EDIT.
// Source: D:/golangProj/tgbot/internal/repository/db.go
//
// Generated by this command:
//
//	mockgen -source=D:/golangProj/tgbot/internal/repository/db.go -destination=mocks/mock_db .go -package=mocks
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockDB is a mock of DB interface.
type MockDB struct {
	ctrl     *gomock.Controller
	recorder *MockDBMockRecorder
	isgomock struct{}
}

// MockDBMockRecorder is the mock recorder for MockDB.
type MockDBMockRecorder struct {
	mock *MockDB
}

// NewMockDB creates a new mock instance.
func NewMockDB(ctrl *gomock.Controller) *MockDB {
	mock := &MockDB{ctrl: ctrl}
	mock.recorder = &MockDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDB) EXPECT() *MockDBMockRecorder {
	return m.recorder
}

// GetContext mocks base method.
func (m *MockDB) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, dest, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetContext", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetContext indicates an expected call of GetContext.
func (mr *MockDBMockRecorder) GetContext(ctx, dest, query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, dest, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContext", reflect.TypeOf((*MockDB)(nil).GetContext), varargs...)
}

// NamedExecContext mocks base method.
func (m *MockDB) NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NamedExecContext", ctx, query, arg)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NamedExecContext indicates an expected call of NamedExecContext.
func (mr *MockDBMockRecorder) NamedExecContext(ctx, query, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NamedExecContext", reflect.TypeOf((*MockDB)(nil).NamedExecContext), ctx, query, arg)
}

// SelectContext mocks base method.
func (m *MockDB) SelectContext(ctx context.Context, dest any, query string, args ...any) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, dest, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SelectContext", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SelectContext indicates an expected call of SelectContext.
func (mr *MockDBMockRecorder) SelectContext(ctx, dest, query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, dest, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectContext", reflect.TypeOf((*MockDB)(nil).SelectContext), varargs...)
}
