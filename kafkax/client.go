package kafkax

import (
	"github.com/boostgo/core/trace"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

const (
	TraceProtocol = "kafka"
	TraceKey      = "X-Lite-Trace-ID"
)

func init() {
	trace.RegisterProtocol(TraceProtocol, TraceKey)
}

type Client struct {
	client sarama.Client
}

// NewClient creates new kafka client with default options as [Option]
func NewClient(cfg Config, opts ...Option) (sarama.Client, error) {
	allEmpty := true
	for _, b := range cfg.Brokers {
		if b != "" {
			allEmpty = false
			break
		}
	}

	if len(cfg.Brokers) == 0 || allEmpty {
		return nil, ErrBrokerListEmpty
	}

	config := sarama.NewConfig()

	apply := make([]Option, 0, len(opts)+1)
	apply = append(apply, clientOption(cfg))
	apply = append(apply, opts...)

	for _, opt := range apply {
		opt(config)
	}

	client, err := sarama.NewClient(cfg.Brokers, config)
	if err != nil {
		return nil, NewConnectError(err, cfg, config)
	}

	return &Client{
		client: client,
	}, nil
}

// MustClient calls [NewClient] and if catches error will be thrown panic
func MustClient(cfg Config, opts ...Option) sarama.Client {
	client, err := NewClient(cfg, opts...)
	if err != nil {
		panic(err)
	}

	return client
}

// NewCluster creates new kafka cluster with default options as [Option] by client
func NewCluster(client sarama.Client) (sarama.ClusterAdmin, error) {
	clusterClient, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		return nil, err
	}

	return clusterClient, nil
}

// MustCluster calls [NewCluster] and if errors is catch throws panic
func MustCluster(client sarama.Client) sarama.ClusterAdmin {
	cluster, err := NewCluster(client)
	if err != nil {
		panic(err)
	}

	return cluster
}

func (c *Client) Config() *sarama.Config {
	return c.client.Config()
}

func (c *Client) Controller() (*sarama.Broker, error) {
	return c.client.Controller()
}

func (c *Client) RefreshController() (*sarama.Broker, error) {
	return c.client.RefreshController()
}

func (c *Client) Brokers() []*sarama.Broker {
	return c.client.Brokers()
}

func (c *Client) Broker(brokerID int32) (*sarama.Broker, error) {
	return c.client.Broker(brokerID)
}

func (c *Client) Topics() ([]string, error) {
	return c.client.Topics()
}

func (c *Client) Partitions(topic string) ([]int32, error) {
	return c.client.Partitions(topic)
}

func (c *Client) WritablePartitions(topic string) ([]int32, error) {
	return c.client.WritablePartitions(topic)
}

func (c *Client) Leader(topic string, partitionID int32) (*sarama.Broker, error) {
	return c.client.Leader(topic, partitionID)
}

func (c *Client) LeaderAndEpoch(topic string, partitionID int32) (*sarama.Broker, int32, error) {
	return c.client.LeaderAndEpoch(topic, partitionID)
}

func (c *Client) Replicas(topic string, partitionID int32) ([]int32, error) {
	return c.client.Replicas(topic, partitionID)
}

func (c *Client) InSyncReplicas(topic string, partitionID int32) ([]int32, error) {
	return c.client.InSyncReplicas(topic, partitionID)
}

func (c *Client) OfflineReplicas(topic string, partitionID int32) ([]int32, error) {
	return c.client.OfflineReplicas(topic, partitionID)
}

func (c *Client) RefreshBrokers(brokers []string) error {
	return c.client.RefreshBrokers(brokers)
}

func (c *Client) RefreshMetadata(topics ...string) error {
	return c.client.RefreshMetadata(topics...)
}

func (c *Client) GetOffset(topic string, partitionID int32, time int64) (int64, error) {
	return c.client.GetOffset(topic, partitionID, time)
}

func (c *Client) Coordinator(consumerGroup string) (*sarama.Broker, error) {
	return c.client.Coordinator(consumerGroup)
}

func (c *Client) RefreshCoordinator(consumerGroup string) error {
	return c.client.RefreshCoordinator(consumerGroup)
}

func (c *Client) TransactionCoordinator(transactionID string) (*sarama.Broker, error) {
	return c.client.TransactionCoordinator(transactionID)
}

func (c *Client) RefreshTransactionCoordinator(transactionID string) error {
	return c.client.RefreshTransactionCoordinator(transactionID)
}

func (c *Client) InitProducerID() (*sarama.InitProducerIDResponse, error) {
	return c.client.InitProducerID()
}

func (c *Client) LeastLoadedBroker() *sarama.Broker {
	return c.client.LeastLoadedBroker()
}

func (c *Client) PartitionNotReadable(topic string, partition int32) bool {
	return c.client.PartitionNotReadable(topic, partition)
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) Closed() bool {
	return c.client.Closed()
}

// clientOption returns default options for client as [Option]
func clientOption(cfg Config) Option {
	return func(config *sarama.Config) {
		config.ClientID = buildClientID()

		if cfg.Username != "" && cfg.Password != "" {
			config.Net.SASL.Enable = true
			config.Net.SASL.Handshake = true
			config.Net.SASL.Mechanism = "PLAIN"
			config.Net.SASL.User = cfg.Username
			config.Net.SASL.Password = cfg.Password
		}
	}
}

var clientIdPrefix = ""

const defaultClientIdPrefix = "lite-app-"

func SetClientIdPrefix(prefix string) {
	clientIdPrefix = prefix
}

func buildClientID() string {
	prefix := clientIdPrefix
	if prefix == "" {
		prefix = defaultClientIdPrefix
	}

	return prefix + uuid.New().String()
}

func joinOptions(options []Option, joined ...Option) []Option {
	return append(options, joined...)
}
