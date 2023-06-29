package app

import (
	"context"
	"fmt"
	"github.com/ardihikaru/go-modules/pkg/utils/httputils"

	deviceSvc "github.com/ardihikaru/go-whatsapp-multi-device/internal/service/device"
	sessionSvc "github.com/ardihikaru/go-whatsapp-multi-device/internal/service/session"
)

// AutoStartLoggedSessions starts all logged sessions
func AutoStartLoggedSessions(deps *Dependencies) {
	// initializes services
	deviceService := deviceSvc.NewService(deps.DB, deps.Log)
	sessionService := sessionSvc.NewService(deviceService, deps.Log, deps.WhatsAppBot, deps.HttpClient,
		deps.Config.WhatsappImageDir, deps.Config.WhatsappQrCodeDir,
		deps.Config.WhatsappWebhookEcho, deps.Config.WhatsappWebhookEnabled, deps.Config.WhatsappQrToTerminal,
		deps.BotClients)

	// builds query parameters
	params := httputils.GetQueryParams{
		Limit:  1000,
		Offset: 0,
		Order:  "ASC",
		Sort:   "",
		Search: "",
	}

	_, devices, err := deviceService.GetDevices(context.Background(), params)
	if err != nil {
		// do nothing here
	} else {
		// loop and start sessions
		for _, device := range devices {
			// removes `+` symbol if exists
			phone := device.Phone
			if phone[0:1] == "+" {
				phone = phone[1:]
			}

			err := sessionService.New(context.Background(), phone)
			if err != nil {
				deps.Log.Warn(fmt.Sprintf("opening session for phone [%s] failed", phone))
			}
		}
	}
}
