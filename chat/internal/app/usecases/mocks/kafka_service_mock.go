// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	usecases "github.com/Nixonxp/discord/chat/internal/app/usecases"
	mock "github.com/stretchr/testify/mock"
)

// KafkaServiceInterface is an autogenerated mock type for the KafkaServiceInterface type
type KafkaServiceInterface struct {
	mock.Mock
}

// SendMessage provides a mock function with given fields: msgData
func (_m *KafkaServiceInterface) SendMessage(msgData usecases.MessageDto) error {
	ret := _m.Called(msgData)

	if len(ret) == 0 {
		panic("no return value specified for SendMessage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(usecases.MessageDto) error); ok {
		r0 = rf(msgData)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewKafkaServiceInterface creates a new instance of KafkaServiceInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewKafkaServiceInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *KafkaServiceInterface {
	mock := &KafkaServiceInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
