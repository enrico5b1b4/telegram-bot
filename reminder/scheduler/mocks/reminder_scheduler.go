// Code generated by MockGen. DO NOT EDIT.
// Source: reminder_scheduler.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	reminder "github.com/enrico5b1b4/telegram-bot/reminder"
	gomock "github.com/golang/mock/gomock"
)

// MockScheduler is a mock of Scheduler interface
type MockScheduler struct {
	ctrl     *gomock.Controller
	recorder *MockSchedulerMockRecorder
}

// MockSchedulerMockRecorder is the mock recorder for MockScheduler
type MockSchedulerMockRecorder struct {
	mock *MockScheduler
}

// NewMockScheduler creates a new mock instance
func NewMockScheduler(ctrl *gomock.Controller) *MockScheduler {
	mock := &MockScheduler{ctrl: ctrl}
	mock.recorder = &MockSchedulerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockScheduler) EXPECT() *MockSchedulerMockRecorder {
	return m.recorder
}

// AddReminder mocks base method
func (m *MockScheduler) AddReminder(r *reminder.Reminder) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddReminder", r)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddReminder indicates an expected call of AddReminder
func (mr *MockSchedulerMockRecorder) AddReminder(r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddReminder", reflect.TypeOf((*MockScheduler)(nil).AddReminder), r)
}

// GetNextScheduleTime mocks base method
func (m *MockScheduler) GetNextScheduleTime(chatID, reminderID int) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNextScheduleTime", chatID, reminderID)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNextScheduleTime indicates an expected call of GetNextScheduleTime
func (mr *MockSchedulerMockRecorder) GetNextScheduleTime(chatID, reminderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNextScheduleTime", reflect.TypeOf((*MockScheduler)(nil).GetNextScheduleTime), chatID, reminderID)
}
