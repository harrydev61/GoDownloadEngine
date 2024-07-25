package common

import "context"

type MqHandlerFunc func(ctx context.Context, queueName string, payload []byte) error

type MqConfig struct {
	Brokers  string
	Version  string
	Group    string
	Assignor string
	Oldest   bool
	Verbose  bool
	ClientId string
}
