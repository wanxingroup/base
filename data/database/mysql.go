package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	defaultMySQLConfig = MySQLConfig{
		host:     "127.0.0.1",
		port:     3306,
		user:     "root",
		password: "",
		database: "",
		logMode:  false,
	}
)

type MySQLOption func(config *MySQLConfig)

func MySQLHost(host string) MySQLOption {
	return func(config *MySQLConfig) {
		config.host = host
	}
}

func MySQLPort(port uint16) MySQLOption {
	return func(config *MySQLConfig) {
		config.port = port
	}
}

func MySQLUsername(username string) MySQLOption {
	return func(config *MySQLConfig) {
		config.user = username
	}
}

func MySQLPassword(password string) MySQLOption {
	return func(config *MySQLConfig) {
		config.password = password
	}
}

func MySQLDatabase(database string) MySQLOption {
	return func(config *MySQLConfig) {
		config.database = database
	}
}

func MySQLLogMode(enable bool) MySQLOption {
	return func(config *MySQLConfig) {
		config.logMode = enable
	}
}

type MySQLConfig struct {
	host     string
	port     uint16
	user     string
	password string
	database string
	logMode  bool
}

func NewMySQLConfig(options ...MySQLOption) *MySQLConfig {

	config := defaultMySQLConfig
	for _, option := range options {
		option(&config)
	}
	logger.Infof("mysql config: %#v", config)

	if conf, ok := configs[config.database]; ok {
		return conf
	}

	configs[config.database] = &config

	return &config
}

func (m *MySQLConfig) Connect() (err error) {

	var conn = &Connection{
		config: m,
	}
	return conn.Connect()
}

func (m MySQLConfig) getConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&loc=%s&parseTime=true",
		m.user, m.password, m.host, m.port, m.database, "Local")
}

type Connection struct {
	*gorm.DB
	config *MySQLConfig
}

func (conn *Connection) Connect() (err error) {

	if conn.isConnected() {
		return
	}

	return conn.reconnect()
}

func (conn *Connection) reconnect() (err error) {

	logger.WithField("connectionString", conn.config.getConnectionString()).
		Info("mysql connect")

	conn.DB, err = gorm.Open("mysql", conn.config.getConnectionString())
	if err != nil {
		logger.Errorf("mysql connect: %v", err)
		return
	}

	conn.DB.SingularTable(true)

	logger.Info("mysql connect succeed")

	if settingError := conn.DB.Exec("SET time_zone='+08:00'").Error; settingError != nil {
		logger.WithError(settingError).Error("set time zone error")
	} else {
		logger.Info("set time zone for +08:00")
	}

	if settingError := conn.DB.Exec("SET NAMES utf8mb4").Error; settingError != nil {
		logger.WithError(settingError).Error("set names error")
	} else {
		logger.Info("set names utf8mb4")
	}

	logger.WithField("logMode", conn.config.logMode).Info("setting gorm log mode")
	conn.DB.LogMode(conn.config.logMode)

	return
}

func (conn *Connection) isConnected() (returnValue bool) {
	if conn.DB == nil {
		return false
	}

	defer func() {

		if err := recover(); err != nil {
			returnValue = false
			logger.WithField("error", err).
				WithField("connection", conn.config.getConnectionString()).
				Errorf("panic by check isConnected")
		}
		return
	}()

	db := conn.DB.DB()
	if db == nil { // 根据 DB() 的实现，这个不可能出现
		return false
	}

	if err := db.Ping(); err != nil {
		return false
	}

	return true
}

func (conn *Connection) Begin() *gorm.DB {

	return conn.DB.Begin()
}

func (conn *Connection) Close() error {

	return conn.DB.Close()
}
