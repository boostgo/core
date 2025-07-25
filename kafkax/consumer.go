package kafkax

import (
	"context"
	"time"

	"github.com/boostgo/core/appx"
	"github.com/boostgo/core/errorx"
	"github.com/boostgo/core/log"

	"github.com/IBM/sarama"
)

type ConsumeHandler func(ctx context.Context, message *sarama.ConsumerMessage) error
type ErrorHandler func(err error)

// Consumer wrap structure for [Consumer]
type Consumer struct {
	consumer     sarama.Consumer
	errorHandler ErrorHandler
}

// ConsumerOption returns default consumer configs
func ConsumerOption() Option {
	return func(config *sarama.Config) {
		config.Consumer.Return.Errors = true
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
		config.Consumer.Offsets.AutoCommit.Enable = true
		config.Consumer.Offsets.AutoCommit.Interval = time.Second

		config.Consumer.Fetch.Default = 1 << 20 // 1MB
		config.Consumer.Fetch.Max = 10 << 20    // 10MB
		config.ChannelBufferSize = 256
	}
}

// NewConsumer creates [Consumer] by options
func NewConsumer(cfg Config, opts ...Option) (*Consumer, error) {
	if err := validateConsumerConfig(cfg); err != nil {
		return nil, err
	}

	config := sarama.NewConfig()
	config.ClientID = buildClientID()
	ConsumerOption()(config)

	for _, opt := range opts {
		opt(config)
	}

	consumer, err := sarama.NewConsumer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: consumer,
	}, nil
}

// NewConsumerFromClient creates [Consumer] by sarama client
func NewConsumerFromClient(client sarama.Client) (*Consumer, error) {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: consumer,
	}, nil
}

// MustConsumer calls [NewConsumer] and if error catch throws panic
func MustConsumer(cfg Config, opts ...Option) *Consumer {
	consumer, err := NewConsumer(cfg, opts...)
	if err != nil {
		panic(err)
	}

	return consumer
}

// MustConsumerFromClient calls [NewConsumerFromClient] and if error catch throws panic
func MustConsumerFromClient(client sarama.Client) *Consumer {
	consumer, err := NewConsumerFromClient(client)
	if err != nil {
		panic(err)
	}

	return consumer
}

func (consumer *Consumer) SetErrorHandler(handler ErrorHandler) {
	consumer.errorHandler = handler
}

// Consume starts consuming topic with consumer.
//
// Catch consumer errors and provided context done (for graceful shutdown).
func (consumer *Consumer) Consume(topic string, handler ConsumeHandler) error {
	partitions, err := consumer.consumer.Partitions(topic)
	if err != nil {
		return err
	}

	for i := 0; i < len(partitions); i++ {
		var partitionConsumer sarama.PartitionConsumer
		partitionConsumer, err = consumer.consumer.ConsumePartition(topic, partitions[i], sarama.OffsetNewest)
		if err != nil {
			return err
		}

		appx.Tear(partitionConsumer.Close)

		go func(partition int32) {
			for {
				select {
				case <-appx.Context().Done():
					log.
						Info().
						Int32("partition", partition).
						Msg("Stop consumer by context")
					return
				case msg, ok := <-partitionConsumer.Messages():
					if !ok {
						log.
							Info().
							Int32("partition", partition).
							Msg("Stop consumer by closing channel")
						return
					}

					traceID := Header(msg, TraceKey)
					ctx := context.WithValue(context.Background(), TraceKey, traceID)

					if err = errorx.Try(func() error {
						return handler(ctx, msg)
					}); err != nil {
						log.
							Error().
							Ctx(ctx).
							Err(err).
							Msg("Handle message error")

						if consumer.errorHandler != nil {
							consumer.errorHandler(err)
						}
					}
				}
			}
		}(partitions[i])
	}

	return nil
}
