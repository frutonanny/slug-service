// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_get_report is a generated GoMock package.
package mock_get_report

import (
	context "context"
	io "io"
	reflect "reflect"
	time "time"

	events "github.com/frutonanny/slug-service/internal/repositories/events"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	minio "github.com/minio/minio-go/v7"
)

// MockgetUserSlugsService is a mock of getUserSlugsService interface.
type MockgetUserSlugsService struct {
	ctrl     *gomock.Controller
	recorder *MockgetUserSlugsServiceMockRecorder
}

// MockgetUserSlugsServiceMockRecorder is the mock recorder for MockgetUserSlugsService.
type MockgetUserSlugsServiceMockRecorder struct {
	mock *MockgetUserSlugsService
}

// NewMockgetUserSlugsService creates a new mock instance.
func NewMockgetUserSlugsService(ctrl *gomock.Controller) *MockgetUserSlugsService {
	mock := &MockgetUserSlugsService{ctrl: ctrl}
	mock.recorder = &MockgetUserSlugsServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockgetUserSlugsService) EXPECT() *MockgetUserSlugsServiceMockRecorder {
	return m.recorder
}

// GetUserSlugs mocks base method.
func (m *MockgetUserSlugsService) GetUserSlugs(ctx context.Context, userID uuid.UUID, sync bool) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserSlugs", ctx, userID, sync)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserSlugs indicates an expected call of GetUserSlugs.
func (mr *MockgetUserSlugsServiceMockRecorder) GetUserSlugs(ctx, userID, sync interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserSlugs", reflect.TypeOf((*MockgetUserSlugsService)(nil).GetUserSlugs), ctx, userID, sync)
}

// MockeventsRepo is a mock of eventsRepo interface.
type MockeventsRepo struct {
	ctrl     *gomock.Controller
	recorder *MockeventsRepoMockRecorder
}

// MockeventsRepoMockRecorder is the mock recorder for MockeventsRepo.
type MockeventsRepoMockRecorder struct {
	mock *MockeventsRepo
}

// NewMockeventsRepo creates a new mock instance.
func NewMockeventsRepo(ctrl *gomock.Controller) *MockeventsRepo {
	mock := &MockeventsRepo{ctrl: ctrl}
	mock.recorder = &MockeventsRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockeventsRepo) EXPECT() *MockeventsRepoMockRecorder {
	return m.recorder
}

// GetReport mocks base method.
func (m *MockeventsRepo) GetReport(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]events.UserReportEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReport", ctx, userID, from, to)
	ret0, _ := ret[0].([]events.UserReportEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReport indicates an expected call of GetReport.
func (mr *MockeventsRepoMockRecorder) GetReport(ctx, userID, from, to interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReport", reflect.TypeOf((*MockeventsRepo)(nil).GetReport), ctx, userID, from, to)
}

// MockminioClient is a mock of minioClient interface.
type MockminioClient struct {
	ctrl     *gomock.Controller
	recorder *MockminioClientMockRecorder
}

// MockminioClientMockRecorder is the mock recorder for MockminioClient.
type MockminioClientMockRecorder struct {
	mock *MockminioClient
}

// NewMockminioClient creates a new mock instance.
func NewMockminioClient(ctrl *gomock.Controller) *MockminioClient {
	mock := &MockminioClient{ctrl: ctrl}
	mock.recorder = &MockminioClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockminioClient) EXPECT() *MockminioClientMockRecorder {
	return m.recorder
}

// PutObject mocks base method.
func (m *MockminioClient) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutObject", ctx, bucketName, objectName, reader, objectSize, opts)
	ret0, _ := ret[0].(minio.UploadInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutObject indicates an expected call of PutObject.
func (mr *MockminioClientMockRecorder) PutObject(ctx, bucketName, objectName, reader, objectSize, opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutObject", reflect.TypeOf((*MockminioClient)(nil).PutObject), ctx, bucketName, objectName, reader, objectSize, opts)
}
