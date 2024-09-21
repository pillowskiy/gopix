// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/delivery/rest/handlers/user.go
//
// Generated by this command:
//
//	mockgen -source=./internal/delivery/rest/handlers/user.go -destination=./internal/delivery/rest/handlers/mock/mock_user.go
//

// Package mock_handlers is a generated GoMock package.
package mock_handlers

import (
	context "context"
	reflect "reflect"

	domain "github.com/pillowskiy/gopix/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockuserUseCase is a mock of userUseCase interface.
type MockuserUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockuserUseCaseMockRecorder
}

// MockuserUseCaseMockRecorder is the mock recorder for MockuserUseCase.
type MockuserUseCaseMockRecorder struct {
	mock *MockuserUseCase
}

// NewMockuserUseCase creates a new mock instance.
func NewMockuserUseCase(ctrl *gomock.Controller) *MockuserUseCase {
	mock := &MockuserUseCase{ctrl: ctrl}
	mock.recorder = &MockuserUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockuserUseCase) EXPECT() *MockuserUseCaseMockRecorder {
	return m.recorder
}

// OverwritePermissions mocks base method.
func (m *MockuserUseCase) OverwritePermissions(ctx context.Context, id domain.ID, deny, allow domain.Permission) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OverwritePermissions", ctx, id, deny, allow)
	ret0, _ := ret[0].(error)
	return ret0
}

// OverwritePermissions indicates an expected call of OverwritePermissions.
func (mr *MockuserUseCaseMockRecorder) OverwritePermissions(ctx, id, deny, allow any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OverwritePermissions", reflect.TypeOf((*MockuserUseCase)(nil).OverwritePermissions), ctx, id, deny, allow)
}

// Update mocks base method.
func (m *MockuserUseCase) Update(ctx context.Context, id domain.ID, user *domain.User) (*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, user)
	ret0, _ := ret[0].(*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockuserUseCaseMockRecorder) Update(ctx, id, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockuserUseCase)(nil).Update), ctx, id, user)
}
