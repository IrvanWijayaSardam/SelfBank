// Code generated by mockery v2.35.4. DO NOT EDIT.

package mocks

import (
	echo "github.com/labstack/echo/v4"
	mock "github.com/stretchr/testify/mock"
)

// WithdrawalController is an autogenerated mock type for the WithdrawalController type
type WithdrawalController struct {
	mock.Mock
}

// All provides a mock function with given fields: context
func (_m *WithdrawalController) All(context echo.Context) error {
	ret := _m.Called(context)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context) error); ok {
		r0 = rf(context)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindWithdrawalByID provides a mock function with given fields: context
func (_m *WithdrawalController) FindWithdrawalByID(context echo.Context) error {
	ret := _m.Called(context)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context) error); ok {
		r0 = rf(context)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Insert provides a mock function with given fields: context
func (_m *WithdrawalController) Insert(context echo.Context) error {
	ret := _m.Called(context)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context) error); ok {
		r0 = rf(context)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewWithdrawalController creates a new instance of WithdrawalController. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWithdrawalController(t interface {
	mock.TestingT
	Cleanup(func())
}) *WithdrawalController {
	mock := &WithdrawalController{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}