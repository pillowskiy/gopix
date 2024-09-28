// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/delivery/rest/handlers/tag.go
//
// Generated by this command:
//
//	mockgen -source=./internal/delivery/rest/handlers/tag.go -destination=./internal/delivery/rest/handlers/mock/mock_tag.go
//

// Package mock_handlers is a generated GoMock package.
package mock_handlers

import (
	context "context"
	reflect "reflect"

	domain "github.com/pillowskiy/gopix/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockTagUseCase is a mock of TagUseCase interface.
type MockTagUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockTagUseCaseMockRecorder
}

// MockTagUseCaseMockRecorder is the mock recorder for MockTagUseCase.
type MockTagUseCaseMockRecorder struct {
	mock *MockTagUseCase
}

// NewMockTagUseCase creates a new mock instance.
func NewMockTagUseCase(ctrl *gomock.Controller) *MockTagUseCase {
	mock := &MockTagUseCase{ctrl: ctrl}
	mock.recorder = &MockTagUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTagUseCase) EXPECT() *MockTagUseCaseMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTagUseCase) Create(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, tag)
	ret0, _ := ret[0].(*domain.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockTagUseCaseMockRecorder) Create(ctx, tag any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTagUseCase)(nil).Create), ctx, tag)
}

// Delete mocks base method.
func (m *MockTagUseCase) Delete(ctx context.Context, tagID domain.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, tagID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTagUseCaseMockRecorder) Delete(ctx, tagID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTagUseCase)(nil).Delete), ctx, tagID)
}

// DeleteImageTag mocks base method.
func (m *MockTagUseCase) DeleteImageTag(ctx context.Context, tagID, imageID domain.ID, executor *domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteImageTag", ctx, tagID, imageID, executor)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteImageTag indicates an expected call of DeleteImageTag.
func (mr *MockTagUseCaseMockRecorder) DeleteImageTag(ctx, tagID, imageID, executor any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteImageTag", reflect.TypeOf((*MockTagUseCase)(nil).DeleteImageTag), ctx, tagID, imageID, executor)
}

// Search mocks base method.
func (m *MockTagUseCase) Search(ctx context.Context, query string) ([]domain.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, query)
	ret0, _ := ret[0].([]domain.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockTagUseCaseMockRecorder) Search(ctx, query any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockTagUseCase)(nil).Search), ctx, query)
}

// UpsertImageTag mocks base method.
func (m *MockTagUseCase) UpsertImageTag(ctx context.Context, tag *domain.Tag, imageID domain.ID, executor *domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertImageTag", ctx, tag, imageID, executor)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertImageTag indicates an expected call of UpsertImageTag.
func (mr *MockTagUseCaseMockRecorder) UpsertImageTag(ctx, tag, imageID, executor any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertImageTag", reflect.TypeOf((*MockTagUseCase)(nil).UpsertImageTag), ctx, tag, imageID, executor)
}
