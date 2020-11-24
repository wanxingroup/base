package verifiable

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSignatureSign(t *testing.T) {

	tests := []struct {
		signature signature
		input     uint64
		want      uint64
	}{
		{
			signature: newSignature(""),
			input:     0x100,
			want:      0x1,
		},
		{
			signature: newSignature(""),
			input:     0xff00,
			want:      0xff,
		},
		{
			signature: newSignature(""),
			input:     0xffff00,
			want:      0x0,
		},
		{
			signature: newSignature(" "),
			input:     0x100,
			want:      0x21,
		},
		{
			signature: newSignature("  "),
			input:     0x100,
			want:      0x1,
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.want, test.signature.sign(test.input), test)
	}
}

func TestVerify(t *testing.T) {

	tests := []struct {
		input     uint64
		secretKey string
		want      bool
	}{
		{
			input:     0x101,
			secretKey: "",
			want:      true,
		},
		{
			input:     0xffff,
			secretKey: "",
			want:      true,
		},
		{
			input:     0xffff00,
			secretKey: "",
			want:      true,
		},
		{
			input:     0x121,
			secretKey: " ",
			want:      true,
		},
		{
			input:     0x101,
			secretKey: "  ",
			want:      true,
		},
		{
			input:     0x102,
			secretKey: "",
			want:      false,
		},
		{
			input:     0x103,
			secretKey: "  ",
			want:      false,
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.want, Verify(test.input, test.secretKey), test)
	}
}

func TestNewIDCreator(t *testing.T) {

	startTime := time.Now()
	tests := []struct {
		input Settings
		want  *IDCreator
	}{
		{
			input: Settings{
				StartTime: startTime,
				SecretKey: "",
			},
			want: &IDCreator{
				signature:   signature{},
				mutex:       sync.Mutex{},
				startTime:   startTime.UnixNano() / timeUnit,
				elapsedTime: 0,
				sequence:    1<<BitLengthSequence - 1,
			},
		},
		{
			input: Settings{
				SecretKey: "",
			},
			want: &IDCreator{
				signature:   signature{},
				mutex:       sync.Mutex{},
				startTime:   time.Date(2014, 9, 1, 0, 0, 0, 0, time.UTC).UnixNano() / timeUnit,
				elapsedTime: 0,
				sequence:    1<<BitLengthSequence - 1,
			},
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.want, NewIDCreator(test.input), test)
	}
}

func TestIDCreator_NextID(t *testing.T) {

	tests := []struct {
		input Settings
		want  []uint64
	}{
		{
			input: Settings{
				StartTime: time.Now(),
				SecretKey: "",
			},
			want: []uint64{0x10001},
		},
		{
			input: Settings{
				StartTime: time.Now(),
				SecretKey: "",
			},
			want: []uint64{0x10001, 0x10100, 0x10203, 0x10302, 0x10405},
		},
		{
			input: Settings{
				StartTime: time.Now(),
				SecretKey: " ",
			},
			want: []uint64{0x10021, 0x10120, 0x10223, 0x10322, 0x10425},
		},
	}

	for _, test := range tests {

		c := NewIDCreator(test.input)
		for _, want := range test.want {
			assert.Equal(t, want, c.NextID())
		}
	}
}
