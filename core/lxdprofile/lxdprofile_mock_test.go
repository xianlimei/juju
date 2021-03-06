// Code generated by MockGen. DO NOT EDIT.
// Source: gopkg.in/juju/charm.v6 (interfaces: LXDProfiler)

// Package lxdprofile_test is a generated GoMock package.
package lxdprofile_test

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	charm_v6 "gopkg.in/juju/charm.v6"
)

// MockLXDProfiler is a mock of LXDProfiler interface
type MockLXDProfiler struct {
	ctrl     *gomock.Controller
	recorder *MockLXDProfilerMockRecorder
}

// MockLXDProfilerMockRecorder is the mock recorder for MockLXDProfiler
type MockLXDProfilerMockRecorder struct {
	mock *MockLXDProfiler
}

// NewMockLXDProfiler creates a new mock instance
func NewMockLXDProfiler(ctrl *gomock.Controller) *MockLXDProfiler {
	mock := &MockLXDProfiler{ctrl: ctrl}
	mock.recorder = &MockLXDProfilerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLXDProfiler) EXPECT() *MockLXDProfilerMockRecorder {
	return m.recorder
}

// LXDProfile mocks base method
func (m *MockLXDProfiler) LXDProfile() *charm_v6.LXDProfile {
	ret := m.ctrl.Call(m, "LXDProfile")
	ret0, _ := ret[0].(*charm_v6.LXDProfile)
	return ret0
}

// LXDProfile indicates an expected call of LXDProfile
func (mr *MockLXDProfilerMockRecorder) LXDProfile() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LXDProfile", reflect.TypeOf((*MockLXDProfiler)(nil).LXDProfile))
}
