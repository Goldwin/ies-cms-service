// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	entities "github.com/Goldwin/ies-pik-cms/pkg/auth/entities"
	mock "github.com/stretchr/testify/mock"
)

// AccountRepository is an autogenerated mock type for the AccountRepository type
type AccountRepository struct {
	mock.Mock
}

type AccountRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *AccountRepository) EXPECT() *AccountRepository_Expecter {
	return &AccountRepository_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: _a0
func (_m *AccountRepository) Delete(_a0 *entities.Account) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*entities.Account) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AccountRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type AccountRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - _a0 *entities.Account
func (_e *AccountRepository_Expecter) Delete(_a0 interface{}) *AccountRepository_Delete_Call {
	return &AccountRepository_Delete_Call{Call: _e.mock.On("Delete", _a0)}
}

func (_c *AccountRepository_Delete_Call) Run(run func(_a0 *entities.Account)) *AccountRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*entities.Account))
	})
	return _c
}

func (_c *AccountRepository_Delete_Call) Return(_a0 error) *AccountRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AccountRepository_Delete_Call) RunAndReturn(run func(*entities.Account) error) *AccountRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: _a0
func (_m *AccountRepository) Get(_a0 string) (*entities.Account, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *entities.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*entities.Account, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) *entities.Account); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AccountRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type AccountRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - _a0 string
func (_e *AccountRepository_Expecter) Get(_a0 interface{}) *AccountRepository_Get_Call {
	return &AccountRepository_Get_Call{Call: _e.mock.On("Get", _a0)}
}

func (_c *AccountRepository_Get_Call) Run(run func(_a0 string)) *AccountRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *AccountRepository_Get_Call) Return(_a0 *entities.Account, _a1 error) *AccountRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AccountRepository_Get_Call) RunAndReturn(run func(string) (*entities.Account, error)) *AccountRepository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: _a0
func (_m *AccountRepository) List(_a0 []string) ([]*entities.Account, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []*entities.Account
	var r1 error
	if rf, ok := ret.Get(0).(func([]string) ([]*entities.Account, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func([]string) []*entities.Account); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entities.Account)
		}
	}

	if rf, ok := ret.Get(1).(func([]string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AccountRepository_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type AccountRepository_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - _a0 []string
func (_e *AccountRepository_Expecter) List(_a0 interface{}) *AccountRepository_List_Call {
	return &AccountRepository_List_Call{Call: _e.mock.On("List", _a0)}
}

func (_c *AccountRepository_List_Call) Run(run func(_a0 []string)) *AccountRepository_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]string))
	})
	return _c
}

func (_c *AccountRepository_List_Call) Return(_a0 []*entities.Account, _a1 error) *AccountRepository_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AccountRepository_List_Call) RunAndReturn(run func([]string) ([]*entities.Account, error)) *AccountRepository_List_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields: _a0
func (_m *AccountRepository) Save(_a0 *entities.Account) (*entities.Account, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 *entities.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(*entities.Account) (*entities.Account, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(*entities.Account) *entities.Account); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(*entities.Account) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AccountRepository_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type AccountRepository_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - _a0 *entities.Account
func (_e *AccountRepository_Expecter) Save(_a0 interface{}) *AccountRepository_Save_Call {
	return &AccountRepository_Save_Call{Call: _e.mock.On("Save", _a0)}
}

func (_c *AccountRepository_Save_Call) Run(run func(_a0 *entities.Account)) *AccountRepository_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*entities.Account))
	})
	return _c
}

func (_c *AccountRepository_Save_Call) Return(_a0 *entities.Account, _a1 error) *AccountRepository_Save_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AccountRepository_Save_Call) RunAndReturn(run func(*entities.Account) (*entities.Account, error)) *AccountRepository_Save_Call {
	_c.Call.Return(run)
	return _c
}

// NewAccountRepository creates a new instance of AccountRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAccountRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *AccountRepository {
	mock := &AccountRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
