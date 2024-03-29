// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package grafana is a generated GoMock package.
package grafana

import (
	reflect "reflect"

	domain "github.com/brotifypacha/grafana-sync/internal/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// GetDashboard mocks base method.
func (m *MockRepository) GetDashboard(uid string) (*domain.GrafanaDashboardDetails, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDashboard", uid)
	ret0, _ := ret[0].(*domain.GrafanaDashboardDetails)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDashboard indicates an expected call of GetDashboard.
func (mr *MockRepositoryMockRecorder) GetDashboard(uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDashboard", reflect.TypeOf((*MockRepository)(nil).GetDashboard), uid)
}

// GetTree mocks base method.
func (m *MockRepository) GetTree() (*domain.GrafanaFolder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTree")
	ret0, _ := ret[0].(*domain.GrafanaFolder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTree indicates an expected call of GetTree.
func (mr *MockRepositoryMockRecorder) GetTree() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTree", reflect.TypeOf((*MockRepository)(nil).GetTree))
}
