package config

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	var err error

	defer os.Clearenv()

	// sets environment variable
	err = os.Setenv(addressEnv, defaultAddress)
	assert.NoError(t, err)
	err = os.Setenv(portEnv, strconv.Itoa(defaultPort))
	assert.NoError(t, err)
	err = os.Setenv(logFormatEnv, defaultLogFormat)
	assert.NoError(t, err)
	err = os.Setenv(logLevelEnv, defaultLogLevel)
	assert.NoError(t, err)

	// loads an env file for test purpose
	c, err := Get()
	assert.NoError(t, err)

	// validates each environment variable
	assert.Equal(t, c.Address, defaultAddress)
	assert.Equal(t, c.Port, defaultPort)
	assert.Equal(t, c.LogLevel, defaultLogLevel)
	assert.Equal(t, c.LogFormat, defaultLogFormat)
	assert.Equal(t, c.DBName, defaultDbName)
}
