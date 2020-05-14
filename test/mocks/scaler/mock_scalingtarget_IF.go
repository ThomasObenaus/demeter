// Code generated by MockGen. DO NOT EDIT.
// Source: scaler/scalingtarget_IF.go

// Package mock_scaler is a generated GoMock package.
package mock_scaler

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockScalingTarget is a mock of ScalingTarget interface
type MockScalingTarget struct {
	ctrl     *gomock.Controller
	recorder *MockScalingTargetMockRecorder
}

// MockScalingTargetMockRecorder is the mock recorder for MockScalingTarget
type MockScalingTargetMockRecorder struct {
	mock *MockScalingTarget
}

// NewMockScalingTarget creates a new mock instance
func NewMockScalingTarget(ctrl *gomock.Controller) *MockScalingTarget {
	mock := &MockScalingTarget{ctrl: ctrl}
	mock.recorder = &MockScalingTargetMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockScalingTarget) EXPECT() *MockScalingTargetMockRecorder {
	return m.recorder
}

// AdjustScalingObjectCount mocks base method
func (m *MockScalingTarget) AdjustScalingObjectCount(scalingObject string, min, max, from, to uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AdjustScalingObjectCount", scalingObject, min, max, from, to)
	ret0, _ := ret[0].(error)
	return ret0
}

// AdjustScalingObjectCount indicates an expected call of AdjustScalingObjectCount
func (mr *MockScalingTargetMockRecorder) AdjustScalingObjectCount(scalingObject, min, max, from, to interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AdjustScalingObjectCount", reflect.TypeOf((*MockScalingTarget)(nil).AdjustScalingObjectCount), scalingObject, min, max, from, to)
}

// GetScalingObjectCount mocks base method
func (m *MockScalingTarget) GetScalingObjectCount(scalingObject string) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetScalingObjectCount", scalingObject)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetScalingObjectCount indicates an expected call of GetScalingObjectCount
func (mr *MockScalingTargetMockRecorder) GetScalingObjectCount(scalingObject interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetScalingObjectCount", reflect.TypeOf((*MockScalingTarget)(nil).GetScalingObjectCount), scalingObject)
}

// IsScalingObjectDead mocks base method
func (m *MockScalingTarget) IsScalingObjectDead(scalingObject string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsScalingObjectDead", scalingObject)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsScalingObjectDead indicates an expected call of IsScalingObjectDead
func (mr *MockScalingTargetMockRecorder) IsScalingObjectDead(scalingObject interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsScalingObjectDead", reflect.TypeOf((*MockScalingTarget)(nil).IsScalingObjectDead), scalingObject)
}

// String mocks base method
func (m *MockScalingTarget) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String
func (mr *MockScalingTargetMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockScalingTarget)(nil).String))
}