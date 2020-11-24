package cache

import (
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {

	tests := []struct {
		inputKey    string
		inputConfig *RedisConfig
		wantError   error
		originPool  map[string]*RedisConnection
		wantPool    map[string]*RedisConnection
	}{
		{
			inputKey: "common",
			inputConfig: &RedisConfig{
				host: "127.0.0.1",
				port: 6379,
			},
			wantError:  nil,
			originPool: make(map[string]*RedisConnection),
			wantPool: map[string]*RedisConnection{
				"common": {
					Client: redis.NewClient(&redis.Options{
						Addr: "127.0.0.1:6379",
					}),
					RedisConfig: &RedisConfig{
						host: "127.0.0.1",
						port: 6379,
					},
				},
			},
		},
		{
			inputKey: "common",
			inputConfig: &RedisConfig{
				host: "127.0.0.1",
				port: 6379,
			},
			wantError: nil,
			originPool: map[string]*RedisConnection{
				"common": {
					Client: redis.NewClient(&redis.Options{
						Addr: "127.0.0.1:6379",
					}),
					RedisConfig: &RedisConfig{
						host: "127.0.0.1",
						port: 6379,
					},
				},
			},
			wantPool: map[string]*RedisConnection{
				"common": {
					Client: redis.NewClient(&redis.Options{
						Addr: "127.0.0.1:6379",
					}),
					RedisConfig: &RedisConfig{
						host: "127.0.0.1",
						port: 6379,
					},
				},
			},
		},
	}

	for _, test := range tests {

		pool = test.originPool
		assert.Equal(t, test.wantError, Connect(test.inputKey, test.inputConfig))
		for key, conn := range test.wantPool {

			assert.EqualValues(t, conn.RedisConfig, pool[key].RedisConfig)
		}
		for key, conn := range pool {

			assert.EqualValues(t, test.wantPool[key].RedisConfig, conn.RedisConfig)
		}
	}
}

func TestDisconnect(t *testing.T) {

	tests := []struct {
		inputKey   string
		originPool map[string]*RedisConnection
		wantPool   map[string]*RedisConnection
		wantError  error
	}{
		{
			inputKey:   "common",
			originPool: make(map[string]*RedisConnection),
			wantPool:   map[string]*RedisConnection{},
			wantError:  nil,
		},
		{
			inputKey: "common",
			originPool: map[string]*RedisConnection{
				"common": {
					Client: redis.NewClient(&redis.Options{
						Addr: "127.0.0.1:6379",
					}),
					RedisConfig: &RedisConfig{
						host: "127.0.0.1",
						port: 6379,
					},
				},
			},
			wantPool: map[string]*RedisConnection{
				"common": {
					Client: redis.NewClient(&redis.Options{
						Addr: "127.0.0.1:6379",
					}),
					RedisConfig: &RedisConfig{
						host: "127.0.0.1",
						port: 6379,
					},
				},
			},
			wantError: nil,
		},
	}

	for _, test := range tests {

		pool = test.originPool
		assert.Equal(t, test.wantError, Disconnect(test.inputKey))
		for key, conn := range test.wantPool {

			assert.EqualValues(t, conn.RedisConfig, pool[key].RedisConfig)
		}
		for key, conn := range pool {

			assert.EqualValues(t, test.wantPool[key].RedisConfig, conn.RedisConfig)
		}
	}
}

func TestGet(t *testing.T) {

	tests := []struct {
		inputKey  string
		pool      map[string]*RedisConnection
		want      *RedisConnection
		wantError error
	}{
		{
			inputKey:  "common",
			pool:      map[string]*RedisConnection{},
			wantError: errorNotHaveConnection,
		},
		{
			inputKey: "common",
			pool: map[string]*RedisConnection{
				"common": {
					Client: redis.NewClient(&redis.Options{
						Addr: "127.0.0.1:6379",
					}),
					RedisConfig: &RedisConfig{
						host: "127.0.0.1",
						port: 6379,
					},
				},
			},
			want: &RedisConnection{
				Client: redis.NewClient(&redis.Options{
					Addr: "127.0.0.1:6379",
				}),
				RedisConfig: &RedisConfig{
					host: "127.0.0.1",
					port: 6379,
				},
			},
			wantError: nil,
		},
	}

	for _, test := range tests {

		pool = test.pool
		actualConnection, actualError := Get(test.inputKey)
		if actualConnection == nil {
			assert.Equal(t, test.want, actualConnection)
		} else {
			assert.Equal(t, test.want.RedisConfig, actualConnection.RedisConfig)
		}

		assert.Equal(t, test.wantError, actualError)
	}
}
