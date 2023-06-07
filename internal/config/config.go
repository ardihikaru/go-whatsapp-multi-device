package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
)

const (
	buildModeEnv          = "BUILD_MODE"
	addressEnv            = "ADDRESS"
	portEnv               = "PORT"
	logLevelEnv           = "LOG_LEVEL"
	logFormatEnv          = "LOG_FORMAT"
	corsAllowOriginsEnv   = "CORS_ALLOW_ORIGINS"
	corsAllowHeadersEnv   = "CORS_ALLOW_HEADERS"
	corsExposedHeadersEnv = "CORS_EXPOSED_HEADERS"
	dbConnUriEnv          = "DB_CONN_URI"
	dbNameEnv             = "DB_NAME"
	dbHBIntervalEnv       = "DB_HEARTBEAT_INTERVAL"
	dbLocalThresholdEnv   = "DB_LOCAL_THRESHOLD"
	dbServerSelTimeoutEnv = "DB_SERVER_SELECTION_TIMEOUT"
	dbMaxPoolSizeEnv      = "DB_MAX_POOL_SIZE"
	jwtSecretEnv          = "JWT_SECRET"
	jwtAlgorithmEnv       = "JWT_ALGORITHM"
	jwtExpiredInSecEnv    = "JWT_EXPIRED_IN_SEC"
)

var defaultCORSAllowOrigins = []string{"*"}
var defaultCORSAllowHeaders = []string{"*"}
var defaultCORSExposedHeaders = []string{"*"}

const (
	defaultBuildMode           = "dev"
	defaultLogLevel            = "info"
	defaultLogFormat           = "json"
	defaultAddress             = "0.0.0.0"
	defaultPort                = 80
	defaultDbConnURI           = "mongodb://localhost:27017"
	defaultDbName              = "apiDb"
	defaultDbConnTimeout       = 30 * time.Second
	defaultDbHeartbeatInterval = 10 * time.Second
	defaultDbLocalThreshold    = 15 * time.Second
	defaultDbServerSelTimeout  = 30 * time.Second
	defaultDbMaxPoolSize       = 100
	defaultJWTSecret           = "secret"
	defaultJWTAlgorithm        = "HS256"
	defaultJWTExpInSec         = 3600 // token will be expired in 1 hour
)

// Config provides all the possible configurations
// FYI: viper will read the field name, such as `DBUser`, instead of `DB_USER`
type Config struct {
	BuildMode           string                 `config:"BUILD_MODE" validate:"oneof=dev stag prod"`
	Address             string                 `config:"ADDRESS"`
	Port                int                    `config:"PORT"`
	CORSAllowOrigins    []string               `config:"CORS_ALLOW_ORIGINS"`
	CORSAllowHeaders    []string               `config:"CORS_ALLOW_HEADERS"`
	CORSExposedHeaders  []string               `config:"CORS_EXPOSED_HEADERS"`
	LogLevel            string                 `config:"LOG_LEVEL" validate:"oneof=debug info warn error fatal panic"`
	LogFormat           string                 `config:"LOG_FORMAT" validate:"oneof=text console json"`
	DbConnURI           string                 `config:"DB_CONN_URI"`
	DBName              string                 `config:"DB_NAME"`
	DbConnTimeout       time.Duration          `config:"DB_CONN_TIMEOUT"`
	DbHeartBeatInterval time.Duration          `config:"DB_HEARTBEAT_INTERVAL"`
	DbLocalThreshold    time.Duration          `config:"DB_LOCAL_THRESHOLD"`
	DbServerSelTimeout  time.Duration          `config:"DB_SERVER_SELECTION_TIMEOUT"`
	DbMaxPoolSize       uint64                 `config:"DB_MAX_POOL_SIZE"`
	JWTSecret           string                 `config:"JWT_SECRET"`
	JWTAlgorithm        jwa.SignatureAlgorithm `config:"JWT_ALGORITHM"`
	JWTExpiredInSec     int64                  `config:"JWT_EXPIRED_IN_SEC"`
}

