package kafkax

import (
	"github.com/IBM/sarama"
)

type Option func(*sarama.Config)

type Config struct {
	Brokers  []string
	Username string
	Password string
}

func (cfg Config) Copy() Config {
	return Config{
		Brokers:  cfg.Brokers,
		Username: cfg.Username,
		Password: cfg.Password,
	}
}

func With(fn func(*sarama.Config)) Option {
	return func(cfg *sarama.Config) {
		fn(cfg)
	}
}

func validateConsumerGroupConfig(config Config) error {
	if len(config.Brokers) == 0 {
		return ErrAtLeastOneBrokerRequired
	}

	return nil
}

func validateConsumerConfig(config Config) error {
	if len(config.Brokers) == 0 {
		return ErrAtLeastOneBrokerRequired
	}

	return nil
}
