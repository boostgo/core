package kafkax

import (
	"context"

	"github.com/IBM/sarama"
)

type Producer interface {
	Produce(ctx context.Context, messages ...*sarama.ProducerMessage) error
}

type mockProducer struct {
	handleMessages func(ctx context.Context, messages ...*sarama.ProducerMessage) error
}

func MockProducer(handle ...func(ctx context.Context, messages ...*sarama.ProducerMessage) error) Producer {
	var handleMessages func(ctx context.Context, messages ...*sarama.ProducerMessage) error
	if len(handle) > 0 {
		handleMessages = handle[0]
	}

	return &mockProducer{
		handleMessages: handleMessages,
	}
}

func (mp *mockProducer) Produce(ctx context.Context, messages ...*sarama.ProducerMessage) error {
	if mp.handleMessages != nil {
		return mp.handleMessages(ctx, messages...)
	}

	return nil
}
