package kafkax

import (
	"context"

	"github.com/boostgo/core/contextx"
	"github.com/boostgo/core/trace"

	"github.com/IBM/sarama"
)

// AsyncProducer producer which produce messages in "async" way
type AsyncProducer struct {
	producer    sarama.AsyncProducer
	traceMaster bool
}

// AsyncProducerOption returns default async producer configuration as [Option]
func AsyncProducerOption() Option {
	return func(config *sarama.Config) {
		config.Producer.Return.Successes = true
		config.Producer.Return.Errors = true
	}
}

// NewAsyncProducer creates [AsyncProducer] with configurations.
//
// Creates async producer with default configuration as [Option] created by [AsyncProducerOption] function.
//
// Adds producer close method to teardown
func NewAsyncProducer(cfg Config, opts ...Option) (Producer, error) {
	config := sarama.NewConfig()
	config.ClientID = buildClientID()
	AsyncProducerOption()(config)

	for _, opt := range opts {
		opt(config)
	}

	producer, err := sarama.NewAsyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}

	return &AsyncProducer{
		producer:    producer,
		traceMaster: trace.AmIMaster(),
	}, nil
}

// NewAsyncProducerFromClient creates [AsyncProducer] by provided client.
//
// Creates async producer with default configuration as [Option] created by [AsyncProducerOption] function.
//
// Adds producer close method to teardown
func NewAsyncProducerFromClient(client sarama.Client) (Producer, error) {
	producer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}

	return &AsyncProducer{
		producer:    producer,
		traceMaster: trace.AmIMaster(),
	}, nil
}

// MustAsyncProducer calls [NewAsyncProducer] function with calls panic if returns error
func MustAsyncProducer(cfg Config, opts ...Option) Producer {
	producer, err := NewAsyncProducer(cfg, opts...)
	if err != nil {
		panic(err)
	}

	return producer
}

// MustAsyncProducerFromClient calls [NewAsyncProducerFromClient] function with calls panic if returns error
func MustAsyncProducerFromClient(client sarama.Client) Producer {
	producer, err := NewAsyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}

	return producer
}

// Produce sends provided message(s) in other goroutine.
//
// Sets trace id to provided messages to header
func (producer *AsyncProducer) Produce(ctx context.Context, messages ...*sarama.ProducerMessage) error {
	if err := contextx.Validate(ctx); err != nil {
		return err
	}

	if len(messages) == 0 {
		return nil
	}

	traceID := trace.Get(ctx)
	if producer.traceMaster && traceID == "" {
		ctx = trace.SetID(ctx, traceID)
	}

	setTrace(ctx, messages...)

	for _, msg := range messages {
		producer.producer.Input() <- msg
	}

	return nil
}
