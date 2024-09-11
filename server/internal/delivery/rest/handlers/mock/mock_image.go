// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/delivery/rest/handlers/image.go
//
// Generated by this command:
//
//	mockgen -source=./internal/delivery/rest/handlers/image.go -destination=./internal/delivery/rest/handlers/mock/mock_image.go
//

// Package mock_handlers is a generated GoMock package.
package mock_handlers

import (
	context "context"
	reflect "reflect"

	domain "github.com/pillowskiy/gopix/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockimageUseCase is a mock of imageUseCase interface.
type MockimageUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockimageUseCaseMockRecorder
}

// MockimageUseCaseMockRecorder is the mock recorder for MockimageUseCase.
type MockimageUseCaseMockRecorder struct {
	mock *MockimageUseCase
}

// NewMockimageUseCase creates a new mock instance.
func NewMockimageUseCase(ctrl *gomock.Controller) *MockimageUseCase {
	mock := &MockimageUseCase{ctrl: ctrl}
	mock.recorder = &MockimageUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockimageUseCase) EXPECT() *MockimageUseCaseMockRecorder {
	return m.recorder
}

// AddLike mocks base method.
func (m *MockimageUseCase) AddLike(ctx context.Context, imageID, userID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddLike", ctx, imageID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddLike indicates an expected call of AddLike.
func (mr *MockimageUseCaseMockRecorder) AddLike(ctx, imageID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLike", reflect.TypeOf((*MockimageUseCase)(nil).AddLike), ctx, imageID, userID)
}

// AddView mocks base method.
func (m *MockimageUseCase) AddView(ctx context.Context, view *domain.ImageView) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddView", ctx, view)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddView indicates an expected call of AddView.
func (mr *MockimageUseCaseMockRecorder) AddView(ctx, view any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddView", reflect.TypeOf((*MockimageUseCase)(nil).AddView), ctx, view)
}

// Create mocks base method.
func (m *MockimageUseCase) Create(ctx context.Context, image *domain.Image, file *domain.FileNode) (*domain.Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, image, file)
	ret0, _ := ret[0].(*domain.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockimageUseCaseMockRecorder) Create(ctx, image, file any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockimageUseCase)(nil).Create), ctx, image, file)
}

// Delete mocks base method.
func (m *MockimageUseCase) Delete(ctx context.Context, id int, executor *domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id, executor)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockimageUseCaseMockRecorder) Delete(ctx, id, executor any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockimageUseCase)(nil).Delete), ctx, id, executor)
}

// Discover mocks base method.
func (m *MockimageUseCase) Discover(ctx context.Context, pagInput *domain.PaginationInput, sort domain.ImageSortMethod) (*domain.Pagination[domain.Image], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Discover", ctx, pagInput, sort)
	ret0, _ := ret[0].(*domain.Pagination[domain.Image])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Discover indicates an expected call of Discover.
func (mr *MockimageUseCaseMockRecorder) Discover(ctx, pagInput, sort any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Discover", reflect.TypeOf((*MockimageUseCase)(nil).Discover), ctx, pagInput, sort)
}

// GetDetailed mocks base method.
func (m *MockimageUseCase) GetDetailed(ctx context.Context, id int) (*domain.DetailedImage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDetailed", ctx, id)
	ret0, _ := ret[0].(*domain.DetailedImage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDetailed indicates an expected call of GetDetailed.
func (mr *MockimageUseCaseMockRecorder) GetDetailed(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDetailed", reflect.TypeOf((*MockimageUseCase)(nil).GetDetailed), ctx, id)
}

// RemoveLike mocks base method.
func (m *MockimageUseCase) RemoveLike(ctx context.Context, imageID, userID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveLike", ctx, imageID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveLike indicates an expected call of RemoveLike.
func (mr *MockimageUseCaseMockRecorder) RemoveLike(ctx, imageID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveLike", reflect.TypeOf((*MockimageUseCase)(nil).RemoveLike), ctx, imageID, userID)
}

// States mocks base method.
func (m *MockimageUseCase) States(ctx context.Context, imageID, userID int) (*domain.ImageStates, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "States", ctx, imageID, userID)
	ret0, _ := ret[0].(*domain.ImageStates)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// States indicates an expected call of States.
func (mr *MockimageUseCaseMockRecorder) States(ctx, imageID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "States", reflect.TypeOf((*MockimageUseCase)(nil).States), ctx, imageID, userID)
}

// Update mocks base method.
func (m *MockimageUseCase) Update(ctx context.Context, id int, image *domain.Image, executor *domain.User) (*domain.Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, image, executor)
	ret0, _ := ret[0].(*domain.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockimageUseCaseMockRecorder) Update(ctx, id, image, executor any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockimageUseCase)(nil).Update), ctx, id, image, executor)
}