package cache

import (
	"errors"
)

var errorNotHaveConfig = errors.New("not set configuration")
var errorPingFailed = errors.New("ping failed")
var errorPingResultNotExpect = errors.New("ping result not expect")
var errorNotHaveConnection = errors.New("not have connection")

func IsNotHaveConfigError(err error) bool {
	return err == errorNotHaveConfig
}

func IsPingFailedError(err error) bool {
	return err == errorPingFailed
}

func IsPingResultNotExpectError(err error) bool {
	return err == errorPingResultNotExpect
}

func IsNotHaveConnectionError(err error) bool {
	return err == errorNotHaveConnection
}