// Get returns the configuration loaded from the environment variable.
// Default values are initialized by the hardcoded values in case the environment variable is not provided
func Get() (*Config, error) {
	var err error

	c := Config{
		BuildMode:           defaultBuildMode,
		Address:             defaultAddress,
		Port:                defaultPort,
		LogLevel:            defaultLogLevel,
		LogFormat:           defaultLogFormat,
		CORSAllowOrigins:    defaultCORSAllowOrigins,
		CORSAllowHeaders:    defaultCORSAllowHeaders,
		CORSExposedHeaders:  defaultCORSExposedHeaders,
		DbConnURI:           defaultDbConnURI,
		DBName:              defaultDbName,
		DbConnTimeout:       defaultDbConnTimeout,
		DbHeartBeatInterval: defaultDbHeartbeatInterval,
		DbLocalThreshold:    defaultDbLocalThreshold,
		DbServerSelTimeout:  defaultDbServerSelTimeout,
		DbMaxPoolSize:       defaultDbMaxPoolSize,
		JWTSecret:           defaultJWTSecret,
		JWTAlgorithm:        defaultJWTAlgorithm,
		JWTExpiredInSec:     defaultJWTExpInSec,
	}

	// try to find the variable inside the environment variable
	err = c.validateAndLoadSystemEnv()
	if err != nil {
		return &c, err
	}

	return &c, nil
}

// validateAndLoadSystemEnv loads any detected environment variable
func (c *Config) validateAndLoadSystemEnv() error {
	var err error

	if os.Getenv(buildModeEnv) != "" {
		c.BuildMode = os.Getenv(buildModeEnv)
	}
	if os.Getenv(addressEnv) != "" {
		c.Address = os.Getenv(addressEnv)
	}
	if os.Getenv(portEnv) != "" {
		c.Port, err = strconv.Atoi(os.Getenv(portEnv))
		if err != nil {
			return err
		}
	}
	if os.Getenv(logLevelEnv) != "" {
		c.LogLevel = os.Getenv(logLevelEnv)
	}
	if os.Getenv(logFormatEnv) != "" {
		c.LogFormat = os.Getenv(logFormatEnv)
	}
	if os.Getenv(corsAllowOriginsEnv) != "" {
		corsAllowOriginsStr := os.Getenv(corsAllowOriginsEnv)
		c.CORSAllowOrigins = strings.Split(corsAllowOriginsStr, ",")
	}
	if os.Getenv(corsAllowHeadersEnv) != "" {
		corsAllowHeadersStr := os.Getenv(corsAllowHeadersEnv)
		c.CORSAllowHeaders = strings.Split(corsAllowHeadersStr, ",")
	}
	if os.Getenv(corsExposedHeadersEnv) != "" {
		corsExposedHeadersStr := os.Getenv(corsExposedHeadersEnv)
		c.CORSExposedHeaders = strings.Split(corsExposedHeadersStr, ",")
	}
	if os.Getenv(dbConnUriEnv) != "" {
		c.DbConnURI = os.Getenv(dbConnUriEnv)
	}
	if os.Getenv(dbNameEnv) != "" {
		c.DBName = os.Getenv(dbNameEnv)
	}
	if os.Getenv(dbHBIntervalEnv) != "" {
		c.DbHeartBeatInterval, err = time.ParseDuration(os.Getenv(dbHBIntervalEnv))
		if err != nil {
			return err
		}
	}
	if os.Getenv(dbLocalThresholdEnv) != "" {
		c.DbLocalThreshold, err = time.ParseDuration(os.Getenv(dbLocalThresholdEnv))
		if err != nil {
			return err
		}
	}
	if os.Getenv(dbServerSelTimeoutEnv) != "" {
		c.DbServerSelTimeout, err = time.ParseDuration(os.Getenv(dbServerSelTimeoutEnv))
		if err != nil {
			return err
		}
	}
	if os.Getenv(dbMaxPoolSizeEnv) != "" {
		c.DbMaxPoolSize, err = strconv.ParseUint(os.Getenv(dbMaxPoolSizeEnv), 10, 64)
		if err != nil {
			return err
		}
	}
	if os.Getenv(jwtSecretEnv) != "" {
		c.JWTSecret = os.Getenv(jwtSecretEnv)
	}
	if os.Getenv(jwtAlgorithmEnv) != "" {
		var alg jwa.SignatureAlgorithm
		c.JWTAlgorithm = jwa.SignatureAlgorithm(os.Getenv(jwtAlgorithmEnv))
		if err := alg.Accept(c.JWTAlgorithm); err != nil {
			return err
		}
	}
	if os.Getenv(jwtExpiredInSecEnv) != "" {
		c.DbMaxPoolSize, err = strconv.ParseUint(os.Getenv(jwtExpiredInSecEnv), 10, 64)
		if err != nil {
			return err
		}
	}

	return nil
}
