package kafkax

import (
	"context"
	"errors"

	"github.com/boostgo/core/contextx"
	"github.com/boostgo/core/trace"

	"github.com/IBM/sarama"
)

// SyncProducer producer which produce messages in current goroutine
type SyncProducer struct {
	producer    sarama.SyncProducer
	traceMaster bool
}

// SyncProducerOption returns default sync producer configuration as [Option]
func SyncProducerOption() Option {
	return func(config *sarama.Config) {
		config.Producer.Return.Successes = true
		config.Producer.Return.Errors = true
		config.Producer.RequiredAcks = sarama.WaitForAll
		config.Producer.Compression = sarama.CompressionSnappy
		config.Producer.Retry.Max = 5
		config.Producer.Partitioner = sarama.NewRandomPartitioner
	}
}

// NewSyncProducer creates [SyncProducer] with configurations.
//
// Creates sync producer with default configuration as [Option] created by [SyncProducerOption] function.
//
// Adds producer close method to teardown
func NewSyncProducer(cfg Config, opts ...Option) (Producer, error) {
	config := sarama.NewConfig()
	config.ClientID = buildClientID()
	SyncProducerOption()(config)

	for _, opt := range opts {
		opt(config)
	}

	producer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}

	return &SyncProducer{
		producer:    producer,
		traceMaster: trace.AmIMaster(),
	}, nil
}

// NewSyncProducerFromClient creates [SyncProducer] by provided client.
//
// Creates sync producer with default configuration as Option created by SyncProducerOption function.
//
// Adds producer close method to teardown
func NewSyncProducerFromClient(client sarama.Client) (Producer, error) {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}

	return &SyncProducer{
		producer:    producer,
		traceMaster: trace.AmIMaster(),
	}, nil
}

// MustSyncProducer calls [NewSyncProducer] function with calls panic if returns error
func MustSyncProducer(cfg Config, opts ...Option) Producer {
	producer, err := NewSyncProducer(cfg, opts...)
	if err != nil {
		panic(err)
	}

	return producer
}

// MustSyncProducerFromClient calls [NewSyncProducerFromClient] function with calls panic if returns error
func MustSyncProducerFromClient(client sarama.Client) Producer {
	producer, err := NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}

	return producer
}

// Produce sends provided message(s) in the same goroutine.
//
// Sets trace id to provided messages to header
func (producer *SyncProducer) Produce(ctx context.Context, messages ...*sarama.ProducerMessage) error {
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

	if err := producer.producer.SendMessages(messages); err != nil {
		var pErrs sarama.ProducerErrors
		if ok := errors.As(err, &pErrs); ok {
			producerErrors := make([]error, 0, len(pErrs))
			for _, pErr := range pErrs {
				producerErrors = append(producerErrors, pErr.Err)
			}

			return ErrProduceMessages.
				SetError(producerErrors...).
				AddParam("size", len(pErrs))
		}

		return err
	}

	return nil
}
