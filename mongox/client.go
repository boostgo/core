package mongox

import (
	"context"
	"crypto/tls"

	"github.com/boostgo/core/timex"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Option func(opts *options.ClientOptions)

type Client interface {
	Close(ctx context.Context) error
	Ping(ctx context.Context) error
	Database() *mongo.Database
	Collection(name string) *mongo.Collection
	EnsureIndexes(ctx context.Context, collection string, indexes []mongo.IndexModel) error
}

type Config struct {
	URI            string         `json:"uri" yaml:"uri"`
	Database       string         `json:"database" yaml:"database"`
	Username       string         `json:"username" yaml:"username"`
	Password       string         `json:"password" yaml:"password"`
	MaxPoolSize    uint64         `json:"max_pool_size" yaml:"maxPoolSize"`
	MinPoolSize    uint64         `json:"min_pool_size" yaml:"minPoolSize"`
	ConnectTimeout timex.Duration `json:"connect_timeout" yaml:"connectTimeout"`
	ServerTimeout  timex.Duration `json:"server_timeout" yaml:"serverTimeout"`
	TLS            bool           `json:"tls" yaml:"tls"`
	AuthSource     string         `json:"auth_source" yaml:"authSource"`
}

type singleClient struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewClient(cfg Config, opts ...Option) (Client, error) {
	clientOptions := options.
		Client().
		ApplyURI(cfg.URI).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize).
		SetConnectTimeout(cfg.ConnectTimeout.Duration()).
		SetServerSelectionTimeout(cfg.ServerTimeout.Duration())

	for _, opt := range opts {
		opt(clientOptions)
	}

	if cfg.Username != "" && cfg.Password != "" {
		credential := options.Credential{
			Username:   cfg.Username,
			Password:   cfg.Password,
			AuthSource: cfg.AuthSource,
		}
		clientOptions.SetAuth(credential)
	}

	if cfg.TLS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
		}
		clientOptions.SetTLSConfig(tlsConfig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout.Duration())
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return &singleClient{
		client:   client,
		database: client.Database(cfg.Database),
	}, nil
}

func MustClient(cfg Config, opts ...Option) Client {
	client, err := NewClient(cfg, opts...)
	if err != nil {
		panic(err)
	}

	return client
}

func (c *singleClient) Close(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

func (c *singleClient) Ping(ctx context.Context) error {
	return c.client.Ping(ctx, readpref.Primary())
}

func (c *singleClient) Database() *mongo.Database {
	return c.database
}

func (c *singleClient) Collection(name string) *mongo.Collection {
	return c.database.Collection(name)
}

func (c *singleClient) EnsureIndexes(ctx context.Context, collection string, indexes []mongo.IndexModel) error {
	coll := c.Collection(collection)

	if _, err := coll.Indexes().CreateMany(ctx, indexes); err != nil {
		return newCreateIndexError(err, collection, indexes)
	}

	return nil
}
