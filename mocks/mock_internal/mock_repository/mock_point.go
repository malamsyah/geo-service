// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/point.go
//
// Generated by this command:
//
//	mockgen -source=internal/repository/point.go -destination=mocks/mock_internal/mock_repository/mock_point.go
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	models "github.com/malamsyah/geo-service/internal/models"
	gomock "go.uber.org/mock/gomock"
)

// MockPointRepository is a mock of PointRepository interface.
type MockPointRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPointRepositoryMockRecorder
}

// MockPointRepositoryMockRecorder is the mock recorder for MockPointRepository.
type MockPointRepositoryMockRecorder struct {
	mock *MockPointRepository
}

// NewMockPointRepository creates a new mock instance.
func NewMockPointRepository(ctrl *gomock.Controller) *MockPointRepository {
	mock := &MockPointRepository{ctrl: ctrl}
	mock.recorder = &MockPointRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPointRepository) EXPECT() *MockPointRepositoryMockRecorder {
	return m.recorder
}

// CreatePoint mocks base method.
func (m *MockPointRepository) CreatePoint(point *models.Point) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePoint", point)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreatePoint indicates an expected call of CreatePoint.
func (mr *MockPointRepositoryMockRecorder) CreatePoint(point any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePoint", reflect.TypeOf((*MockPointRepository)(nil).CreatePoint), point)
}

// DeletePoint mocks base method.
func (m *MockPointRepository) DeletePoint(id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePoint", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePoint indicates an expected call of DeletePoint.
func (mr *MockPointRepositoryMockRecorder) DeletePoint(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePoint", reflect.TypeOf((*MockPointRepository)(nil).DeletePoint), id)
}

// GetPointByID mocks base method.
func (m *MockPointRepository) GetPointByID(id uint) (*models.Point, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPointByID", id)
	ret0, _ := ret[0].(*models.Point)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPointByID indicates an expected call of GetPointByID.
func (mr *MockPointRepositoryMockRecorder) GetPointByID(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPointByID", reflect.TypeOf((*MockPointRepository)(nil).GetPointByID), id)
}

// GetPoints mocks base method.
func (m *MockPointRepository) GetPoints(offset, limit int) ([]models.Point, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPoints", offset, limit)
	ret0, _ := ret[0].([]models.Point)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPoints indicates an expected call of GetPoints.
func (mr *MockPointRepositoryMockRecorder) GetPoints(offset, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPoints", reflect.TypeOf((*MockPointRepository)(nil).GetPoints), offset, limit)
}

// GetPointsByContourID mocks base method.
func (m *MockPointRepository) GetPointsByContourID(contourID uint) ([]models.Point, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPointsByContourID", contourID)
	ret0, _ := ret[0].([]models.Point)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPointsByContourID indicates an expected call of GetPointsByContourID.
func (mr *MockPointRepositoryMockRecorder) GetPointsByContourID(contourID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPointsByContourID", reflect.TypeOf((*MockPointRepository)(nil).GetPointsByContourID), contourID)
}

// UpdatePoint mocks base method.
func (m *MockPointRepository) UpdatePoint(point *models.Point) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePoint", point)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePoint indicates an expected call of UpdatePoint.
func (mr *MockPointRepositoryMockRecorder) UpdatePoint(point any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePoint", reflect.TypeOf((*MockPointRepository)(nil).UpdatePoint), point)
}