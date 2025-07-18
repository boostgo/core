package kafkax

import (
	"github.com/boostgo/core/errorx"

	"github.com/IBM/sarama"
)

var (
	ErrConnect             = errorx.New("kafkax.connect")
	ErrCreateConsumerGroup = errorx.New("kafkax.create_consumer_group")

	ErrBrokerListEmpty          = errorx.New("kafkax.config.broker_list_empty")
	ErrAtLeastOneBrokerRequired = errorx.New("kafkax.config.at_least_one_broker")
	ErrConsumerGroupIdRequired  = errorx.New("kafkax.config.group_id_required")

	ErrGroupHandlerClosed = errorx.New("kafkax.consumer_group.handler_closed")
	ErrProduceMessages    = errorx.New("kafkax.sync_producer.produce_messages")
)

type connectContext struct {
	Brokers      []string `json:"brokers"`
	Username     string   `json:"username"`
	Password     string   `json:"password"`
	SaramaConfig any      `json:"sarama_config"`
}

func NewConnectError(err error, cfg Config, saramaCfg *sarama.Config) error {
	return ErrConnect.
		SetError(err).
		SetData(connectContext{
			Brokers:      cfg.Brokers,
			Username:     cfg.Username,
			Password:     cfg.Password,
			SaramaConfig: saramaCfg,
		})
}
