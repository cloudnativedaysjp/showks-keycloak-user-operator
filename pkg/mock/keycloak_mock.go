// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/keycloak/keycloak.go

// Package mock_keycloak is a generated GoMock package.
package mock_keycloak

import (
	gocloak "github.com/Nerzal/gocloak"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockKeyCloakClientInterface is a mock of KeyCloakClientInterface interface
type MockKeyCloakClientInterface struct {
	ctrl     *gomock.Controller
	recorder *MockKeyCloakClientInterfaceMockRecorder
}

// MockKeyCloakClientInterfaceMockRecorder is the mock recorder for MockKeyCloakClientInterface
type MockKeyCloakClientInterfaceMockRecorder struct {
	mock *MockKeyCloakClientInterface
}

// NewMockKeyCloakClientInterface creates a new mock instance
func NewMockKeyCloakClientInterface(ctrl *gomock.Controller) *MockKeyCloakClientInterface {
	mock := &MockKeyCloakClientInterface{ctrl: ctrl}
	mock.recorder = &MockKeyCloakClientInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockKeyCloakClientInterface) EXPECT() *MockKeyCloakClientInterfaceMockRecorder {
	return m.recorder
}

// GetUsers mocks base method
func (m *MockKeyCloakClientInterface) GetUsers(realm string, param gocloak.GetUsersParams) (*[]gocloak.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", realm, param)
	ret0, _ := ret[0].(*[]gocloak.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers
func (mr *MockKeyCloakClientInterfaceMockRecorder) GetUsers(realm, param interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockKeyCloakClientInterface)(nil).GetUsers), realm, param)
}

// GetUserByID mocks base method
func (m *MockKeyCloakClientInterface) GetUserByID(realm, id string) (*gocloak.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", realm, id)
	ret0, _ := ret[0].(*gocloak.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID
func (mr *MockKeyCloakClientInterfaceMockRecorder) GetUserByID(realm, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockKeyCloakClientInterface)(nil).GetUserByID), realm, id)
}

// CreateUser mocks base method
func (m *MockKeyCloakClientInterface) CreateUser(realm string, user gocloak.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", realm, user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser
func (mr *MockKeyCloakClientInterfaceMockRecorder) CreateUser(realm, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockKeyCloakClientInterface)(nil).CreateUser), realm, user)
}

// DeleteUser mocks base method
func (m *MockKeyCloakClientInterface) DeleteUser(realm, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", realm, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser
func (mr *MockKeyCloakClientInterfaceMockRecorder) DeleteUser(realm, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockKeyCloakClientInterface)(nil).DeleteUser), realm, id)
}
