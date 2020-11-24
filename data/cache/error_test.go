package cache

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNotHaveConfigError(t *testing.T) {

	tests := []struct {
		input error
		want  bool
	}{
		{
			input: errors.New("other error"),
			want:  false,
		},
		{
			input: errorNotHaveConfig,
			want:  true,
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.want, IsNotHaveConfigError(test.input))
	}
}

func TestIsNotHaveConnectionError(t *testing.T) {

	tests := []struct {
		input error
		want  bool
	}{
		{
			input: errors.New("other error"),
			want:  false,
		},
		{
			input: errorNotHaveConnection,
			want:  true,
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.want, IsNotHaveConnectionError(test.input))
	}
}

func TestIsPingFailedError(t *testing.T) {

	tests := []struct {
		input error
		want  bool
	}{
		{
			input: errors.New("other error"),
			want:  false,
		},
		{
			input: errorPingFailed,
			want:  true,
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.want, IsPingFailedError(test.input))
	}
}

func TestIsPingResultNotExpectError(t *testing.T) {

	tests := []struct {
		input error
		want  bool
	}{
		{
			input: errors.New("other error"),
			want:  false,
		},
		{
			input: errorPingResultNotExpect,
			want:  true,
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.want, IsPingResultNotExpectError(test.input))
	}
}
