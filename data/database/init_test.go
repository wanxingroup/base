package database

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {

	InitDatabase()

	code := m.Run()

	if code == 0 {
		ReleaseDatabase()
	}

	os.Exit(code)
}

const databaseNamePrefix = "test"

var databaseName string

func InitDatabase() {

	rand.Seed(time.Now().UnixNano())

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("dir = %s\n", dir)
	for {
		dir = path.Dir(dir)
		if strings.LastIndex(dir, "/pkg") <= 0 {
			break
		}
	}

	fmt.Printf("processed dir = %s\n", dir)

	databaseName = fmt.Sprintf("%s_%d", databaseNamePrefix, rand.Uint64())

	logrus.WithField("database", databaseName).Debug("init database")

	db, err := getTestDatabaseConnection()
	if err != nil {
		panic(err)
	}

	err = db.Exec(fmt.Sprintf("CREATE DATABASE `%s` character set UTF8mb4 collate utf8mb4_general_ci", databaseName)).Error
	if err != nil {
		panic(err)
	}
	_ = db.Close()
}

func ReleaseDatabase() {

	var err error

	db, err := getTestDatabaseConnection()
	if err != nil {
		panic(err)
	}

	err = db.Exec(fmt.Sprintf("DROP DATABASE %s", databaseName)).Error
	if err != nil {
		logrus.WithField("error", err).Error("drop database error")
	}

	_ = db.Close()
}

func getTestDatabaseConnection() (db *gorm.DB, err error) {

	db, err = gorm.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&loc=Local&parseTime=true",
			"root",
			"",
			"localhost",
			3306,
		))
	return
}
