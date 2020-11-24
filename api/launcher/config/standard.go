package config

import (
	"fmt"
)

const (
	DefaultServiceId = 0
)

type StandardConfig struct {
	Web       GinConfig              `json:"web" yaml:"web"`
	RPC       RPCConfig              `json:"rpc" yaml:"rpc"`
	MySQL     map[string]MySQLConfig `json:"mysql" yaml:"mysql"`
	Redis     map[string]RedisConfig `json:"redis" yaml:"redis"`
	Log       LogConfig              `json:"log" yaml:"log"`
	ServiceId uint16                 `json:"-" yaml:"-"` // used to distinguish between different services when highly available. no parse from configuration file, because services will use the same configuration file.
}

func (config StandardConfig) String() string {

	return fmt.Sprintf("%#v", config)
}

func (config *StandardConfig) GetServiceId() uint16 {
	if config.ServiceId == 0 {
		return DefaultServiceId
	}
	return config.ServiceId
}
