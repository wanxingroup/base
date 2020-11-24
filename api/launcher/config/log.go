package config

import (
	"github.com/sirupsen/logrus"
)

type LogConfig struct {
	Level string `json:"level" yaml:"level"`
}

func (config LogConfig) GetLogLevel() logrus.Level {

	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		return logrus.InfoLevel
	}

	return level
}
