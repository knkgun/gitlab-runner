// Code generated by mockery v1.1.0. DO NOT EDIT.

package referees

import (
	bytes "bytes"
	context "context"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// MockReferee is an autogenerated mock type for the Referee type
type MockReferee struct {
	mock.Mock
}

// ArtifactBaseName provides a mock function with given fields:
func (_m *MockReferee) ArtifactBaseName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ArtifactFormat provides a mock function with given fields:
func (_m *MockReferee) ArtifactFormat() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ArtifactType provides a mock function with given fields:
func (_m *MockReferee) ArtifactType() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Execute provides a mock function with given fields: ctx, startTime, endTime
func (_m *MockReferee) Execute(ctx context.Context, startTime time.Time, endTime time.Time) (*bytes.Reader, error) {
	ret := _m.Called(ctx, startTime, endTime)

	var r0 *bytes.Reader
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, time.Time) *bytes.Reader); ok {
		r0 = rf(ctx, startTime, endTime)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bytes.Reader)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, time.Time, time.Time) error); ok {
		r1 = rf(ctx, startTime, endTime)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
