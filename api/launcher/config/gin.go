package config

import (
	"time"
)

type GinConfig struct {
	Enable           bool           `json:"enable,omitempty" yaml:"enable,omitempty"`
	IP               string         `json:"ip" yaml:"ip"`
	Port             uint16         `json:"port" yaml:"port"`
	Mode             WebServiceMode `json:"mode" yaml:"mode"`
	ReadWriteTimeout time.Duration  `json:"readWriteTimeout" yaml:"readWriteTimeout"`
}

// Gin Service Mode
type WebServiceMode string

const (
	WebServiceModeDebug   WebServiceMode = "debug"
	WebServiceModeRelease WebServiceMode = "release"
)
