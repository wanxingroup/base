package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisOption func(config *RedisConfig)

type RedisConfig struct {
	host     string
	port     uint16
	password string
	database uint8
}

func (config *RedisConfig) GetConfigKey() string {
	return fmt.Sprintf("%s:%d", config.host, config.port)
}

func (config *RedisConfig) GetClientOptions() *redis.Options {
	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.host, config.port),
		Password: config.password,      // no password set
		DB:       int(config.database), // use default DB
	}
}

var defaultRedisConfig = RedisConfig{
	host:     "127.0.0.1",
	port:     6379,
	password: "",
	database: 0,
}

var configs = make(map[string]*RedisConfig)

var defaultConfigKey = "127.0.0.1:6379"

func RedisHost(host string) RedisOption {
	return func(config *RedisConfig) {
		config.host = host
	}
}

func RedisPort(port uint16) RedisOption {
	return func(config *RedisConfig) {
		config.port = port
	}
}

func RedisPassword(password string) RedisOption {
	return func(config *RedisConfig) {
		config.password = password
	}
}

func RedisDatabase(database uint8) RedisOption {
	return func(config *RedisConfig) {
		config.database = database
	}
}

func NewRedisConfig(options ...RedisOption) *RedisConfig {

	config := defaultRedisConfig
	for _, option := range options {
		option(&config)
	}
	logger.Infof("redis config: %#v", config)

	configKey := config.GetConfigKey()
	if conf, ok := configs[configKey]; ok {
		return conf
	}

	if len(configs) == 0 {
		defaultConfigKey = configKey
	}

	configs[defaultConfigKey] = &config

	return &config
}

type RedisConnection struct {
	*redis.Client
	*RedisConfig
}

func (conn *RedisConnection) TryConnect() error {

	ctx := context.Background()
	result, err := conn.Ping(ctx).Result()
	if err != nil {

		logger.WithError(err).
			WithContext(ctx).
			WithField("config", conn.RedisConfig).
			Error("ping failed")
		return errorPingFailed
	}

	if result != "PONG" {
		logger.WithField("result", result).
			WithField("config", conn.RedisConfig).
			Error("result not expect")
		return errorPingResultNotExpect
	}

	return nil
}

func (conn *RedisConnection) Connect() error {

	if conn.Client != nil {
		return nil
	}

	if conn.RedisConfig == nil {
		return errorNotHaveConfig
	}

	conn.Client = redis.NewClient(conn.RedisConfig.GetClientOptions())
	return conn.TryConnect()
}
