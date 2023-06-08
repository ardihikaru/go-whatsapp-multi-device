package app

import (
	"net/http"

	"github.com/ardihikaru/go-modules/pkg/logger"
	botHook "github.com/ardihikaru/go-modules/pkg/whatsappbot/wawebhook"

	"github.com/ardihikaru/go-whatsapp-multi-device/internal/config"
	"github.com/ardihikaru/go-whatsapp-multi-device/internal/storage"
)

// Dependencies holds the primitives and structs and/or interfaces that are required
// for the application's business logic.
type Dependencies struct {
	Config      *config.Config
	DB          *storage.DataStoreMongo
	Log         *logger.Logger
	HttpClient  *http.Client
	WhatsAppBot *botHook.WaManager
	BotClients  *botHook.BotClientList
}
