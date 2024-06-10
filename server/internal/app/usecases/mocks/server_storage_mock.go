// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/Nixonxp/discord/server/internal/app/models"
	mock "github.com/stretchr/testify/mock"
)

// ServerStorage is an autogenerated mock type for the ServerStorage type
type ServerStorage struct {
	mock.Mock
}

// CreateServer provides a mock function with given fields: ctx, server
func (_m *ServerStorage) CreateServer(ctx context.Context, server models.ServerInfo) error {
	ret := _m.Called(ctx, server)

	if len(ret) == 0 {
		panic("no return value specified for CreateServer")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.ServerInfo) error); ok {
		r0 = rf(ctx, server)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetServerById provides a mock function with given fields: ctx, id
func (_m *ServerStorage) GetServerById(ctx context.Context, id string) (*models.ServerInfo, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetServerById")
	}

	var r0 *models.ServerInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*models.ServerInfo, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.ServerInfo); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ServerInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchServers provides a mock function with given fields: ctx, serverName
func (_m *ServerStorage) SearchServers(ctx context.Context, serverName string) ([]*models.ServerInfo, error) {
	ret := _m.Called(ctx, serverName)

	if len(ret) == 0 {
		panic("no return value specified for SearchServers")
	}

	var r0 []*models.ServerInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]*models.ServerInfo, error)); ok {
		return rf(ctx, serverName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []*models.ServerInfo); ok {
		r0 = rf(ctx, serverName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.ServerInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, serverName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewServerStorage creates a new instance of ServerStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewServerStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *ServerStorage {
	mock := &ServerStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
