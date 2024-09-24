// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/delivery/rest/handlers/subscription.go
//
// Generated by this command:
//
//	mockgen -source=./internal/delivery/rest/handlers/subscription.go -destination=./internal/delivery/rest/handlers/mock/mock_subscription.go
//

// Package mock_handlers is a generated GoMock package.
package mock_handlers

import (
	context "context"
	reflect "reflect"

	domain "github.com/pillowskiy/gopix/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockSubscriptionUseCase is a mock of SubscriptionUseCase interface.
type MockSubscriptionUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockSubscriptionUseCaseMockRecorder
}

// MockSubscriptionUseCaseMockRecorder is the mock recorder for MockSubscriptionUseCase.
type MockSubscriptionUseCaseMockRecorder struct {
	mock *MockSubscriptionUseCase
}

// NewMockSubscriptionUseCase creates a new mock instance.
func NewMockSubscriptionUseCase(ctrl *gomock.Controller) *MockSubscriptionUseCase {
	mock := &MockSubscriptionUseCase{ctrl: ctrl}
	mock.recorder = &MockSubscriptionUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubscriptionUseCase) EXPECT() *MockSubscriptionUseCaseMockRecorder {
	return m.recorder
}

// Follow mocks base method.
func (m *MockSubscriptionUseCase) Follow(ctx context.Context, userID domain.ID, executor *domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Follow", ctx, userID, executor)
	ret0, _ := ret[0].(error)
	return ret0
}

// Follow indicates an expected call of Follow.
func (mr *MockSubscriptionUseCaseMockRecorder) Follow(ctx, userID, executor any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Follow", reflect.TypeOf((*MockSubscriptionUseCase)(nil).Follow), ctx, userID, executor)
}

// Unfollow mocks base method.
func (m *MockSubscriptionUseCase) Unfollow(ctx context.Context, userID domain.ID, executor *domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unfollow", ctx, userID, executor)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unfollow indicates an expected call of Unfollow.
func (mr *MockSubscriptionUseCaseMockRecorder) Unfollow(ctx, userID, executor any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unfollow", reflect.TypeOf((*MockSubscriptionUseCase)(nil).Unfollow), ctx, userID, executor)
}