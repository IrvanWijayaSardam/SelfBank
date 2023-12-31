// Code generated by mockery v2.35.4. DO NOT EDIT.

package mocks

import (
	echo "github.com/labstack/echo/v4"
	mock "github.com/stretchr/testify/mock"
)

// DepositController is an autogenerated mock type for the DepositController type
type DepositController struct {
	mock.Mock
}

// All provides a mock function with given fields: context
func (_m *DepositController) All(context echo.Context) error {
	ret := _m.Called(context)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context) error); ok {
		r0 = rf(context)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindDepositByID provides a mock function with given fields: context
func (_m *DepositController) FindDepositByID(context echo.Context) error {
	ret := _m.Called(context)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context) error); ok {
		r0 = rf(context)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HandleMidtransNotification provides a mock function with given fields: context
func (_m *DepositController) HandleMidtransNotification(context echo.Context) error {
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
func (_m *DepositController) Insert(context echo.Context) error {
	ret := _m.Called(context)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context) error); ok {
		r0 = rf(context)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Refund provides a mock function with given fields: context
func (_m *DepositController) Refund(context echo.Context) error {
	ret := _m.Called(context)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context) error); ok {
		r0 = rf(context)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewDepositController creates a new instance of DepositController. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDepositController(t interface {
	mock.TestingT
	Cleanup(func())
}) *DepositController {
	mock := &DepositController{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
