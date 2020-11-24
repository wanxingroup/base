package auth

import (
	"errors"
)

var errorEmptyString = errors.New("token string empty")
var errorTokenExpired = errors.New("token was expired")
var errorTokenInvalid = errors.New("token is invalid")

func IsTokenExpired(err error) bool {
	return err == errorTokenExpired
}

func IsTokenInvalid(err error) bool {
	return err == errorTokenInvalid
}
