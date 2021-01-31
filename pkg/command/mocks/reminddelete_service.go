// Code generated by MockGen. DO NOT EDIT.
// Source: reminddelete_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRemindDeleteServicer is a mock of RemindDeleteServicer interface
type MockRemindDeleteServicer struct {
	ctrl     *gomock.Controller
	recorder *MockRemindDeleteServicerMockRecorder
}

// MockRemindDeleteServicerMockRecorder is the mock recorder for MockRemindDeleteServicer
type MockRemindDeleteServicerMockRecorder struct {
	mock *MockRemindDeleteServicer
}

// NewMockRemindDeleteServicer creates a new mock instance
func NewMockRemindDeleteServicer(ctrl *gomock.Controller) *MockRemindDeleteServicer {
	mock := &MockRemindDeleteServicer{ctrl: ctrl}
	mock.recorder = &MockRemindDeleteServicerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRemindDeleteServicer) EXPECT() *MockRemindDeleteServicerMockRecorder {
	return m.recorder
}

// DeleteReminder mocks base method
func (m *MockRemindDeleteServicer) DeleteReminder(chatID, ID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteReminder", chatID, ID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteReminder indicates an expected call of DeleteReminder
func (mr *MockRemindDeleteServicerMockRecorder) DeleteReminder(chatID, ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteReminder", reflect.TypeOf((*MockRemindDeleteServicer)(nil).DeleteReminder), chatID, ID)
}
