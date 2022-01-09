// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DuckLuckBreakout/ozonBackend/internal/pkg/user (interfaces: UseCase)

// Package mock is a generated GoMock package.
package mock

import (
	multipart "mime/multipart"
	reflect "reflect"

	models "github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models"
	gomock "github.com/golang/mock/gomock"
)

// MockUseCase is a mock of UseCase interface.
type MockUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUseCaseMockRecorder
}

// MockUseCaseMockRecorder is the mock recorder for MockUseCase.
type MockUseCaseMockRecorder struct {
	mock *MockUseCase
}

// NewMockUseCase creates a new mock instance.
func NewMockUseCase(ctrl *gomock.Controller) *MockUseCase {
	mock := &MockUseCase{ctrl: ctrl}
	mock.recorder = &MockUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUseCase) EXPECT() *MockUseCaseMockRecorder {
	return m.recorder
}

// AddUser mocks base method.
func (m *MockUseCase) AddUser(arg0 *models.SignupUser) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", arg0)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUser indicates an expected call of AddUser.
func (mr *MockUseCaseMockRecorder) AddUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockUseCase)(nil).AddUser), arg0)
}

// Authorize mocks base method.
func (m *MockUseCase) Authorize(arg0 *models.LoginUser) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authorize", arg0)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Authorize indicates an expected call of Authorize.
func (mr *MockUseCaseMockRecorder) Authorize(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authorize", reflect.TypeOf((*MockUseCase)(nil).Authorize), arg0)
}

// GetAvatar mocks base method.
func (m *MockUseCase) GetAvatar(arg0 uint64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAvatar", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAvatar indicates an expected call of GetAvatar.
func (mr *MockUseCaseMockRecorder) GetAvatar(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvatar", reflect.TypeOf((*MockUseCase)(nil).GetAvatar), arg0)
}

// GetUserById mocks base method.
func (m *MockUseCase) GetUserById(arg0 uint64) (*models.ProfileUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserById", arg0)
	ret0, _ := ret[0].(*models.ProfileUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserById indicates an expected call of GetUserById.
func (mr *MockUseCaseMockRecorder) GetUserById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserById", reflect.TypeOf((*MockUseCase)(nil).GetUserById), arg0)
}

// SetAvatar mocks base method.
func (m *MockUseCase) SetAvatar(arg0 uint64, arg1 *multipart.File, arg2 *multipart.FileHeader) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetAvatar", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetAvatar indicates an expected call of SetAvatar.
func (mr *MockUseCaseMockRecorder) SetAvatar(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAvatar", reflect.TypeOf((*MockUseCase)(nil).SetAvatar), arg0, arg1, arg2)
}

// UpdateProfile mocks base method.
func (m *MockUseCase) UpdateProfile(arg0 uint64, arg1 *models.UpdateUser) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProfile", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProfile indicates an expected call of UpdateProfile.
func (mr *MockUseCaseMockRecorder) UpdateProfile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfile", reflect.TypeOf((*MockUseCase)(nil).UpdateProfile), arg0, arg1)
}
