// Code generated by mockery v1.0.0. DO NOT EDIT.
package mocks

import mock "github.com/stretchr/testify/mock"

// Connector is an autogenerated mock type for the Connector type
type Connector struct {
	mock.Mock
}

// EstablishConnection provides a mock function with given fields:
func (_m *Connector) EstablishConnection() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
