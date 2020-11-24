package verifiable

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitCreator(t *testing.T) {

	tests := []struct {
		input string
		want  *IDCreator
	}{
		{
			input: "",
			want: &IDCreator{
				signature:   newSignature(""),
				mutex:       sync.Mutex{},
				startTime:   toTime(time.Date(2014, 9, 1, 0, 0, 0, 0, time.UTC)),
				elapsedTime: 0,
				sequence:    255,
			},
		},
	}

	originalIDCreator := idCreator

	for _, test := range tests {

		InitCreator(test.input)

		assert.Equal(t, test.want, idCreator)
	}

	idCreator = originalIDCreator
}

func TestNextID(t *testing.T) {

	tests := []struct {
		input *IDCreator
		want  uint64
	}{
		{
			input: NewIDCreator(Settings{
				StartTime: time.Now(),
				SecretKey: " ",
			}),
			want: 0x10021,
		},
	}

	originalIDCreator := idCreator

	for _, test := range tests {

		idCreator = test.input
		assert.Equal(t, test.want, NextID())
	}

	idCreator = originalIDCreator
}

func TestNextHexString(t *testing.T) {

	tests := []struct {
		input *IDCreator
		want  string
	}{
		{
			input: NewIDCreator(Settings{
				StartTime: time.Now(),
				SecretKey: " ",
			}),
			want: "10021",
		},
	}

	originalIDCreator := idCreator

	for _, test := range tests {

		idCreator = test.input
		assert.Equal(t, test.want, NextHexString())
	}

	idCreator = originalIDCreator
}

func TestNextDecimalString(t *testing.T) {

	tests := []struct {
		input *IDCreator
		want  string
	}{
		{
			input: NewIDCreator(Settings{
				StartTime: time.Now(),
				SecretKey: " ",
			}),
			want: "65569",
		},
	}

	originalIDCreator := idCreator

	for _, test := range tests {

		idCreator = test.input
		assert.Equal(t, test.want, NextDecimalString())
	}

	idCreator = originalIDCreator
}
