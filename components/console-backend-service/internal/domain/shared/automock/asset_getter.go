// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	mock "github.com/stretchr/testify/mock"

	v1beta1 "github.com/kyma-project/rafter/pkg/apis/rafter/v1beta1"
)

// AssetGetter is an autogenerated mock type for the AssetGetter type
type AssetGetter struct {
	mock.Mock
}

// ListForAssetGroupByType provides a mock function with given fields: namespace, assetGroupName, types
func (_m *AssetGetter) ListForAssetGroupByType(namespace string, assetGroupName string, types []string) ([]*v1beta1.Asset, error) {
	ret := _m.Called(namespace, assetGroupName, types)

	var r0 []*v1beta1.Asset
	if rf, ok := ret.Get(0).(func(string, string, []string) []*v1beta1.Asset); ok {
		r0 = rf(namespace, assetGroupName, types)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*v1beta1.Asset)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, []string) error); ok {
		r1 = rf(namespace, assetGroupName, types)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
