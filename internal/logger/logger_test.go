package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		desc      string
		logLevel  string
		logFormat string
		isNil     assert.ValueAssertionFunc
		isError   assert.ErrorAssertionFunc
	}{
		{
			desc:      "tested with invalid log-level, returned an error",
			logLevel:  "unknown",
			logFormat: logFormatJSON,
			isNil:     assert.Nil,
			isError:   assert.Error,
		},
		{
			desc:      "tested with invalid log-level, returned an error",
			logLevel:  LogLevelError,
			logFormat: "plain",
			isNil:     assert.Nil,
			isError:   assert.Error,
		},
		{
			desc:      "tested with empty log-level, returned with no error. Got INFO log level as the default",
			logLevel:  "",
			logFormat: logFormatJSON,
			isNil:     assert.NotNil,
			isError:   assert.NoError,
		},
		{
			desc:      "tested with valid log-level, returned with no error. Set log level to [error]",
			logLevel:  LogLevelError,
			logFormat: logFormatJSON,
			isNil:     assert.NotNil,
			isError:   assert.NoError,
		},
		{
			desc:      "tested with valid log-level, returned with no error. Set log level to [warn]",
			logLevel:  LogLevelWarn,
			logFormat: logFormatJSON,
			isNil:     assert.NotNil,
			isError:   assert.NoError,
		},
		{
			desc:      "tested with valid log-level, returned with no error. Set log level to [debug]",
			logLevel:  LogLevelDebug,
			logFormat: logFormatJSON,
			isNil:     assert.NotNil,
			isError:   assert.NoError,
		},
		{
			desc:      "tested with valid log-level, returned with no error. Set log level to [fatal]",
			logLevel:  LogLevelFatal,
			logFormat: logFormatJSON,
			isNil:     assert.NotNil,
			isError:   assert.NoError,
		},
		{
			desc:      "tested with valid log-level, returned with no error. Set log level to [panic]",
			logLevel:  LogLevelPanic,
			logFormat: logFormatJSON,
			isNil:     assert.NotNil,
			isError:   assert.NoError,
		},
		{
			desc:      "tested with valid log-level, returned with no error. Set log level to [info]",
			logLevel:  LogLevelInfo,
			logFormat: logFormatJSON,
			isNil:     assert.NotNil,
			isError:   assert.NoError,
		},
		{
			desc:      "tested with legacy log-level, returned with no error. Set Logger to use console",
			logLevel:  LogLevelError,
			logFormat: logFormatText,
			isNil:     assert.NotNil,
			isError:   assert.NoError,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			actual, err := New(test.logLevel, test.logFormat)
			test.isNil(t, actual)
			test.isError(t, err)
		})
	}
}
