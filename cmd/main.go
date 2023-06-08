package main

import (
	"fmt"
	botHook "github.com/ardihikaru/go-modules/pkg/whatsappbot/wawebhook"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ardihikaru/go-modules/pkg/logger"
	e "github.com/ardihikaru/go-modules/pkg/utils/error"
	hc "github.com/ardihikaru/go-modules/pkg/utils/httpclient"
	wBot "github.com/ardihikaru/go-modules/pkg/whatsappbot"
	"go.uber.org/zap"

	"github.com/ardihikaru/go-whatsapp-multi-device/internal/app"
	"github.com/ardihikaru/go-whatsapp-multi-device/internal/config"
	"github.com/ardihikaru/go-whatsapp-multi-device/internal/router"
)

// Version sets the default build version
var Version = "development"

func main() {
	// loads configuration
	cfg, err := config.Get()
	if err != nil {
		e.FatalOnError(err, "error loading configuration")
	}

	// configures logger
	log, err := logger.New(cfg.LogLevel, cfg.LogFormat)
	if err != nil {
		e.FatalOnError(err, "failed to prepare the logger")
	}

	// shows the build version
	log.Info("starting WhatsApp multi-device API service. ",
		zap.String("Version", Version),
		zap.String("BuildMode", cfg.BuildMode),
	)

	// gracefully exit on keyboard interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// initializes persistent store
	db := app.InitializeDB(cfg, log)

	// initializes http client
	httpClient := hc.BuildHttpClient(cfg.HttpClientTLS)

	// creates list to store created whatsapp bot clients
	botClients := make(botHook.BotClientList)

	// initializes whatsapp bot
	whatsAppBot := wBot.InitWhatsappContainer(cfg.WhatsappDbName, log)

	// initializes the dependency parameters
	deps := &app.Dependencies{
		Config:      cfg,
		DB:          db,
		Log:         log,
		HttpClient:  httpClient,
		WhatsAppBot: whatsAppBot,
		BotClients:  &botClients,
	}

	// starts the api server
	initializeHandler(deps, cfg.Address, cfg.Port)

	// logs that application is ready
	log.Info("preparing to serve the request in => " + fmt.Sprintf("%s:%v", cfg.Address, cfg.Port))

	// shutdowns the RESTApi Server
	<-c
	log.Info("gracefully shutting down the system")

	// exit app
	os.Exit(0)
}

func initializeHandler(deps *app.Dependencies, address string, port int) {
	r := router.GetRouter(deps)
	go func() {
		// stops the application if any error found
		if err := http.ListenAndServe(fmt.Sprintf("%s:%v", address, port), r); err != nil {
			e.FatalOnError(err, "failed to start server")
			os.Exit(1)
		}
	}()
}
