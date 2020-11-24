package database

import (
	"math/rand"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Transaction(t *testing.T) {

	rand.Seed(time.Now().UnixNano())
	repo1 := &Connection{config: NewMySQLConfig(MySQLDatabase(databaseName))}
	err := repo1.Connect()
	assert.Nil(t, err)

	type testModel struct {
		Id   uint64 `gorm:"primary_key;auto_increment:false"`
		Name string
	}

	repo1.AutoMigrate(&testModel{})

	var m = testModel{Id: rand.Uint64(), Name: "test"}
	transaction := repo1.Begin()
	transaction.Model(&testModel{}).Create(m)
	assert.False(t, repo1.Model(&testModel{}).NewRecord(m))
	transaction.Commit()

	var readModel testModel
	repo1.Where(&testModel{Id: m.Id}).First(&readModel)
	assert.EqualValues(t, m, readModel)
}

func TestMySQLHost(t *testing.T) {

	tests := []struct {
		inputOption MySQLOption
		inputConfig MySQLConfig
		want        MySQLConfig
	}{
		{
			inputOption: MySQLHost("localhost"),
			inputConfig: defaultMySQLConfig,
			want: func() MySQLConfig {

				config := defaultMySQLConfig
				config.host = "localhost"
				return config
			}(),
		},
	}

	for _, test := range tests {

		test.inputOption(&test.inputConfig)
		assert.Equal(t, test.want, test.inputConfig)
	}
}

func TestMySQLPort(t *testing.T) {

	tests := []struct {
		inputOption MySQLOption
		inputConfig MySQLConfig
		want        MySQLConfig
	}{
		{
			inputOption: MySQLPort(3308),
			inputConfig: defaultMySQLConfig,
			want: func() MySQLConfig {

				config := defaultMySQLConfig
				config.port = 3308
				return config
			}(),
		},
	}

	for _, test := range tests {

		test.inputOption(&test.inputConfig)
		assert.Equal(t, test.want, test.inputConfig)
	}
}

func TestMySQLUsername(t *testing.T) {

	tests := []struct {
		inputOption MySQLOption
		inputConfig MySQLConfig
		want        MySQLConfig
	}{
		{
			inputOption: MySQLUsername("testUser"),
			inputConfig: defaultMySQLConfig,
			want: func() MySQLConfig {

				config := defaultMySQLConfig
				config.user = "testUser"
				return config
			}(),
		},
	}

	for _, test := range tests {

		test.inputOption(&test.inputConfig)
		assert.Equal(t, test.want, test.inputConfig)
	}
}

func TestMySQLPassword(t *testing.T) {

	tests := []struct {
		inputOption MySQLOption
		inputConfig MySQLConfig
		want        MySQLConfig
	}{
		{
			inputOption: MySQLPassword("password"),
			inputConfig: defaultMySQLConfig,
			want: func() MySQLConfig {

				config := defaultMySQLConfig
				config.password = "password"
				return config
			}(),
		},
	}

	for _, test := range tests {

		test.inputOption(&test.inputConfig)
		assert.Equal(t, test.want, test.inputConfig)
	}
}

func TestMySQLDatabase(t *testing.T) {

	tests := []struct {
		inputOption MySQLOption
		inputConfig MySQLConfig
		want        MySQLConfig
	}{
		{
			inputOption: MySQLDatabase("db"),
			inputConfig: defaultMySQLConfig,
			want: func() MySQLConfig {

				config := defaultMySQLConfig
				config.database = "db"
				return config
			}(),
		},
	}

	for _, test := range tests {

		test.inputOption(&test.inputConfig)
		assert.Equal(t, test.want, test.inputConfig)
	}
}

func TestMySQLLogMode(t *testing.T) {

	tests := []struct {
		inputOption MySQLOption
		inputConfig MySQLConfig
		want        MySQLConfig
	}{
		{
			inputOption: nil,
			inputConfig: defaultMySQLConfig,
			want: func() MySQLConfig {

				config := defaultMySQLConfig
				config.logMode = false
				return config
			}(),
		},
		{
			inputOption: MySQLLogMode(true),
			inputConfig: defaultMySQLConfig,
			want: func() MySQLConfig {

				config := defaultMySQLConfig
				config.logMode = true
				return config
			}(),
		},
		{
			inputOption: MySQLLogMode(false),
			inputConfig: defaultMySQLConfig,
			want: func() MySQLConfig {

				config := defaultMySQLConfig
				config.logMode = false
				return config
			}(),
		},
	}

	for _, test := range tests {

		if test.inputOption != nil {
			test.inputOption(&test.inputConfig)
		}

		assert.Equal(t, test.want, test.inputConfig)
	}
}

func TestConnection_Connect(t *testing.T) {

	tests := []struct {
		inputOption []MySQLOption
		hasError    bool
	}{
		{
			inputOption: []MySQLOption{},
			hasError:    false,
		},
		{
			inputOption: []MySQLOption{
				MySQLPort(3308),
				MySQLDatabase("test-connect-failed"),
			},
			hasError: true,
		},
		{
			inputOption: []MySQLOption{},
			hasError:    false,
		},
	}

	for _, test := range tests {

		config := NewMySQLConfig(test.inputOption...)

		if test.hasError {
			assert.NotNil(t, config.Connect())
		} else {
			assert.Nil(t, config.Connect())
		}
	}
}

func TestConnection_isConnected(t *testing.T) {

	disconnected, _ := gorm.Open("mysql", defaultMySQLConfig.getConnectionString())
	_ = disconnected.Close()
	connected, _ := gorm.Open("mysql", defaultMySQLConfig.getConnectionString())
	tests := []struct {
		inputConnection *Connection
		want            bool
	}{
		{
			inputConnection: &Connection{
				config: NewMySQLConfig(),
			},
			want: false,
		},
		{
			inputConnection: &Connection{
				DB:     &gorm.DB{},
				config: NewMySQLConfig(),
			},
			want: false,
		},
		{
			inputConnection: &Connection{
				DB:     disconnected,
				config: NewMySQLConfig(),
			},
			want: false,
		},
		{
			inputConnection: &Connection{
				DB:     connected,
				config: NewMySQLConfig(),
			},
			want: true,
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.want, test.inputConnection.isConnected())
	}
}

func TestConnection_Close(t *testing.T) {

	tests := []struct {
		input *Connection
		err   error
	}{
		{
			input: func() *Connection {
				testDatabase := &Connection{
					config: NewMySQLConfig(MySQLDatabase(databaseName)),
				}
				testDatabase.DB, _ = gorm.Open("mysql", NewMySQLConfig(MySQLDatabase(databaseName)).getConnectionString())
				return testDatabase
			}(),
			err: nil,
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.err, test.input.Close())
	}
}
