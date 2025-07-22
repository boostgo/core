package mongox

import (
	"context"
	"crypto/tls"
	"strconv"

	"github.com/boostgo/core/timex"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
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
	URI        string `json:"uri" yaml:"uri"`
	Database   string `json:"database" yaml:"database"`
	Username   string `json:"username" yaml:"username"`
	Password   string `json:"password" yaml:"password"`
	TLS        bool   `json:"tls" yaml:"tls"`
	AuthSource string `json:"auth_source" yaml:"authSource"`

	// connection pool
	MaxPoolSize     uint64         `json:"max_pool_size" yaml:"maxPoolSize"`
	MinPoolSize     uint64         `json:"min_pool_size" yaml:"minPoolSize"`
	MaxConnIdleTime timex.Duration `json:"max_conn_idle_time" yaml:"maxConnIdleTime"`
	MaxConnecting   uint64         `json:"max_connecting" yaml:"maxConnecting"`

	// timeouts
	ConnectTimeout   timex.Duration `json:"connect_timeout" yaml:"connectTimeout"`
	ServerTimeout    timex.Duration `json:"server_timeout" yaml:"serverTimeout"`
	SocketTimeout    timex.Duration `json:"socket_timeout" yaml:"socketTimeout"`
	HeartbeatTimeout timex.Duration `json:"heartbeat_timeout" yaml:"heartbeatTimeout"`

	// retry
	RetryWrites  bool           `json:"retry_writes" yaml:"retryWrites"`
	RetryReads   bool           `json:"retry_reads" yaml:"retryReads"`
	MaxRetryTime timex.Duration `json:"max_retry_time" yaml:"maxRetryTime"`

	// read preference
	ReadPreference string `json:"read_preference" yaml:"readPreference"`
	ReadConcern    string `json:"read_concern" yaml:"readConcern"`
	WriteConcern   string `json:"write_concern" yaml:"writeConcern"`

	// compression
	Compressors []string `json:"compressors" yaml:"compressors"`
	ZlibLevel   int      `json:"zlib_level" yaml:"zlibLevel"`
}

type singleClient struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewClient(cfg Config, opts ...Option) (Client, error) {
	clientOptions := options.
		Client().
		ApplyURI(cfg.URI).
		// connection pool
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize).
		SetMaxConnIdleTime(cfg.MaxConnIdleTime.Duration()).
		SetMaxPoolSize(cfg.MaxPoolSize).
		// timeouts
		SetConnectTimeout(cfg.ConnectTimeout.Duration()).
		SetServerSelectionTimeout(cfg.ServerTimeout.Duration()).
		SetSocketTimeout(cfg.SocketTimeout.Duration()).
		SetHeartbeatInterval(cfg.HeartbeatTimeout.Duration()).
		// retry
		SetRetryWrites(cfg.RetryWrites).
		SetRetryReads(cfg.RetryReads)

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

	if cfg.ReadPreference != "" {
		readPrefMode, err := readpref.ModeFromString(cfg.ReadPreference)
		if err != nil {
			return nil, ErrReadPrefInvalidMode.SetError(err)
		}

		readPref, err := readpref.New(readPrefMode)
		if err != nil {
			return nil, ErrReadPrefCreate.SetError(err).AddParam("mode", readPrefMode.String())
		}

		clientOptions.SetReadPreference(readPref)
	}

	if cfg.ReadConcern != "" {
		var readConcern *readconcern.ReadConcern
		switch cfg.ReadConcern {
		case "local":
			readConcern = readconcern.Local()
		case "available":
			readConcern = readconcern.Available()
		case "majority":
			readConcern = readconcern.Majority()
		case "linearizable":
			readConcern = readconcern.Linearizable()
		case "snapshot":
			readConcern = readconcern.Snapshot()
		default:
			return nil, newUnsupportedConcernError(cfg.ReadConcern, "read")
		}

		clientOptions.SetReadConcern(readConcern)
	}

	if cfg.WriteConcern != "" {
		var writeConcern *writeconcern.WriteConcern
		switch cfg.WriteConcern {
		case "majority":
			writeConcern = writeconcern.New(writeconcern.WMajority())
		case "acknowledged":
			writeConcern = writeconcern.New(writeconcern.W(1))
		case "unacknowledged":
			writeConcern = writeconcern.New(writeconcern.W(0))
		default:
			if w, err := strconv.Atoi(cfg.WriteConcern); err == nil {
				writeConcern = writeconcern.New(writeconcern.W(w))
			} else {
				return nil, newUnsupportedConcernError(cfg.WriteConcern, "write")
			}
		}

		clientOptions.SetWriteConcern(writeConcern)
	}

	if len(cfg.Compressors) > 0 {
		clientOptions.SetCompressors(cfg.Compressors)
		if cfg.ZlibLevel > 0 {
			clientOptions.SetZlibLevel(cfg.ZlibLevel)
		}
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
