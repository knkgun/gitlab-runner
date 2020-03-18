// Code generated by mockery v1.0.0. DO NOT EDIT.

package scope

import mock "github.com/stretchr/testify/mock"

// MockEntriesMapBuilder is an autogenerated mock type for the EntriesMapBuilder type
type MockEntriesMapBuilder struct {
	mock.Mock
}

// Build provides a mock function with given fields: objects
func (_m *MockEntriesMapBuilder) Build(objects []Object) (EntriesMap, error) {
	ret := _m.Called(objects)

	var r0 EntriesMap
	if rf, ok := ret.Get(0).(func([]Object) EntriesMap); ok {
		r0 = rf(objects)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(EntriesMap)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]Object) error); ok {
		r1 = rf(objects)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}