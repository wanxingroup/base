package database

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

var (
	pool    = make(map[string]*Connection)
	configs = make(map[string]*MySQLConfig)
)

func Connect(key string, config *MySQLConfig) error {

	if conn, ok := pool[key]; ok {
		return conn.Connect()
	}

	conn := &Connection{
		config: config,
	}

	// 将连接放入缓存
	pool[key] = conn

	return conn.Connect()
}

func Disconnect(key string) error {

	conn, ok := pool[key]
	if !ok {
		return errors.New(fmt.Sprintf("%s database connection not exist", key))
	}

	return conn.Close()
}

func GetDB(key string) *gorm.DB {

	if conn, ok := pool[key]; ok {
		return conn.DB
	}

	return nil
}
