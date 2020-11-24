package config

const defaultRedisHost = "127.0.0.1"
const defaultRedisPort = 6379

type RedisConfig struct {
	Host     string `json:"host" yaml:"host"`
	Port     uint16 `json:"port" yaml:"port"`
	Password string `json:"password" yaml:"password"`
	Database uint8  `json:"database" yaml:"database"`
}

func (config RedisConfig) GetHost() string {

	if config.Host == "" {
		return defaultRedisHost
	}
	return config.Host
}

func (config RedisConfig) GetPort() uint16 {

	if config.Port == 0 {
		return defaultRedisPort
	}
	return config.Port
}

func (config RedisConfig) GetPassword() string {

	return config.Password
}

func (config RedisConfig) GetDatabase() uint8 {

	return config.Database
}
