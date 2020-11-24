package database

import (
	"errors"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {

	tests := []struct {
		inputKey    string
		inputConfig *MySQLConfig
		err         error
	}{
		{
			inputKey:    "test-connect",
			inputConfig: NewMySQLConfig(MySQLDatabase(databaseName)),
			err:         nil,
		},
		{
			inputKey:    "test-connect",
			inputConfig: NewMySQLConfig(MySQLDatabase(databaseName)),
			err:         nil,
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.err, Connect(test.inputKey, test.inputConfig))
	}
}

func TestDisconnect(t *testing.T) {

	pool = map[string]*Connection{
		"test": &Connection{
			config: NewMySQLConfig(MySQLDatabase(databaseName)),
		},
	}
	_ = pool["test"].Connect()

	tests := []struct {
		inputKey string
		err      error
	}{
		{
			inputKey: "test",
			err:      nil,
		},
		{
			inputKey: "notExist",
			err:      errors.New("notExist database connection not exist"),
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.err, Disconnect(test.inputKey))
	}
}

func TestGetDB(t *testing.T) {

	pool = map[string]*Connection{
		"test": &Connection{
			config: NewMySQLConfig(MySQLDatabase(databaseName)),
		},
	}
	_ = pool["test"].Connect()

	tests := []struct {
		inputKey string
		want     *gorm.DB
	}{
		{
			inputKey: "test",
			want:     pool["test"].DB,
		},
		{
			inputKey: "notExist",
			want:     nil,
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.want, GetDB(test.inputKey))
	}
}
