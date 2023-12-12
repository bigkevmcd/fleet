// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/rancher/gitjob/pkg/controller (interfaces: GitPoller)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "github.com/rancher/gitjob/pkg/apis/gitjob.cattle.io/v1"
)

// MockGitPoller is a mock of GitPoller interface.
type MockGitPoller struct {
	ctrl     *gomock.Controller
	recorder *MockGitPollerMockRecorder
}

// MockGitPollerMockRecorder is the mock recorder for MockGitPoller.
type MockGitPollerMockRecorder struct {
	mock *MockGitPoller
}

// NewMockGitPoller creates a new mock instance.
func NewMockGitPoller(ctrl *gomock.Controller) *MockGitPoller {
	mock := &MockGitPoller{ctrl: ctrl}
	mock.recorder = &MockGitPollerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGitPoller) EXPECT() *MockGitPollerMockRecorder {
	return m.recorder
}

// AddOrModifyGitRepoWatch mocks base method.
func (m *MockGitPoller) AddOrModifyGitRepoWatch(arg0 context.Context, arg1 v1.GitJob) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddOrModifyGitRepoWatch", arg0, arg1)
}

// AddOrModifyGitRepoWatch indicates an expected call of AddOrModifyGitRepoWatch.
func (mr *MockGitPollerMockRecorder) AddOrModifyGitRepoWatch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOrModifyGitRepoWatch", reflect.TypeOf((*MockGitPoller)(nil).AddOrModifyGitRepoWatch), arg0, arg1)
}

// CleanUpWatches mocks base method.
func (m *MockGitPoller) CleanUpWatches(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CleanUpWatches", arg0)
}

// CleanUpWatches indicates an expected call of CleanUpWatches.
func (mr *MockGitPollerMockRecorder) CleanUpWatches(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CleanUpWatches", reflect.TypeOf((*MockGitPoller)(nil).CleanUpWatches), arg0)
}