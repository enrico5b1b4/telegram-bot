// Code generated by MockGen. DO NOT EDIT.
// Source: reminddetail_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	command "github.com/enrico5b1b4/telegram-bot/pkg/command"
	gomock "github.com/golang/mock/gomock"
)

// MockRemindDetailServicer is a mock of RemindDetailServicer interface
type MockRemindDetailServicer struct {
	ctrl     *gomock.Controller
	recorder *MockRemindDetailServicerMockRecorder
}

// MockRemindDetailServicerMockRecorder is the mock recorder for MockRemindDetailServicer
type MockRemindDetailServicerMockRecorder struct {
	mock *MockRemindDetailServicer
}

// NewMockRemindDetailServicer creates a new mock instance
func NewMockRemindDetailServicer(ctrl *gomock.Controller) *MockRemindDetailServicer {
	mock := &MockRemindDetailServicer{ctrl: ctrl}
	mock.recorder = &MockRemindDetailServicerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRemindDetailServicer) EXPECT() *MockRemindDetailServicerMockRecorder {
	return m.recorder
}

// GetReminder mocks base method
func (m *MockRemindDetailServicer) GetReminder(chatID, reminderID int) (*command.ReminderDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReminder", chatID, reminderID)
	ret0, _ := ret[0].(*command.ReminderDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReminder indicates an expected call of GetReminder
func (mr *MockRemindDetailServicerMockRecorder) GetReminder(chatID, reminderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReminder", reflect.TypeOf((*MockRemindDetailServicer)(nil).GetReminder), chatID, reminderID)
}

// DeleteReminder mocks base method
func (m *MockRemindDetailServicer) DeleteReminder(chatID, ID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteReminder", chatID, ID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteReminder indicates an expected call of DeleteReminder
func (mr *MockRemindDetailServicerMockRecorder) DeleteReminder(chatID, ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteReminder", reflect.TypeOf((*MockRemindDetailServicer)(nil).DeleteReminder), chatID, ID)
}
