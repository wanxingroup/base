package service

import (
	"fmt"
	"time"
)

const (
	defaultGinListenHost    = "127.0.0.1"
	defaultGinListenPort    = 8080
	defaultReadWriteTimeout = time.Minute
)

type GinListenOption func(config *GinListenConfig)

type GinListenConfig struct {
	IP               string
	Port             uint16
	ReadWriteTimeout time.Duration
}

func GinListenConfigIP(ip string) GinListenOption {

	return func(config *GinListenConfig) {
		config.IP = ip
	}
}

func GinListenConfigPort(port uint16) GinListenOption {

	return func(config *GinListenConfig) {
		config.Port = port
	}
}

func GinListenConfigReadWriteTimeout(duration time.Duration) GinListenOption {

	return func(config *GinListenConfig) {
		config.ReadWriteTimeout = duration
	}
}

func NewGinListenConfig(options ...GinListenOption) *GinListenConfig {

	config := &GinListenConfig{}

	for _, option := range options {

		option(config)
	}

	return config
}

func (config GinListenConfig) GetIP() string {

	if config.IP == "" {

		return defaultGinListenHost
	}

	return config.IP
}

func (config GinListenConfig) GetPort() uint16 {

	if config.Port == 0 {

		return defaultGinListenPort
	}

	return config.Port
}

func (config GinListenConfig) GetReadWriteTimeout() time.Duration {

	if config.ReadWriteTimeout == 0 {

		return defaultReadWriteTimeout
	}

	return config.ReadWriteTimeout
}

// Gin Service Mode
type WebServiceMode string

const (
	WebServiceModeDebug   WebServiceMode = "debug"
	WebServiceModeRelease WebServiceMode = "release"
)

type GinConfigOption func(config *GinConfig)

type GinConfig struct {
	ListenConfig   *GinListenConfig
	WebServiceMode WebServiceMode
}

func NewGinConfig(options ...GinConfigOption) *GinConfig {

	config := &GinConfig{}

	for _, option := range options {

		option(config)
	}

	return config
}

func (config *GinConfig) GetMode() WebServiceMode {

	return config.WebServiceMode
}

func (config *GinConfig) String() string {

	return fmt.Sprintf("%#v", config)
}

func GinConfigListenConfig(ListenConfig *GinListenConfig) GinConfigOption {

	return func(config *GinConfig) {

		config.ListenConfig = ListenConfig
	}
}

func GinConfigWebServiceMode(mode WebServiceMode) GinConfigOption {
	return func(config *GinConfig) {

		config.WebServiceMode = mode
	}
}
