package handlers

import (
	"net/http"

	"github.com/ardihikaru/go-modules/pkg/logger"
	"github.com/ardihikaru/go-modules/pkg/utils/httputils"
	botHook "github.com/ardihikaru/go-modules/pkg/whatsappbot/wawebhook"
	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/ardihikaru/go-whatsapp-multi-device/internal/config"
	m "github.com/ardihikaru/go-whatsapp-multi-device/internal/middleware"
	deviceSvc "github.com/ardihikaru/go-whatsapp-multi-device/internal/service/device"
	sessionSvc "github.com/ardihikaru/go-whatsapp-multi-device/internal/service/session"
	"github.com/ardihikaru/go-whatsapp-multi-device/internal/storage"
)

// SessionMainHandler handles all authentication related routes
func SessionMainHandler(cfg *config.Config, db *storage.DataStoreMongo, log *logger.Logger,
	whatsAppBot *botHook.WaManager, httpClient *http.Client) http.Handler {
	r := chi.NewRouter()

	// Initialize services
	deviceService := deviceSvc.NewService(db, log)
	sessionService := sessionSvc.NewService(deviceService, log, whatsAppBot, httpClient, cfg.WhatsappWebhook,
		cfg.WhatsappQrCodeDir, cfg.WhatsappWebhookEcho, cfg.WhatsappWebhookEnabled)

	r.Route("/", func(r chi.Router) {
		r.Route("/{phone}", func(r chi.Router) {
			// extracts the phone on the URL parameter
			r.Use(m.MiddlewarePhoneCtx)

			r.Get("/", sessionRegister(sessionService, log)) // POST /auth/register - register a new WhatsApp account
		})
	})

	return r
}

// sessionRegister processes the request to create new whatsapp session
func sessionRegister(sessionService *sessionSvc.Service, log *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// extracts phone from the context and cast them into a string
		var phoneKey m.Phone = m.PhoneKey
		phone := r.Context().Value(phoneKey).(string)

		// creates new whatsapp session
		err := sessionService.New(r.Context(), phone)
		if err != nil {
			log.Debug(httputils.ResponseText("", httputils.UpdateDataFailed), zap.Error(err))
			httputils.RenderErrResponse(w, r,
				httputils.ResponseText("", httputils.UpdateDataFailed),
				httputils.UpdateDataFailed,
				http.StatusBadRequest, err)
			return
		}

		// prepares response body
		respBody := httputils.Response{
			Success:     true,
			Data:        nil,
			MessageText: "action request has been successfully executed",
			Total:       1,
		}

		// renders OK response
		_ = httputils.RenderOKResponse(w, r, respBody)
	}
}
