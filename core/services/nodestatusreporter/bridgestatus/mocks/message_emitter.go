package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
)

// MessageEmitter is a mock implementation of custmsg.MessageEmitter for testing
type MessageEmitter struct {
	mock.Mock
}

func (m *MessageEmitter) With(kvs ...string) custmsg.MessageEmitter {
	args := m.Called(kvs)
	return args.Get(0).(custmsg.MessageEmitter)
}

func (m *MessageEmitter) WithMapLabels(labels map[string]string) custmsg.MessageEmitter {
	args := m.Called(labels)
	return args.Get(0).(custmsg.MessageEmitter)
}

func (m *MessageEmitter) Emit(ctx context.Context, msg string) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *MessageEmitter) Labels() map[string]string {
	args := m.Called()
	return args.Get(0).(map[string]string)
}

// NewMessageEmitter creates a new message emitter mock for testing
func NewMessageEmitter() *MessageEmitter {
	return &MessageEmitter{}
}
