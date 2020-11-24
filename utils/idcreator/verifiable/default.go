package verifiable

import (
	"fmt"
)

const (
	DefaultSecretKey = ""
)

var idCreator *IDCreator

func init() {
	InitCreator(DefaultSecretKey)
}

func InitCreator(secretKey string) {
	idCreator = NewIDCreator(
		Settings{
			SecretKey: secretKey,
		},
	)
}

func NextID() uint64 {
	return idCreator.NextID()
}

func NextHexString() string {
	return fmt.Sprintf("%x", NextID())
}

func NextDecimalString() string {
	return fmt.Sprintf("%d", NextID())
}
