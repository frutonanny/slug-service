// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_modify_slug is a generated GoMock package.
package mock_modify_slug

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockusersRepo is a mock of usersRepo interface.
type MockusersRepo struct {
	ctrl     *gomock.Controller
	recorder *MockusersRepoMockRecorder
}

// MockusersRepoMockRecorder is the mock recorder for MockusersRepo.
type MockusersRepoMockRecorder struct {
	mock *MockusersRepo
}

// NewMockusersRepo creates a new mock instance.
func NewMockusersRepo(ctrl *gomock.Controller) *MockusersRepo {
	mock := &MockusersRepo{ctrl: ctrl}
	mock.recorder = &MockusersRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockusersRepo) EXPECT() *MockusersRepoMockRecorder {
	return m.recorder
}

// AddUserSlug mocks base method.
func (m *MockusersRepo) AddUserSlug(ctx context.Context, user uuid.UUID, slugID int64, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserSlug", ctx, user, slugID, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUserSlug indicates an expected call of AddUserSlug.
func (mr *MockusersRepoMockRecorder) AddUserSlug(ctx, user, slugID, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserSlug", reflect.TypeOf((*MockusersRepo)(nil).AddUserSlug), ctx, user, slugID, name)
}

// AddUserSlugWithTtl mocks base method.
func (m *MockusersRepo) AddUserSlugWithTtl(ctx context.Context, userID uuid.UUID, slugID int64, name string, ttl time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserSlugWithTtl", ctx, userID, slugID, name, ttl)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUserSlugWithTtl indicates an expected call of AddUserSlugWithTtl.
func (mr *MockusersRepoMockRecorder) AddUserSlugWithTtl(ctx, userID, slugID, name, ttl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserSlugWithTtl", reflect.TypeOf((*MockusersRepo)(nil).AddUserSlugWithTtl), ctx, userID, slugID, name, ttl)
}

// CreateUserIfNotExist mocks base method.
func (m *MockusersRepo) CreateUserIfNotExist(ctx context.Context, userID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserIfNotExist", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUserIfNotExist indicates an expected call of CreateUserIfNotExist.
func (mr *MockusersRepoMockRecorder) CreateUserIfNotExist(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserIfNotExist", reflect.TypeOf((*MockusersRepo)(nil).CreateUserIfNotExist), ctx, userID)
}

// DeleteUserSlug mocks base method.
func (m *MockusersRepo) DeleteUserSlug(ctx context.Context, user uuid.UUID, slugID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserSlug", ctx, user, slugID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserSlug indicates an expected call of DeleteUserSlug.
func (mr *MockusersRepoMockRecorder) DeleteUserSlug(ctx, user, slugID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserSlug", reflect.TypeOf((*MockusersRepo)(nil).DeleteUserSlug), ctx, user, slugID)
}

// MockslugRepo is a mock of slugRepo interface.
type MockslugRepo struct {
	ctrl     *gomock.Controller
	recorder *MockslugRepoMockRecorder
}

// MockslugRepoMockRecorder is the mock recorder for MockslugRepo.
type MockslugRepoMockRecorder struct {
	mock *MockslugRepo
}

// NewMockslugRepo creates a new mock instance.
func NewMockslugRepo(ctrl *gomock.Controller) *MockslugRepo {
	mock := &MockslugRepo{ctrl: ctrl}
	mock.recorder = &MockslugRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockslugRepo) EXPECT() *MockslugRepoMockRecorder {
	return m.recorder
}

// GetID mocks base method.
func (m *MockslugRepo) GetID(ctx context.Context, name string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetID", ctx, name)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetID indicates an expected call of GetID.
func (mr *MockslugRepoMockRecorder) GetID(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetID", reflect.TypeOf((*MockslugRepo)(nil).GetID), ctx, name)
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

// AddEvents mocks base method.
func (m *MockeventsRepo) AddEvents(ctx context.Context, userID uuid.UUID, slugID int64, event string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddEvents", ctx, userID, slugID, event)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddEvents indicates an expected call of AddEvents.
func (mr *MockeventsRepoMockRecorder) AddEvents(ctx, userID, slugID, event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEvents", reflect.TypeOf((*MockeventsRepo)(nil).AddEvents), ctx, userID, slugID, event)
}

// Mocktransactor is a mock of transactor interface.
type Mocktransactor struct {
	ctrl     *gomock.Controller
	recorder *MocktransactorMockRecorder
}

// MocktransactorMockRecorder is the mock recorder for Mocktransactor.
type MocktransactorMockRecorder struct {
	mock *Mocktransactor
}

// NewMocktransactor creates a new mock instance.
func NewMocktransactor(ctrl *gomock.Controller) *Mocktransactor {
	mock := &Mocktransactor{ctrl: ctrl}
	mock.recorder = &MocktransactorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mocktransactor) EXPECT() *MocktransactorMockRecorder {
	return m.recorder
}

// RunInTx mocks base method.
func (m *Mocktransactor) RunInTx(ctx context.Context, f func(context.Context) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunInTx", ctx, f)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunInTx indicates an expected call of RunInTx.
func (mr *MocktransactorMockRecorder) RunInTx(ctx, f interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunInTx", reflect.TypeOf((*Mocktransactor)(nil).RunInTx), ctx, f)
}