// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import mock "github.com/stretchr/testify/mock"
import resource "github.com/kyma-project/kyma/components/console-backend-service/pkg/resource"
import v1alpha2 "github.com/kyma-project/kyma/components/api-controller/pkg/apis/gateway.kyma-project.io/v1alpha2"

// apiSvc is an autogenerated mock type for the apiSvc type
type apiSvc struct {
	mock.Mock
}

// Create provides a mock function with given fields: api
func (_m *apiSvc) Create(api *v1alpha2.Api) (*v1alpha2.Api, error) {
	ret := _m.Called(api)

	var r0 *v1alpha2.Api
	if rf, ok := ret.Get(0).(func(*v1alpha2.Api) *v1alpha2.Api); ok {
		r0 = rf(api)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1alpha2.Api)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*v1alpha2.Api) error); ok {
		r1 = rf(api)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: name, namespace
func (_m *apiSvc) Delete(name string, namespace string) error {
	ret := _m.Called(name, namespace)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(name, namespace)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: name, namespace
func (_m *apiSvc) Find(name string, namespace string) (*v1alpha2.Api, error) {
	ret := _m.Called(name, namespace)

	var r0 *v1alpha2.Api
	if rf, ok := ret.Get(0).(func(string, string) *v1alpha2.Api); ok {
		r0 = rf(name, namespace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1alpha2.Api)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(name, namespace)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: namespace, serviceName, hostname
func (_m *apiSvc) List(namespace string, serviceName *string, hostname *string) ([]*v1alpha2.Api, error) {
	ret := _m.Called(namespace, serviceName, hostname)

	var r0 []*v1alpha2.Api
	if rf, ok := ret.Get(0).(func(string, *string, *string) []*v1alpha2.Api); ok {
		r0 = rf(namespace, serviceName, hostname)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*v1alpha2.Api)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *string, *string) error); ok {
		r1 = rf(namespace, serviceName, hostname)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Subscribe provides a mock function with given fields: listener
func (_m *apiSvc) Subscribe(listener resource.Listener) {
	_m.Called(listener)
}

// Unsubscribe provides a mock function with given fields: listener
func (_m *apiSvc) Unsubscribe(listener resource.Listener) {
	_m.Called(listener)
}

// Update provides a mock function with given fields: api
func (_m *apiSvc) Update(api *v1alpha2.Api) (*v1alpha2.Api, error) {
	ret := _m.Called(api)

	var r0 *v1alpha2.Api
	if rf, ok := ret.Get(0).(func(*v1alpha2.Api) *v1alpha2.Api); ok {
		r0 = rf(api)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1alpha2.Api)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*v1alpha2.Api) error); ok {
		r1 = rf(api)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
