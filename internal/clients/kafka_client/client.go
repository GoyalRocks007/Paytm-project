package kafkaclient

import (
	"context"
)

type IKafkaClient interface {
	Listen(ctx context.Context, handler func(key, value []byte))
	Stop() error
}
