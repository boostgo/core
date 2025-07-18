package samples

import (
	"fmt"
	"time"

	"github.com/boostgo/core/echox"
	"github.com/boostgo/core/log"
	"github.com/boostgo/core/log/logx"
	"github.com/boostgo/core/redis"
	"github.com/boostgo/core/sql"
	"github.com/boostgo/core/timex"
	"github.com/boostgo/core/translate"
	translateEcho "github.com/boostgo/core/translate/echo"

	"github.com/jmoiron/sqlx"
	"github.com/swaggo/swag"
)

type Settings struct {
	PrettyLog bool `json:"pretty_log" yaml:"prettyLog"`
}

func (s Settings) Init() {
	if s.PrettyLog {
		logx.Pretty()
	}
}

type Server struct {
	Host         string         `json:"host" yaml:"host"`
	Port         int            `json:"port" yaml:"port"`
	ShutdownWait timex.Duration `json:"shutdown_wait" yaml:"shutdownWait"`
}

func (s Server) Address() string {
	host := s.Host
	port := s.Port

	if host == "" {
		host = "0.0.0.0"
	}

	if port == 0 {
		port = 80
	}

	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (s Server) Log() {
	log.
		Info().
		Str("host", s.Host).
		Int("port", s.Port).
		Object("shutdown_wait", s.ShutdownWait).
		Msg("Server config")
}

type GrpcServer struct {
	Port int `json:"port" yaml:"port"`
}

func (s GrpcServer) Log() {
	log.
		Info().
		Int("port", s.Port).
		Msg("GRPC server config")
}

type Swagger struct {
	Host             string   `json:"host" yaml:"host"`
	Version          string   `json:"version" yaml:"version"`
	BasePath         string   `json:"base_path" yaml:"basePath"`
	Schemes          []string `json:"schemes" yaml:"schemes"`
	Title            string   `json:"title" yaml:"title"`
	Description      string   `json:"description" yaml:"description"`
	InfoInstanceName string   `json:"info_instance_name" yaml:"infoInstanceName"`
	SwaggerTemplate  string   `json:"swagger_template" yaml:"swaggerTemplate"`
	LeftDelim        string   `json:"left_delim" yaml:"leftDelim"`
	RightDelim       string   `json:"right_delim" yaml:"rightDelim"`
}

func (s Swagger) Init(cfg *swag.Spec) {
	if s.Host != "" {
		cfg.Host = s.Host
	}

	if s.Version != "" {
		cfg.Version = s.Version
	}

	if s.BasePath != "" {
		cfg.BasePath = s.BasePath
	}

	if s.Schemes != nil {
		cfg.Schemes = s.Schemes
	}

	if s.Title != "" {
		cfg.Title = s.Title
	}

	if s.Description != "" {
		cfg.Description = s.Description
	}

	if s.InfoInstanceName != "" {
		cfg.InfoInstanceName = s.InfoInstanceName
	}

	if s.SwaggerTemplate != "" {
		cfg.SwaggerTemplate = s.SwaggerTemplate
	}

	if s.LeftDelim != "" {
		cfg.LeftDelim = s.LeftDelim
	}

	if s.RightDelim != "" {
		cfg.RightDelim = s.RightDelim
	}
}

func (s Swagger) Log() {
	log.
		Info().
		Str("host", s.Host).
		Str("version", s.Version).
		Str("base_path", s.BasePath).
		Strs("schemes", s.Schemes).
		Str("title", s.Title).
		Str("description", s.Description).
		Str("info_instance_name", s.InfoInstanceName).
		Str("swagger_template", s.SwaggerTemplate).
		Str("left_delim", s.LeftDelim).
		Str("right_delim", s.RightDelim).
		Msg("Swagger config")
}

type Auth struct {
	Secret    string `json:"secret" yaml:"secret" env:"AUTH_SECRET"`
	PublicKey string `json:"public_key" yaml:"publicKey" env:"PUBLIC_SECRET_KEY"`
	Algorithm string `json:"algorithm" yaml:"algorithm"`
}

func (a Auth) Log() {
	log.
		Info().
		Str("algorithm", a.Algorithm).
		Str("secret", a.Secret).
		Str("public_key", a.PublicKey).
		Msg("Auth config")
}

type Translate struct {
	Path      string             `json:"path" yaml:"path"`
	Extension string             `json:"extension" yaml:"extension"`
	Header    string             `json:"header" yaml:"header"`
	Locales   []translate.Locale `json:"locales" yaml:"locales"`
}

func (t Translate) NewTranslator() *translate.Translator {
	path := t.Path
	if path == "" {
		path = "./translations"
	}

	ext := t.Extension
	if ext == "" {
		ext = translate.ExtYaml
	}

	return translate.NewTranslator(path, ext, t.Locales...)
}

func (t Translate) RegisterFailureMiddleware() {
	header := t.Header
	if header == "" {
		header = "Content-Language"
	}

	translator := t.NewTranslator()
	if err := translator.Read(); err != nil {
		panic(err)
	}

	echox.RegisterFailureMiddleware(translateEcho.FailureMiddleware(translator, header))
}

func (t Translate) Log() {
	locales := make([]string, len(t.Locales))
	for i := range t.Locales {
		locales[i] = t.Locales[i].String()
	}

	log.
		Info().
		Str("path", t.Path).
		Str("extension", t.Extension).
		Str("header", t.Header).
		Strs("locales", locales).
		Msg("Translate config")
}

type SQL struct {
	Host               string `json:"host" yaml:"host"`
	Port               int    `json:"port" yaml:"port"`
	Username           string `json:"username" yaml:"username"`
	Password           string `json:"password" yaml:"password"`
	Database           string `json:"database" yaml:"database"`
	BinaryParameters   bool   `json:"binary_parameters" yaml:"binaryParameters"`
	MaxOpenConnections int    `json:"max_open_connections" yaml:"maxOpenConnections"`
	MaxIdleConnections int    `json:"max_idle_connections" yaml:"maxIdleConnections"`
	MaxLifetime        int    `json:"max_lifetime" yaml:"maxLifetime"`
	MaxIdleTime        int    `json:"max_idle_time" yaml:"maxIdleTime"`
	ReadTimeout        int    `json:"read_timeout" yaml:"readTimeout"`
	WriteTimeout       int    `json:"write_timeout" yaml:"writeTimeout"`
	Driver             string `json:"driver" yaml:"driver"`
}

func (s SQL) ConnectionString() string {
	connector := sql.
		NewConnector().
		Host(s.Host).
		Port(s.Port).
		Username(s.Username).
		Password(s.Password).
		Database(s.Database).
		BinaryParameters(s.BinaryParameters).
		MaxOpenConnections(s.MaxOpenConnections).
		MaxIdleConnections(s.MaxIdleConnections).
		ConnectionMaxLifetime(time.Duration(s.MaxLifetime) * time.Second).
		MaxIdleTime(time.Duration(s.MaxIdleTime) * time.Second).
		ReadTimeout(s.ReadTimeout).
		WriteTimeout(s.WriteTimeout)

	if s.Driver == sql.ChDriver {
		return connector.BuildClickhouse()
	}

	return connector.Build()
}

func (s SQL) Connect(timeout time.Duration) (*sqlx.DB, error) {
	return sql.Connect(s.Driver, s.ConnectionString(), timeout)
}

func (s SQL) MustConnect(timeout time.Duration) *sqlx.DB {
	conn, err := s.Connect(timeout)
	if err != nil {
		panic(err)
	}

	return conn
}

func (s SQL) Log() {
	log.
		Info().
		Str("host", s.Host).
		Int("port", s.Port).
		Str("database", s.Database).
		Str("username", s.Username).
		Bool("binary_parameters", s.BinaryParameters).
		Int("max_open_connections", s.MaxOpenConnections).
		Int("max_idle_connections", s.MaxIdleConnections).
		Int("max_lifetime", s.MaxLifetime).
		Int("max_idle_time", s.MaxIdleTime).
		Str("driver", s.Driver).
		Msg("SQL config")
}

type RedisSingle struct {
	Address  string `json:"address" yaml:"address"`
	Port     int    `json:"port" yaml:"port"`
	Password string `json:"password" yaml:"password"`
	DB       int    `json:"db" yaml:"db"`
}

func (r RedisSingle) Connect(password ...string) (redis.Client, error) {
	pwd := r.Password
	if len(password) > 0 {
		pwd = password[0]
	}

	return redis.New(r.Address, r.Port, r.DB, pwd)
}

func (r RedisSingle) MustConnect(password ...string) redis.Client {
	client, err := r.Connect(password...)
	if err != nil {
		panic(err)
	}

	return client
}

func (r RedisSingle) Log() {
	log.
		Info().
		Str("host", r.Address).
		Int("port", r.Port).
		Int("db", r.DB).
		Msg("Redis single config")
}

type RedisShards []redis.ShardConnectConfig

func (r RedisShards) Connect(selector redis.ClientSelector, password ...string) (redis.Client, error) {
	if len(password) > 0 {
		var pwd string
		for i := 0; i < len(r); i++ {
			r[i].Password = pwd
		}
	}

	clients, err := redis.ConnectShards(r, selector)
	if err != nil {
		return nil, err
	}

	return redis.NewShard(clients), nil
}

func (r RedisShards) MustConnect(selector redis.ClientSelector, password ...string) redis.Client {
	client, err := r.Connect(selector, password...)
	if err != nil {
		panic(err)
	}

	return client
}

func (r RedisShards) Log() {
	for idx, shard := range r {
		log.
			Info().
			Str("key", shard.Key).
			Str("address", shard.Address).
			Int("port", shard.Port).
			Int("db", shard.DB).
			Msgf("Redis shard #%d config", idx+1)
	}
}

type Redis struct {
	Single  RedisSingle `json:"single" yaml:"single"`
	Cluster RedisSingle `json:"cluster" yaml:"cluster"`
	Shards  RedisShards `json:"shards" yaml:"shards"`
}

func (r Redis) Log() {
	r.Single.Log()
	r.Cluster.Log()
	r.Shards.Log()
}

type Keycloak struct {
	Host string `json:"host" yaml:"host"`
}

func (k Keycloak) Log() {
	log.
		Info().
		Str("host", k.Host).
		Msg("Keycloak config")
}
