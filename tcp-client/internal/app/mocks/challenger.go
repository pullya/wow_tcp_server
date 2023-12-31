// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Challenger is an autogenerated mock type for the Challenger type
type Challenger struct {
	mock.Mock
}

// GenerateSolution provides a mock function with given fields: ctx, challenge
func (_m *Challenger) GenerateSolution(ctx context.Context, challenge string) string {
	ret := _m.Called(ctx, challenge)

	if len(ret) == 0 {
		panic("no return value specified for GenerateSolution")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, challenge)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// SetDifficulty provides a mock function with given fields: diff
func (_m *Challenger) SetDifficulty(diff int) {
	_m.Called(diff)
}

// NewChallenger creates a new instance of Challenger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewChallenger(t interface {
	mock.TestingT
	Cleanup(func())
}) *Challenger {
	mock := &Challenger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
