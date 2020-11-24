package cache

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSetLogger(t *testing.T) {

	tests := []struct {
		input *logrus.Logger
	}{
		{
			input: func() *logrus.Logger {

				log := logrus.New()
				log.SetLevel(logrus.TraceLevel)
				return log
			}(),
		},
	}

	for _, test := range tests {

		SetLogger(test.input)
		assert.Equal(t, logger, test.input)
	}
}
