// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	entities "github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	mock "github.com/stretchr/testify/mock"
)

// PasswordResetCodeRepository is an autogenerated mock type for the PasswordResetCodeRepository type
type PasswordResetCodeRepository struct {
	mock.Mock
}

type PasswordResetCodeRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *PasswordResetCodeRepository) EXPECT() *PasswordResetCodeRepository_Expecter {
	return &PasswordResetCodeRepository_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: _a0
func (_m *PasswordResetCodeRepository) Delete(_a0 *entities.PasswordResetCode) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*entities.PasswordResetCode) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PasswordResetCodeRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type PasswordResetCodeRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - _a0 *entities.PasswordResetCode
func (_e *PasswordResetCodeRepository_Expecter) Delete(_a0 interface{}) *PasswordResetCodeRepository_Delete_Call {
	return &PasswordResetCodeRepository_Delete_Call{Call: _e.mock.On("Delete", _a0)}
}

func (_c *PasswordResetCodeRepository_Delete_Call) Run(run func(_a0 *entities.PasswordResetCode)) *PasswordResetCodeRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*entities.PasswordResetCode))
	})
	return _c
}

func (_c *PasswordResetCodeRepository_Delete_Call) Return(_a0 error) *PasswordResetCodeRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PasswordResetCodeRepository_Delete_Call) RunAndReturn(run func(*entities.PasswordResetCode) error) *PasswordResetCodeRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: _a0
func (_m *PasswordResetCodeRepository) Get(_a0 string) (*entities.PasswordResetCode, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *entities.PasswordResetCode
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*entities.PasswordResetCode, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) *entities.PasswordResetCode); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.PasswordResetCode)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PasswordResetCodeRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type PasswordResetCodeRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - _a0 string
func (_e *PasswordResetCodeRepository_Expecter) Get(_a0 interface{}) *PasswordResetCodeRepository_Get_Call {
	return &PasswordResetCodeRepository_Get_Call{Call: _e.mock.On("Get", _a0)}
}

func (_c *PasswordResetCodeRepository_Get_Call) Run(run func(_a0 string)) *PasswordResetCodeRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *PasswordResetCodeRepository_Get_Call) Return(_a0 *entities.PasswordResetCode, _a1 error) *PasswordResetCodeRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PasswordResetCodeRepository_Get_Call) RunAndReturn(run func(string) (*entities.PasswordResetCode, error)) *PasswordResetCodeRepository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: _a0
func (_m *PasswordResetCodeRepository) List(_a0 []string) ([]*entities.PasswordResetCode, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []*entities.PasswordResetCode
	var r1 error
	if rf, ok := ret.Get(0).(func([]string) ([]*entities.PasswordResetCode, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func([]string) []*entities.PasswordResetCode); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entities.PasswordResetCode)
		}
	}

	if rf, ok := ret.Get(1).(func([]string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PasswordResetCodeRepository_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type PasswordResetCodeRepository_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - _a0 []string
func (_e *PasswordResetCodeRepository_Expecter) List(_a0 interface{}) *PasswordResetCodeRepository_List_Call {
	return &PasswordResetCodeRepository_List_Call{Call: _e.mock.On("List", _a0)}
}

func (_c *PasswordResetCodeRepository_List_Call) Run(run func(_a0 []string)) *PasswordResetCodeRepository_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]string))
	})
	return _c
}

func (_c *PasswordResetCodeRepository_List_Call) Return(_a0 []*entities.PasswordResetCode, _a1 error) *PasswordResetCodeRepository_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PasswordResetCodeRepository_List_Call) RunAndReturn(run func([]string) ([]*entities.PasswordResetCode, error)) *PasswordResetCodeRepository_List_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields: _a0
func (_m *PasswordResetCodeRepository) Save(_a0 *entities.PasswordResetCode) (*entities.PasswordResetCode, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 *entities.PasswordResetCode
	var r1 error
	if rf, ok := ret.Get(0).(func(*entities.PasswordResetCode) (*entities.PasswordResetCode, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(*entities.PasswordResetCode) *entities.PasswordResetCode); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.PasswordResetCode)
		}
	}

	if rf, ok := ret.Get(1).(func(*entities.PasswordResetCode) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PasswordResetCodeRepository_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type PasswordResetCodeRepository_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - _a0 *entities.PasswordResetCode
func (_e *PasswordResetCodeRepository_Expecter) Save(_a0 interface{}) *PasswordResetCodeRepository_Save_Call {
	return &PasswordResetCodeRepository_Save_Call{Call: _e.mock.On("Save", _a0)}
}

func (_c *PasswordResetCodeRepository_Save_Call) Run(run func(_a0 *entities.PasswordResetCode)) *PasswordResetCodeRepository_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*entities.PasswordResetCode))
	})
	return _c
}

func (_c *PasswordResetCodeRepository_Save_Call) Return(_a0 *entities.PasswordResetCode, _a1 error) *PasswordResetCodeRepository_Save_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PasswordResetCodeRepository_Save_Call) RunAndReturn(run func(*entities.PasswordResetCode) (*entities.PasswordResetCode, error)) *PasswordResetCodeRepository_Save_Call {
	_c.Call.Return(run)
	return _c
}

// NewPasswordResetCodeRepository creates a new instance of PasswordResetCodeRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPasswordResetCodeRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *PasswordResetCodeRepository {
	mock := &PasswordResetCodeRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
