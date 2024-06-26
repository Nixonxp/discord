// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/Nixonxp/discord/server/internal/app/models"
	mock "github.com/stretchr/testify/mock"
)

// SubscribeStorage is an autogenerated mock type for the SubscribeStorage type
type SubscribeStorage struct {
	mock.Mock
}

// CreateSubscribe provides a mock function with given fields: ctx, server
func (_m *SubscribeStorage) CreateSubscribe(ctx context.Context, server models.SubscribeInfo) error {
	ret := _m.Called(ctx, server)

	if len(ret) == 0 {
		panic("no return value specified for CreateSubscribe")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.SubscribeInfo) error); ok {
		r0 = rf(ctx, server)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteSubscribe provides a mock function with given fields: ctx, serverId, userId
func (_m *SubscribeStorage) DeleteSubscribe(ctx context.Context, serverId models.ServerID, userId models.UserID) error {
	ret := _m.Called(ctx, serverId, userId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteSubscribe")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.ServerID, models.UserID) error); ok {
		r0 = rf(ctx, serverId, userId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByUserId provides a mock function with given fields: ctx, userId
func (_m *SubscribeStorage) GetByUserId(ctx context.Context, userId models.UserID) ([]*models.SubscribeInfo, error) {
	ret := _m.Called(ctx, userId)

	if len(ret) == 0 {
		panic("no return value specified for GetByUserId")
	}

	var r0 []*models.SubscribeInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.UserID) ([]*models.SubscribeInfo, error)); ok {
		return rf(ctx, userId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.UserID) []*models.SubscribeInfo); ok {
		r0 = rf(ctx, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.SubscribeInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.UserID) error); ok {
		r1 = rf(ctx, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSubscribeStorage creates a new instance of SubscribeStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSubscribeStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *SubscribeStorage {
	mock := &SubscribeStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
