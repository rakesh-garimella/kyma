// Code generated by mockery v1.0.0. DO NOT EDIT.
package mocks

import mock "github.com/stretchr/testify/mock"

// Supervisor is an autogenerated mock type for the Supervisor type
type Supervisor struct {
	mock.Mock
}

// InitializeCompassConnectionCR provides a mock function with given fields:
func (_m *Supervisor) InitializeCompassConnectionCR() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
