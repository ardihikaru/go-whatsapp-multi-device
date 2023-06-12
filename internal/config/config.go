package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
)

const (
	buildModeEnv              = "BUILD_MODE"
	addressEnv                = "ADDRESS"
	portEnv                   = "PORT"
	logLevelEnv               = "LOG_LEVEL"
	logFormatEnv              = "LOG_FORMAT"
	corsAllowOriginsEnv       = "CORS_ALLOW_ORIGINS"
	corsAllowHeadersEnv       = "CORS_ALLOW_HEADERS"
	corsExposedHeadersEnv     = "CORS_EXPOSED_HEADERS"
	dbConnUriEnv              = "DB_CONN_URI"
	dbNameEnv                 = "DB_NAME"
	dbHBIntervalEnv           = "DB_HEARTBEAT_INTERVAL"
	dbLocalThresholdEnv       = "DB_LOCAL_THRESHOLD"
	dbServerSelTimeoutEnv     = "DB_SERVER_SELECTION_TIMEOUT"
	dbMaxPoolSizeEnv          = "DB_MAX_POOL_SIZE"
	jwtSecretEnv              = "JWT_SECRET"
	jwtAlgorithmEnv           = "JWT_ALGORITHM"
	jwtExpiredInSecEnv        = "JWT_EXPIRED_IN_SEC"
	whatsappDbNameEnv         = "WHATSAPP_DB_NAME"
	whatsappQrCodeDirEnv      = "WHATSAPP_QC_CODE_DIR"
	whatsappWebhookEnv        = "WHATSAPP_WEBHOOK"
	whatsappWebhookEnabledEnv = "WHATSAPP_WEBHOOK_ENABLED"
	WhatsappQrToTerminalEnv   = "WHATSAPP_QR_TO_TERMINAL"
	whatsappWebhookEchoEnv    = "WHATSAPP_WEBHOOK_ECHO"
	httpClientTlsEnv          = "HTTP_CLIENT_TLS"
)

var defaultCORSAllowOrigins = []string{"*"}
var defaultCORSAllowHeaders = []string{"*"}
var defaultCORSExposedHeaders = []string{"*"}

// Config provides all the possible configurations
// FYI: viper will read the field name, such as `DBUser`, instead of `DB_USER`
type Config struct {
	BuildMode              string                 `config:"BUILD_MODE" validate:"oneof=dev stag prod"`
	Address                string                 `config:"ADDRESS"`
	Port                   int                    `config:"PORT"`
	CORSAllowOrigins       []string               `config:"CORS_ALLOW_ORIGINS"`
	CORSAllowHeaders       []string               `config:"CORS_ALLOW_HEADERS"`
	CORSExposedHeaders     []string               `config:"CORS_EXPOSED_HEADERS"`
	LogLevel               string                 `config:"LOG_LEVEL" validate:"oneof=debug info warn error fatal panic"`
	LogFormat              string                 `config:"LOG_FORMAT" validate:"oneof=text console json"`
	DbConnURI              string                 `config:"DB_CONN_URI"`
	DBName                 string                 `config:"DB_NAME"`
	DbConnTimeout          time.Duration          `config:"DB_CONN_TIMEOUT"`
	DbHeartBeatInterval    time.Duration          `config:"DB_HEARTBEAT_INTERVAL"`
	DbLocalThreshold       time.Duration          `config:"DB_LOCAL_THRESHOLD"`
	DbServerSelTimeout     time.Duration          `config:"DB_SERVER_SELECTION_TIMEOUT"`
	DbMaxPoolSize          uint64                 `config:"DB_MAX_POOL_SIZE"`
	JWTSecret              string                 `config:"JWT_SECRET"`
	JWTAlgorithm           jwa.SignatureAlgorithm `config:"JWT_ALGORITHM"`
	JWTExpiredInSec        int64                  `config:"JWT_EXPIRED_IN_SEC"`
	WhatsappDbName         string                 `config:"WHATSAPP_DB_NAME"`
	WhatsappQrCodeDir      string                 `config:"WHATSAPP_QC_CODE_DIR"`
	WhatsappWebhook        string                 `config:"WHATSAPP_WEBHOOK"`
	WhatsappWebhookEnabled bool                   `config:"WHATSAPP_WEBHOOK_ENABLED"`
	WhatsappQrToTerminal   bool                   `config:"WHATSAPP_QR_TO_TERMINAL"`
	WhatsappWebhookEcho    bool                   `config:"WHATSAPP_WEBHOOK_ECHO"`
	HttpClientTLS          bool                   `config:"HTTP_CLIENT_TLS"`
}

