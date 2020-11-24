package config

const defaultMySQLHost = "127.0.0.1"
const defaultMySQLPort = 3306
const defaultMySQLUsername = "root"

type MySQLConfig struct {
	Host     string `json:"host" yaml:"host"`
	Port     uint16 `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Database string `json:"database" yaml:"database"`
	LogMode  bool   `json:"logMode" yaml:"logMode"`
}

func (config MySQLConfig) GetHost() string {

	if config.Host == "" {
		return defaultMySQLHost
	}
	return config.Host
}

func (config MySQLConfig) GetPort() uint16 {

	if config.Port == 0 {
		return defaultMySQLPort
	}
	return config.Port
}

func (config MySQLConfig) GetUsername() string {

	if config.Username == "" {
		return defaultMySQLUsername
	}
	return config.Username
}

func (config MySQLConfig) GetPassword() string {

	return config.Password
}

func (config MySQLConfig) GetDatabase() string {

	return config.Database
}

func (config MySQLConfig) GetLogMode() bool {

	return config.LogMode
}
