// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/SimonRichardson/discourse/pkg/metrics (interfaces: Gauge,HistogramVec,Counter,Observer)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	prometheus "github.com/prometheus/client_golang/prometheus"
	reflect "reflect"
)

// MockGauge is a mock of Gauge interface
type MockGauge struct {
	ctrl     *gomock.Controller
	recorder *MockGaugeMockRecorder
}

// MockGaugeMockRecorder is the mock recorder for MockGauge
type MockGaugeMockRecorder struct {
	mock *MockGauge
}

// NewMockGauge creates a new mock instance
func NewMockGauge(ctrl *gomock.Controller) *MockGauge {
	mock := &MockGauge{ctrl: ctrl}
	mock.recorder = &MockGaugeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGauge) EXPECT() *MockGaugeMockRecorder {
	return m.recorder
}

// Dec mocks base method
func (m *MockGauge) Dec() {
	m.ctrl.Call(m, "Dec")
}

// Dec indicates an expected call of Dec
func (mr *MockGaugeMockRecorder) Dec() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dec", reflect.TypeOf((*MockGauge)(nil).Dec))
}

// Inc mocks base method
func (m *MockGauge) Inc() {
	m.ctrl.Call(m, "Inc")
}

// Inc indicates an expected call of Inc
func (mr *MockGaugeMockRecorder) Inc() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Inc", reflect.TypeOf((*MockGauge)(nil).Inc))
}

// MockHistogramVec is a mock of HistogramVec interface
type MockHistogramVec struct {
	ctrl     *gomock.Controller
	recorder *MockHistogramVecMockRecorder
}

// MockHistogramVecMockRecorder is the mock recorder for MockHistogramVec
type MockHistogramVecMockRecorder struct {
	mock *MockHistogramVec
}

// NewMockHistogramVec creates a new mock instance
func NewMockHistogramVec(ctrl *gomock.Controller) *MockHistogramVec {
	mock := &MockHistogramVec{ctrl: ctrl}
	mock.recorder = &MockHistogramVecMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHistogramVec) EXPECT() *MockHistogramVecMockRecorder {
	return m.recorder
}

// WithLabelValues mocks base method
func (m *MockHistogramVec) WithLabelValues(arg0 ...string) prometheus.Observer {
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "WithLabelValues", varargs...)
	ret0, _ := ret[0].(prometheus.Observer)
	return ret0
}

// WithLabelValues indicates an expected call of WithLabelValues
func (mr *MockHistogramVecMockRecorder) WithLabelValues(arg0 ...interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithLabelValues", reflect.TypeOf((*MockHistogramVec)(nil).WithLabelValues), arg0...)
}

// MockCounter is a mock of Counter interface
type MockCounter struct {
	ctrl     *gomock.Controller
	recorder *MockCounterMockRecorder
}

// MockCounterMockRecorder is the mock recorder for MockCounter
type MockCounterMockRecorder struct {
	mock *MockCounter
}

// NewMockCounter creates a new mock instance
func NewMockCounter(ctrl *gomock.Controller) *MockCounter {
	mock := &MockCounter{ctrl: ctrl}
	mock.recorder = &MockCounterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCounter) EXPECT() *MockCounterMockRecorder {
	return m.recorder
}

// Add mocks base method
func (m *MockCounter) Add(arg0 float64) {
	m.ctrl.Call(m, "Add", arg0)
}

// Add indicates an expected call of Add
func (mr *MockCounterMockRecorder) Add(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockCounter)(nil).Add), arg0)
}

// Inc mocks base method
func (m *MockCounter) Inc() {
	m.ctrl.Call(m, "Inc")
}

// Inc indicates an expected call of Inc
func (mr *MockCounterMockRecorder) Inc() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Inc", reflect.TypeOf((*MockCounter)(nil).Inc))
}

// MockObserver is a mock of Observer interface
type MockObserver struct {
	ctrl     *gomock.Controller
	recorder *MockObserverMockRecorder
}

// MockObserverMockRecorder is the mock recorder for MockObserver
type MockObserverMockRecorder struct {
	mock *MockObserver
}

// NewMockObserver creates a new mock instance
func NewMockObserver(ctrl *gomock.Controller) *MockObserver {
	mock := &MockObserver{ctrl: ctrl}
	mock.recorder = &MockObserverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockObserver) EXPECT() *MockObserverMockRecorder {
	return m.recorder
}

// Observe mocks base method
func (m *MockObserver) Observe(arg0 float64) {
	m.ctrl.Call(m, "Observe", arg0)
}

// Observe indicates an expected call of Observe
func (mr *MockObserverMockRecorder) Observe(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Observe", reflect.TypeOf((*MockObserver)(nil).Observe), arg0)
}
