// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// IStorage is an autogenerated mock type for the IStorage type
type IStorage struct {
	mock.Mock
}

// GetRandomWoW provides a mock function with given fields: ctx
func (_m *IStorage) GetRandomWoW(ctx context.Context) string {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetRandomWoW")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context) string); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewIStorage creates a new instance of IStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *IStorage {
	mock := &IStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}