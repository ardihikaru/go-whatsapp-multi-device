package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/satumedishub/sea-cucumber-api-service/internal/app"
	"github.com/satumedishub/sea-cucumber-api-service/internal/config"
	"github.com/satumedishub/sea-cucumber-api-service/internal/logger"
	"github.com/satumedishub/sea-cucumber-api-service/internal/router"
)

// Version sets the default build version
var Version = "development"

func main() {
	// loads configuration
	cfg, err := config.Get()
	if err != nil {
		app.FatalOnError(err, "error loading configuration")
	}

	// configures logger
	log, err := logger.New(cfg.LogLevel, cfg.LogFormat)
	if err != nil {
		app.FatalOnError(err, "failed to prepare the logger")
	}

	// shows the build version
	log.Info("starting Sea Cucumber API service. ",
		zap.String("Version", Version),
		zap.String("BuildMode", cfg.BuildMode),
	)

	// gracefully exit on keyboard interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// initializes persistent store
	db := app.InitializeDB(cfg, log)

	// initializes JWT Authenticator
	tokenAuth := app.GetTokenAuthentication(cfg, log)

	// initializes http client
	httpClient := app.BuildHttpClient()

	// initializes the dependency parameters
	deps := &app.Dependencies{
		Config:     cfg,
		DB:         db,
		Log:        log,
		TokenAuth:  tokenAuth,
		HttpClient: httpClient,
	}

	// starts the api server
	initializeHandler(deps, cfg.Address, cfg.Port)

	// logs that application is ready
	log.Info("preparing to serve the request in => " + fmt.Sprintf("%s:%v", cfg.Address, cfg.Port))

	// shutdowns the application
	<-c
	log.Info("gracefully shutting down the system")
	os.Exit(0)
}

func initializeHandler(deps *app.Dependencies, address string, port int) {
	r := router.GetRouter(deps)
	go func() {
		// stops the application if any error found
		if err := http.ListenAndServe(fmt.Sprintf("%s:%v", address, port), r); err != nil {
			app.FatalOnError(err, "failed to start server")
			os.Exit(1)
		}
	}()
}