// Get returns the configuration loaded from the environment variable.
// Default values are initialized by the hardcoded values in case the environment variable is not provided
func Get() (*Config, error) {
	var err error

	c := Config{
		BuildMode:              "dev",
		Address:                "0.0.0.0",
		Port:                   80,
		LogLevel:               "info",
		LogFormat:              "json",
		CORSAllowOrigins:       defaultCORSAllowOrigins,
		CORSAllowHeaders:       defaultCORSAllowHeaders,
		CORSExposedHeaders:     defaultCORSExposedHeaders,
		DbConnURI:              "mongodb://localhost:27017",
		DBName:                 "whatsappDb",
		DbConnTimeout:          30 * time.Second,
		DbHeartBeatInterval:    10 * time.Second,
		DbLocalThreshold:       15 * time.Second,
		DbServerSelTimeout:     30 * time.Second,
		DbMaxPoolSize:          100,
		JWTSecret:              "secret",
		JWTAlgorithm:           "HS256",
		JWTExpiredInSec:        3600, // token will be expired in 1 hour,
		WhatsappDbName:         "./data/sqlitedb/datastore",
		WhatsappQrCodeDir:      "./data/qrcode",
		WhatsappWebhookEnabled: false,
		WhatsappQrToTerminal:   true,
		WhatsappWebhook:        "http://localhost:8500/webhook",
		WhatsappWebhookEcho:    true,
		HttpClientTLS:          true,
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

	if os.Getenv(whatsappDbNameEnv) != "" {
		c.WhatsappDbName = os.Getenv(whatsappDbNameEnv)
	}
	if os.Getenv(whatsappWebhookEnabledEnv) != "" {
		// validates the boolean value
		boolWhatsappWebhookEnabled, err := strconv.ParseBool(os.Getenv(whatsappWebhookEnabledEnv))
		if err != nil {
			return err
		}
		c.WhatsappWebhookEnabled = boolWhatsappWebhookEnabled
	}
	if os.Getenv(WhatsappQrToTerminalEnv) != "" {
		// validates the boolean value
		boolWhatsappQrToTerminal, err := strconv.ParseBool(os.Getenv(WhatsappQrToTerminalEnv))
		if err != nil {
			return err
		}
		c.WhatsappQrToTerminal = boolWhatsappQrToTerminal
	}
	if os.Getenv(whatsappQrCodeDirEnv) != "" {
		c.WhatsappQrCodeDir = os.Getenv(whatsappQrCodeDirEnv)
	}
	if os.Getenv(whatsappWebhookEnv) != "" {
		c.WhatsappWebhook = os.Getenv(whatsappWebhookEnv)
	}
	if os.Getenv(whatsappWebhookEchoEnv) != "" {
		// validates the boolean value
		boolWhatsappWebhookEcho, err := strconv.ParseBool(os.Getenv(whatsappWebhookEchoEnv))
		if err != nil {
			return err
		}
		c.WhatsappWebhookEcho = boolWhatsappWebhookEcho
	}

	// http client
	if os.Getenv(httpClientTlsEnv) != "" {
		// validates the boolean value
		boolHttpClientTLS, err := strconv.ParseBool(os.Getenv(httpClientTlsEnv))
		if err != nil {
			return err
		}
		c.HttpClientTLS = boolHttpClientTLS
	}

	return nil
}
