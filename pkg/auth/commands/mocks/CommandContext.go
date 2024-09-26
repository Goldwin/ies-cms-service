// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	repositories "github.com/Goldwin/ies-pik-cms/pkg/auth/repositories"
	mock "github.com/stretchr/testify/mock"
)

// CommandContext is an autogenerated mock type for the CommandContext type
type CommandContext struct {
	mock.Mock
}

type CommandContext_Expecter struct {
	mock *mock.Mock
}

func (_m *CommandContext) EXPECT() *CommandContext_Expecter {
	return &CommandContext_Expecter{mock: &_m.Mock}
}

// AccountRepository provides a mock function with given fields:
func (_m *CommandContext) AccountRepository() repositories.AccountRepository {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for AccountRepository")
	}

	var r0 repositories.AccountRepository
	if rf, ok := ret.Get(0).(func() repositories.AccountRepository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repositories.AccountRepository)
		}
	}

	return r0
}

// CommandContext_AccountRepository_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AccountRepository'
type CommandContext_AccountRepository_Call struct {
	*mock.Call
}

// AccountRepository is a helper method to define mock.On call
func (_e *CommandContext_Expecter) AccountRepository() *CommandContext_AccountRepository_Call {
	return &CommandContext_AccountRepository_Call{Call: _e.mock.On("AccountRepository")}
}

func (_c *CommandContext_AccountRepository_Call) Run(run func()) *CommandContext_AccountRepository_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *CommandContext_AccountRepository_Call) Return(_a0 repositories.AccountRepository) *CommandContext_AccountRepository_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CommandContext_AccountRepository_Call) RunAndReturn(run func() repositories.AccountRepository) *CommandContext_AccountRepository_Call {
	_c.Call.Return(run)
	return _c
}

// OtpRepository provides a mock function with given fields:
func (_m *CommandContext) OtpRepository() repositories.OtpRepository {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for OtpRepository")
	}

	var r0 repositories.OtpRepository
	if rf, ok := ret.Get(0).(func() repositories.OtpRepository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repositories.OtpRepository)
		}
	}

	return r0
}

// CommandContext_OtpRepository_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OtpRepository'
type CommandContext_OtpRepository_Call struct {
	*mock.Call
}

// OtpRepository is a helper method to define mock.On call
func (_e *CommandContext_Expecter) OtpRepository() *CommandContext_OtpRepository_Call {
	return &CommandContext_OtpRepository_Call{Call: _e.mock.On("OtpRepository")}
}

func (_c *CommandContext_OtpRepository_Call) Run(run func()) *CommandContext_OtpRepository_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *CommandContext_OtpRepository_Call) Return(_a0 repositories.OtpRepository) *CommandContext_OtpRepository_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CommandContext_OtpRepository_Call) RunAndReturn(run func() repositories.OtpRepository) *CommandContext_OtpRepository_Call {
	_c.Call.Return(run)
	return _c
}

// PasswordRepository provides a mock function with given fields:
func (_m *CommandContext) PasswordRepository() repositories.PasswordRepository {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for PasswordRepository")
	}

	var r0 repositories.PasswordRepository
	if rf, ok := ret.Get(0).(func() repositories.PasswordRepository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repositories.PasswordRepository)
		}
	}

	return r0
}

// CommandContext_PasswordRepository_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PasswordRepository'
type CommandContext_PasswordRepository_Call struct {
	*mock.Call
}

// PasswordRepository is a helper method to define mock.On call
func (_e *CommandContext_Expecter) PasswordRepository() *CommandContext_PasswordRepository_Call {
	return &CommandContext_PasswordRepository_Call{Call: _e.mock.On("PasswordRepository")}
}

func (_c *CommandContext_PasswordRepository_Call) Run(run func()) *CommandContext_PasswordRepository_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *CommandContext_PasswordRepository_Call) Return(_a0 repositories.PasswordRepository) *CommandContext_PasswordRepository_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CommandContext_PasswordRepository_Call) RunAndReturn(run func() repositories.PasswordRepository) *CommandContext_PasswordRepository_Call {
	_c.Call.Return(run)
	return _c
}

// PasswordResetCodeRepository provides a mock function with given fields:
func (_m *CommandContext) PasswordResetCodeRepository() repositories.PasswordResetCodeRepository {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for PasswordResetCodeRepository")
	}

	var r0 repositories.PasswordResetCodeRepository
	if rf, ok := ret.Get(0).(func() repositories.PasswordResetCodeRepository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repositories.PasswordResetCodeRepository)
		}
	}

	return r0
}

// CommandContext_PasswordResetCodeRepository_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PasswordResetCodeRepository'
type CommandContext_PasswordResetCodeRepository_Call struct {
	*mock.Call
}

// PasswordResetCodeRepository is a helper method to define mock.On call
func (_e *CommandContext_Expecter) PasswordResetCodeRepository() *CommandContext_PasswordResetCodeRepository_Call {
	return &CommandContext_PasswordResetCodeRepository_Call{Call: _e.mock.On("PasswordResetCodeRepository")}
}

func (_c *CommandContext_PasswordResetCodeRepository_Call) Run(run func()) *CommandContext_PasswordResetCodeRepository_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *CommandContext_PasswordResetCodeRepository_Call) Return(_a0 repositories.PasswordResetCodeRepository) *CommandContext_PasswordResetCodeRepository_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CommandContext_PasswordResetCodeRepository_Call) RunAndReturn(run func() repositories.PasswordResetCodeRepository) *CommandContext_PasswordResetCodeRepository_Call {
	_c.Call.Return(run)
	return _c
}

// NewCommandContext creates a new instance of CommandContext. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCommandContext(t interface {
	mock.TestingT
	Cleanup(func())
}) *CommandContext {
	mock := &CommandContext{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}