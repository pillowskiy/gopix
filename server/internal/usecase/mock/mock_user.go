// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/usecase/user.go
//
// Generated by this command:
//
//	mockgen -source=./internal/usecase/user.go -destination=./internal/usecase/mock/mock_user.go
//

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"

	domain "github.com/pillowskiy/gopix/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockUserCache is a mock of UserCache interface.
type MockUserCache struct {
	ctrl     *gomock.Controller
	recorder *MockUserCacheMockRecorder
}

// MockUserCacheMockRecorder is the mock recorder for MockUserCache.
type MockUserCacheMockRecorder struct {
	mock *MockUserCache
}

// NewMockUserCache creates a new mock instance.
func NewMockUserCache(ctrl *gomock.Controller) *MockUserCache {
	mock := &MockUserCache{ctrl: ctrl}
	mock.recorder = &MockUserCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserCache) EXPECT() *MockUserCacheMockRecorder {
	return m.recorder
}

// Del mocks base method.
func (m *MockUserCache) Del(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Del", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Del indicates an expected call of Del.
func (mr *MockUserCacheMockRecorder) Del(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Del", reflect.TypeOf((*MockUserCache)(nil).Del), ctx, id)
}

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// GetByID mocks base method.
func (m *MockUserRepository) GetByID(ctx context.Context, id domain.ID) (*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockUserRepositoryMockRecorder) GetByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockUserRepository)(nil).GetByID), ctx, id)
}

// GetUnique mocks base method.
func (m *MockUserRepository) GetUnique(ctx context.Context, user *domain.User) (*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnique", ctx, user)
	ret0, _ := ret[0].(*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnique indicates an expected call of GetUnique.
func (mr *MockUserRepositoryMockRecorder) GetUnique(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnique", reflect.TypeOf((*MockUserRepository)(nil).GetUnique), ctx, user)
}

// SetPermissions mocks base method.
func (m *MockUserRepository) SetPermissions(ctx context.Context, id domain.ID, permissions int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetPermissions", ctx, id, permissions)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetPermissions indicates an expected call of SetPermissions.
func (mr *MockUserRepositoryMockRecorder) SetPermissions(ctx, id, permissions any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPermissions", reflect.TypeOf((*MockUserRepository)(nil).SetPermissions), ctx, id, permissions)
}

// Update mocks base method.
func (m *MockUserRepository) Update(ctx context.Context, id domain.ID, user *domain.User) (*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, user)
	ret0, _ := ret[0].(*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockUserRepositoryMockRecorder) Update(ctx, id, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUserRepository)(nil).Update), ctx, id, user)
}

// MockUserFollowingUseCase is a mock of UserFollowingUseCase interface.
type MockUserFollowingUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUserFollowingUseCaseMockRecorder
}

// MockUserFollowingUseCaseMockRecorder is the mock recorder for MockUserFollowingUseCase.
type MockUserFollowingUseCaseMockRecorder struct {
	mock *MockUserFollowingUseCase
}

// NewMockUserFollowingUseCase creates a new mock instance.
func NewMockUserFollowingUseCase(ctrl *gomock.Controller) *MockUserFollowingUseCase {
	mock := &MockUserFollowingUseCase{ctrl: ctrl}
	mock.recorder = &MockUserFollowingUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserFollowingUseCase) EXPECT() *MockUserFollowingUseCaseMockRecorder {
	return m.recorder
}

// Stats mocks base method.
func (m *MockUserFollowingUseCase) Stats(ctx context.Context, userID domain.ID, executorID *domain.ID) (*domain.FollowingStats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stats", ctx, userID, executorID)
	ret0, _ := ret[0].(*domain.FollowingStats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stats indicates an expected call of Stats.
func (mr *MockUserFollowingUseCaseMockRecorder) Stats(ctx, userID, executorID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stats", reflect.TypeOf((*MockUserFollowingUseCase)(nil).Stats), ctx, userID, executorID)
}
