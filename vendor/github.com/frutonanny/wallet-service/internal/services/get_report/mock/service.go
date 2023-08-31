// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_get_report is a generated GoMock package.
package mock_get_report

import (
	context "context"
	io "io"
	reflect "reflect"

	postgres "github.com/frutonanny/wallet-service/internal/postgres"
	report "github.com/frutonanny/wallet-service/internal/repositories/report"
	get_report "github.com/frutonanny/wallet-service/internal/services/get_report"
	gomock "github.com/golang/mock/gomock"
	minio "github.com/minio/minio-go/v7"
)

// Mocklogger is a mock of logger interface.
type Mocklogger struct {
	ctrl     *gomock.Controller
	recorder *MockloggerMockRecorder
}

// MockloggerMockRecorder is the mock recorder for Mocklogger.
type MockloggerMockRecorder struct {
	mock *Mocklogger
}

// NewMocklogger creates a new mock instance.
func NewMocklogger(ctrl *gomock.Controller) *Mocklogger {
	mock := &Mocklogger{ctrl: ctrl}
	mock.recorder = &MockloggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mocklogger) EXPECT() *MockloggerMockRecorder {
	return m.recorder
}

// Error mocks base method.
func (m *Mocklogger) Error(msg string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Error", msg)
}

// Error indicates an expected call of Error.
func (mr *MockloggerMockRecorder) Error(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*Mocklogger)(nil).Error), msg)
}

// Info mocks base method.
func (m *Mocklogger) Info(msg string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Info", msg)
}

// Info indicates an expected call of Info.
func (mr *MockloggerMockRecorder) Info(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*Mocklogger)(nil).Info), msg)
}

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

// GetReport mocks base method.
func (m *MockRepository) GetReport(ctx context.Context, period string) ([]report.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReport", ctx, period)
	ret0, _ := ret[0].([]report.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReport indicates an expected call of GetReport.
func (mr *MockRepositoryMockRecorder) GetReport(ctx, period interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReport", reflect.TypeOf((*MockRepository)(nil).GetReport), ctx, period)
}

// MockMinioClient is a mock of MinioClient interface.
type MockMinioClient struct {
	ctrl     *gomock.Controller
	recorder *MockMinioClientMockRecorder
}

// MockMinioClientMockRecorder is the mock recorder for MockMinioClient.
type MockMinioClientMockRecorder struct {
	mock *MockMinioClient
}

// NewMockMinioClient creates a new mock instance.
func NewMockMinioClient(ctrl *gomock.Controller) *MockMinioClient {
	mock := &MockMinioClient{ctrl: ctrl}
	mock.recorder = &MockMinioClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMinioClient) EXPECT() *MockMinioClientMockRecorder {
	return m.recorder
}

// PutObject mocks base method.
func (m *MockMinioClient) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutObject", ctx, bucketName, objectName, reader, objectSize, opts)
	ret0, _ := ret[0].(minio.UploadInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutObject indicates an expected call of PutObject.
func (mr *MockMinioClientMockRecorder) PutObject(ctx, bucketName, objectName, reader, objectSize, opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutObject", reflect.TypeOf((*MockMinioClient)(nil).PutObject), ctx, bucketName, objectName, reader, objectSize, opts)
}

// Mockdependencies is a mock of dependencies interface.
type Mockdependencies struct {
	ctrl     *gomock.Controller
	recorder *MockdependenciesMockRecorder
}

// MockdependenciesMockRecorder is the mock recorder for Mockdependencies.
type MockdependenciesMockRecorder struct {
	mock *Mockdependencies
}

// NewMockdependencies creates a new mock instance.
func NewMockdependencies(ctrl *gomock.Controller) *Mockdependencies {
	mock := &Mockdependencies{ctrl: ctrl}
	mock.recorder = &MockdependenciesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockdependencies) EXPECT() *MockdependenciesMockRecorder {
	return m.recorder
}

// NewRepository mocks base method.
func (m *Mockdependencies) NewRepository(db postgres.Database) get_report.Repository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRepository", db)
	ret0, _ := ret[0].(get_report.Repository)
	return ret0
}

// NewRepository indicates an expected call of NewRepository.
func (mr *MockdependenciesMockRecorder) NewRepository(db interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRepository", reflect.TypeOf((*Mockdependencies)(nil).NewRepository), db)
}
