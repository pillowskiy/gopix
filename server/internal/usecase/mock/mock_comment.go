// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/usecase/comment.go
//
// Generated by this command:
//
//	mockgen -source=./internal/usecase/comment.go -destination=./internal/usecase/mock/mock_comment.go
//

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"

	domain "github.com/pillowskiy/gopix/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockCommentRepository is a mock of CommentRepository interface.
type MockCommentRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCommentRepositoryMockRecorder
}

// MockCommentRepositoryMockRecorder is the mock recorder for MockCommentRepository.
type MockCommentRepositoryMockRecorder struct {
	mock *MockCommentRepository
}

// NewMockCommentRepository creates a new mock instance.
func NewMockCommentRepository(ctrl *gomock.Controller) *MockCommentRepository {
	mock := &MockCommentRepository{ctrl: ctrl}
	mock.recorder = &MockCommentRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommentRepository) EXPECT() *MockCommentRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockCommentRepository) Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, comment)
	ret0, _ := ret[0].(*domain.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockCommentRepositoryMockRecorder) Create(ctx, comment any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCommentRepository)(nil).Create), ctx, comment)
}

// Delete mocks base method.
func (m *MockCommentRepository) Delete(ctx context.Context, commentID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, commentID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCommentRepositoryMockRecorder) Delete(ctx, commentID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCommentRepository)(nil).Delete), ctx, commentID)
}

// GetByID mocks base method.
func (m *MockCommentRepository) GetByID(ctx context.Context, imageID int) (*domain.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, imageID)
	ret0, _ := ret[0].(*domain.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockCommentRepositoryMockRecorder) GetByID(ctx, imageID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockCommentRepository)(nil).GetByID), ctx, imageID)
}

// GetByImageID mocks base method.
func (m *MockCommentRepository) GetByImageID(ctx context.Context, imageID int, pagInput *domain.PaginationInput, sort domain.CommentSortMethod) (*domain.Pagination[domain.DetailedComment], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByImageID", ctx, imageID, pagInput, sort)
	ret0, _ := ret[0].(*domain.Pagination[domain.DetailedComment])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByImageID indicates an expected call of GetByImageID.
func (mr *MockCommentRepositoryMockRecorder) GetByImageID(ctx, imageID, pagInput, sort any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByImageID", reflect.TypeOf((*MockCommentRepository)(nil).GetByImageID), ctx, imageID, pagInput, sort)
}

// HasUserCommented mocks base method.
func (m *MockCommentRepository) HasUserCommented(ctx context.Context, commentID, userID int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasUserCommented", ctx, commentID, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasUserCommented indicates an expected call of HasUserCommented.
func (mr *MockCommentRepositoryMockRecorder) HasUserCommented(ctx, commentID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasUserCommented", reflect.TypeOf((*MockCommentRepository)(nil).HasUserCommented), ctx, commentID, userID)
}

// Update mocks base method.
func (m *MockCommentRepository) Update(ctx context.Context, commentID int, comment *domain.Comment) (*domain.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, commentID, comment)
	ret0, _ := ret[0].(*domain.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockCommentRepositoryMockRecorder) Update(ctx, commentID, comment any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCommentRepository)(nil).Update), ctx, commentID, comment)
}

// MockCommentAccessPolicy is a mock of CommentAccessPolicy interface.
type MockCommentAccessPolicy struct {
	ctrl     *gomock.Controller
	recorder *MockCommentAccessPolicyMockRecorder
}

// MockCommentAccessPolicyMockRecorder is the mock recorder for MockCommentAccessPolicy.
type MockCommentAccessPolicyMockRecorder struct {
	mock *MockCommentAccessPolicy
}

// NewMockCommentAccessPolicy creates a new mock instance.
func NewMockCommentAccessPolicy(ctrl *gomock.Controller) *MockCommentAccessPolicy {
	mock := &MockCommentAccessPolicy{ctrl: ctrl}
	mock.recorder = &MockCommentAccessPolicyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommentAccessPolicy) EXPECT() *MockCommentAccessPolicyMockRecorder {
	return m.recorder
}

// CanModify mocks base method.
func (m *MockCommentAccessPolicy) CanModify(user *domain.User, comment *domain.Comment) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CanModify", user, comment)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CanModify indicates an expected call of CanModify.
func (mr *MockCommentAccessPolicyMockRecorder) CanModify(user, comment any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CanModify", reflect.TypeOf((*MockCommentAccessPolicy)(nil).CanModify), user, comment)
}

// MockCommentImageUseCase is a mock of CommentImageUseCase interface.
type MockCommentImageUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockCommentImageUseCaseMockRecorder
}

// MockCommentImageUseCaseMockRecorder is the mock recorder for MockCommentImageUseCase.
type MockCommentImageUseCaseMockRecorder struct {
	mock *MockCommentImageUseCase
}

// NewMockCommentImageUseCase creates a new mock instance.
func NewMockCommentImageUseCase(ctrl *gomock.Controller) *MockCommentImageUseCase {
	mock := &MockCommentImageUseCase{ctrl: ctrl}
	mock.recorder = &MockCommentImageUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommentImageUseCase) EXPECT() *MockCommentImageUseCaseMockRecorder {
	return m.recorder
}

// GetByID mocks base method.
func (m *MockCommentImageUseCase) GetByID(ctx context.Context, imageID int) (*domain.Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, imageID)
	ret0, _ := ret[0].(*domain.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockCommentImageUseCaseMockRecorder) GetByID(ctx, imageID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockCommentImageUseCase)(nil).GetByID), ctx, imageID)
}
