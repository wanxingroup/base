package service

import (
	"fmt"
)

const (
	defaultRPCListenHost = "127.0.0.1"
	defaultRPCListenPort = 8088
)

type RPCListenConfig struct {
	IP   string `json:"ip" yaml:"ip"`
	Port uint16 `json:"port" yaml:"port"`
}

func (config RPCListenConfig) GetIP() string {

	if config.IP == "" {

		return defaultRPCListenHost
	}

	return config.IP
}

func (config RPCListenConfig) GetPort() uint16 {

	if config.Port == 0 {

		return defaultRPCListenPort
	}

	return config.Port
}

type RPCListenOption func(config *RPCListenConfig)

func RPCListenConfigIP(ip string) RPCListenOption {

	return func(config *RPCListenConfig) {
		config.IP = ip
	}
}

func RPCListenConfigPort(port uint16) RPCListenOption {

	return func(config *RPCListenConfig) {
		config.Port = port
	}
}

func NewRPCListenConfig(options ...RPCListenOption) *RPCListenConfig {

	config := &RPCListenConfig{}

	for _, option := range options {

		option(config)
	}

	return config
}

type RPCConfig struct {
	ListenConfig *RPCListenConfig
}

func (config *RPCConfig) String() string {

	return fmt.Sprintf("%#v", config)
}

type RPCConfigOption func(config *RPCConfig)

func NewRPCConfig(options ...RPCConfigOption) *RPCConfig {

	config := &RPCConfig{}

	for _, option := range options {

		option(config)
	}

	return config
}

func RPCConfigListenConfig(ListenConfig *RPCListenConfig) RPCConfigOption {

	return func(config *RPCConfig) {

		config.ListenConfig = ListenConfig
	}
}
