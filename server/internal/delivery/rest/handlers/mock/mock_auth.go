// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/delivery/rest/handlers/auth.go
//
// Generated by this command:
//
//	mockgen -source=./internal/delivery/rest/handlers/auth.go -destination=./internal/delivery/rest/handlers/mock/mock_auth.go
//

// Package mock_handlers is a generated GoMock package.
package mock_handlers

import (
	context "context"
	reflect "reflect"

	domain "github.com/pillowskiy/gopix/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockauthUseCase is a mock of authUseCase interface.
type MockauthUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockauthUseCaseMockRecorder
}

// MockauthUseCaseMockRecorder is the mock recorder for MockauthUseCase.
type MockauthUseCaseMockRecorder struct {
	mock *MockauthUseCase
}

// NewMockauthUseCase creates a new mock instance.
func NewMockauthUseCase(ctrl *gomock.Controller) *MockauthUseCase {
	mock := &MockauthUseCase{ctrl: ctrl}
	mock.recorder = &MockauthUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockauthUseCase) EXPECT() *MockauthUseCaseMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *MockauthUseCase) Login(ctx context.Context, user *domain.User) (*domain.UserWithToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, user)
	ret0, _ := ret[0].(*domain.UserWithToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockauthUseCaseMockRecorder) Login(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockauthUseCase)(nil).Login), ctx, user)
}

// Register mocks base method.
func (m *MockauthUseCase) Register(ctx context.Context, user *domain.User) (*domain.UserWithToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, user)
	ret0, _ := ret[0].(*domain.UserWithToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockauthUseCaseMockRecorder) Register(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockauthUseCase)(nil).Register), ctx, user)
}
