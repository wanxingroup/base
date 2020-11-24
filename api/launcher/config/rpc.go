package config

import (
	"time"
)

type RPCConfig struct {
	Enable           bool          `json:"enable,omitempty" yaml:"enable,omitempty"`
	IP               string        `json:"ip" yaml:"ip"`
	Port             uint16        `json:"port" yaml:"port"`
	ReadWriteTimeout time.Duration `json:"readWriteTimeout" yaml:"readWriteTimeout"`
	TTL              time.Duration `json:"ttl" yaml:"ttl"`
	Interval         time.Duration `json:"interval" yaml:"interval"`
}
