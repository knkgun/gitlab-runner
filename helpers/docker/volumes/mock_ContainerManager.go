// Code generated by mockery v1.0.0. DO NOT EDIT.

package volumes

import container "github.com/docker/docker/api/types/container"
import mock "github.com/stretchr/testify/mock"
import network "github.com/docker/docker/api/types/network"
import types "github.com/docker/docker/api/types"

// MockContainerManager is an autogenerated mock type for the ContainerManager type
type MockContainerManager struct {
	mock.Mock
}

// CreateContainer provides a mock function with given fields: config, hostConfig, networkingConfig, containerName
func (_m *MockContainerManager) CreateContainer(config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.ContainerCreateCreatedBody, error) {
	ret := _m.Called(config, hostConfig, networkingConfig, containerName)

	var r0 container.ContainerCreateCreatedBody
	if rf, ok := ret.Get(0).(func(*container.Config, *container.HostConfig, *network.NetworkingConfig, string) container.ContainerCreateCreatedBody); ok {
		r0 = rf(config, hostConfig, networkingConfig, containerName)
	} else {
		r0 = ret.Get(0).(container.ContainerCreateCreatedBody)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*container.Config, *container.HostConfig, *network.NetworkingConfig, string) error); ok {
		r1 = rf(config, hostConfig, networkingConfig, containerName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InspectContainer provides a mock function with given fields: containerName
func (_m *MockContainerManager) InspectContainer(containerName string) (types.ContainerJSON, error) {
	ret := _m.Called(containerName)

	var r0 types.ContainerJSON
	if rf, ok := ret.Get(0).(func(string) types.ContainerJSON); ok {
		r0 = rf(containerName)
	} else {
		r0 = ret.Get(0).(types.ContainerJSON)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(containerName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LabelContainer provides a mock function with given fields: _a0, containerType, otherLabels
func (_m *MockContainerManager) LabelContainer(_a0 *container.Config, containerType string, otherLabels ...string) {
	_va := make([]interface{}, len(otherLabels))
	for _i := range otherLabels {
		_va[_i] = otherLabels[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0, containerType)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// RemoveContainer provides a mock function with given fields: id
func (_m *MockContainerManager) RemoveContainer(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StartContainer provides a mock function with given fields: containerID, options
func (_m *MockContainerManager) StartContainer(containerID string, options types.ContainerStartOptions) error {
	ret := _m.Called(containerID, options)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, types.ContainerStartOptions) error); ok {
		r0 = rf(containerID, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WaitForContainer provides a mock function with given fields: id
func (_m *MockContainerManager) WaitForContainer(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}