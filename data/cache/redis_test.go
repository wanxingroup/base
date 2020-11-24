package cache

import (
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestRedisConfig_GetConfigKey(t *testing.T) {

	tests := []struct {
		input RedisConfig
		want  string
	}{
		{
			input: RedisConfig{
				host:     "127.0.0.1",
				port:     6379,
				password: "",
				database: 0,
			},
			want: "127.0.0.1:6379",
		},
		{
			input: RedisConfig{
				host:     "redis-01.internal.com",
				port:     16379,
				password: "",
				database: 0,
			},
			want: "redis-01.internal.com:16379",
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.want, test.input.GetConfigKey())
	}
}

func TestRedisConfig_GetClientOptions(t *testing.T) {

	tests := []struct {
		input RedisConfig
		want  *redis.Options
	}{
		{
			input: RedisConfig{
				host:     "127.0.0.1",
				port:     6379,
				password: "",
				database: 0,
			},
			want: &redis.Options{
				Addr:     "127.0.0.1:6379",
				Password: "",
				DB:       0,
			},
		},
		{
			input: RedisConfig{
				host:     "redis-01.internal.com",
				port:     16379,
				password: "",
				database: 0,
			},
			want: &redis.Options{
				Addr:     "redis-01.internal.com:16379",
				Password: "",
				DB:       0,
			},
		},
		{
			input: RedisConfig{
				host:     "redis-01.internal.com",
				port:     16379,
				password: "123456",
				database: 0,
			},
			want: &redis.Options{
				Addr:     "redis-01.internal.com:16379",
				Password: "123456",
				DB:       0,
			},
		},
		{
			input: RedisConfig{
				host:     "redis-01.internal.com",
				port:     16379,
				password: "123456",
				database: 2,
			},
			want: &redis.Options{
				Addr:     "redis-01.internal.com:16379",
				Password: "123456",
				DB:       2,
			},
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.want, test.input.GetClientOptions())
	}
}

func TestRedisHost(t *testing.T) {

	tests := []struct {
		inputHost   string
		inputConfig RedisConfig
		want        RedisConfig
	}{
		{
			inputHost:   "127.0.0.1",
			inputConfig: RedisConfig{},
			want: RedisConfig{
				host: "127.0.0.1",
			},
		},
		{
			inputHost:   "redis.host",
			inputConfig: RedisConfig{},
			want: RedisConfig{
				host: "redis.host",
			},
		},
	}

	for _, test := range tests {

		RedisHost(test.inputHost)(&test.inputConfig)
		assert.Equal(t, test.want, test.inputConfig)
	}
}

func TestRedisPort(t *testing.T) {

	tests := []struct {
		inputPort   uint16
		inputConfig RedisConfig
		want        RedisConfig
	}{
		{
			inputPort:   6379,
			inputConfig: RedisConfig{},
			want: RedisConfig{
				port: 6379,
			},
		},
		{
			inputPort:   16379,
			inputConfig: RedisConfig{},
			want: RedisConfig{
				port: 16379,
			},
		},
	}

	for _, test := range tests {

		RedisPort(test.inputPort)(&test.inputConfig)
		assert.Equal(t, test.want, test.inputConfig)
	}
}

func TestRedisPassword(t *testing.T) {

	tests := []struct {
		input       string
		inputConfig RedisConfig
		want        RedisConfig
	}{
		{
			input:       "6379",
			inputConfig: RedisConfig{},
			want: RedisConfig{
				password: "6379",
			},
		},
		{
			input:       "16379",
			inputConfig: RedisConfig{},
			want: RedisConfig{
				password: "16379",
			},
		},
	}

	for _, test := range tests {

		RedisPassword(test.input)(&test.inputConfig)
		assert.Equal(t, test.want, test.inputConfig)
	}
}

func TestRedisDatabase(t *testing.T) {

	tests := []struct {
		input       uint8
		inputConfig RedisConfig
		want        RedisConfig
	}{
		{
			input:       0,
			inputConfig: RedisConfig{},
			want: RedisConfig{
				database: 0,
			},
		},
		{
			input:       1,
			inputConfig: RedisConfig{},
			want: RedisConfig{
				database: 1,
			},
		},
	}

	for _, test := range tests {

		RedisDatabase(test.input)(&test.inputConfig)
		assert.Equal(t, test.want, test.inputConfig)
	}
}

func TestNewRedisConfig(t *testing.T) {

	tests := []struct {
		input []RedisOption
		want  *RedisConfig
	}{
		{
			input: []RedisOption{
				RedisHost("127.0.0.1"),
				RedisPort(6379),
			},
			want: &RedisConfig{
				host: "127.0.0.1",
				port: 6379,
			},
		},
		{
			input: []RedisOption{
				RedisHost("127.0.0.1"),
				RedisPort(6379),
			},
			want: &RedisConfig{
				host: "127.0.0.1",
				port: 6379,
			},
		},
		{
			input: []RedisOption{
				RedisHost("127.0.0.1"),
				RedisPort(16379),
				RedisPassword("password"),
			},
			want: &RedisConfig{
				host:     "127.0.0.1",
				port:     16379,
				password: "password",
			},
		},
		{
			input: []RedisOption{
				RedisHost("127.0.0.1"),
				RedisPort(26379),
				RedisDatabase(2),
			},
			want: &RedisConfig{
				host:     "127.0.0.1",
				port:     26379,
				database: 2,
			},
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.want, NewRedisConfig(test.input...))
	}
}

func TestRedisConnection_TryConnect(t *testing.T) {

	tests := []struct {
		input *RedisConnection
		want  error
	}{
		{
			input: &RedisConnection{
				Client: redis.NewClient(&redis.Options{
					Addr: "127.0.0.1:6379",
				}),
				RedisConfig: &RedisConfig{
					host: "127.0.0.1",
					port: 6379,
				},
			},
			want: nil,
		},
		{
			input: &RedisConnection{
				Client: redis.NewClient(&redis.Options{
					Addr: "127.0.0.1:16379",
				}),
				RedisConfig: &RedisConfig{
					host: "127.0.0.1",
					port: 16379,
				},
			},
			want: errorPingFailed,
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.want, test.input.TryConnect())
	}
}

func TestRedisConnection_Connect(t *testing.T) {

	tests := []struct {
		input *RedisConnection
		want  error
	}{
		{
			input: &RedisConnection{
				Client: redis.NewClient(&redis.Options{
					Addr: "127.0.0.1:6379",
				}),
				RedisConfig: &RedisConfig{
					host: "127.0.0.1",
					port: 6379,
				},
			},
			want: nil,
		},
		{
			input: &RedisConnection{},
			want:  errorNotHaveConfig,
		},
		{
			input: &RedisConnection{
				RedisConfig: &RedisConfig{
					host: "127.0.0.1",
					port: 6379,
				},
			},
			want: nil,
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.want, test.input.Connect())
	}
}
