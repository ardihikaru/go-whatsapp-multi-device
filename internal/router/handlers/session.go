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

// SessionMainHandler handles all session related routes
func SessionMainHandler(cfg *config.Config, db *storage.DataStoreMongo, log *logger.Logger,
	whatsAppBot *botHook.WaManager, httpClient *http.Client, bcList *botHook.BotClientList) http.Handler {
	r := chi.NewRouter()

	// initializes services
	deviceService := deviceSvc.NewService(db, log)
	sessionService := sessionSvc.NewService(deviceService, log, whatsAppBot, httpClient,
		cfg.WhatsappImageDir, cfg.WhatsappQrCodeDir, cfg.WhatsappWebhookEcho, cfg.WhatsappWebhookEnabled,
		cfg.WhatsappQrToTerminal, bcList)

	// initializes middleware resources
	waM := m.Resource{
		Log:        log,
		BotClients: bcList,
	}

	r.Route("/", func(r chi.Router) {

		r.Route("/{phone}", func(r chi.Router) {
			// extracts the phone on the URL parameter
			r.Use(waM.WhatsappCtx)

			r.Get("/", sessionConnect(sessionService, log)) // GET /api/session/{phone} - register a new session
		})

		r.Route("/phone/{phone}", func(r chi.Router) {
			// extracts the phone on the URL parameter
			r.Use(m.PhoneMiddlewareCtx)

			r.Delete("/", sessionDisconnect(sessionService, log)) // DELETE /api/session/{phone} - delete session
		})

		r.Route("/on-whatsapp/{phone}", func(r chi.Router) {
			// extracts the phone on the URL parameter
			r.Use(m.PhoneMiddlewareCtx)

			r.Get("/", isOnWhatsapp(sessionService, log)) // GET /api/session/on-whatsapp/{phone} - is on WA?
		})
	})

	return r
}

// sessionConnect processes the request to create new whatsapp session
func sessionConnect(sessionService *sessionSvc.Service, log *logger.Logger) func(http.ResponseWriter, *http.Request) {
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

// sessionDisconnect processes the request to delete an existing whatsapp session
func sessionDisconnect(sessionService *sessionSvc.Service, log *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// extracts phone from the context and cast them into a string
		var phoneKey m.Phone = m.PhoneKey
		phone := r.Context().Value(phoneKey).(string)

		// if key exists, disconnect and remove the key first
		msg := sessionService.Disconnect(phone)

		// prepares response body
		respBody := httputils.Response{
			Success:     true,
			Data:        nil,
			MessageText: msg,
			Total:       1,
		}

		// renders OK response
		_ = httputils.RenderOKResponse(w, r, respBody)
	}
}

// isOnWhatsapp processes the request to verify if the designated phone on whatsapp or not
func isOnWhatsapp(sessionService *sessionSvc.Service, log *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// extracts phone from the context and cast them into a string
		var phoneKey m.Phone = m.PhoneKey
		phone := r.Context().Value(phoneKey).(string)

		// checks if on WA or not
		onWhatsapp, err := sessionService.IsOnWhatsapp(phone)
		if err != nil {
			httputils.RenderErrResponse(w, r,
				"failed to check on the Whatsapp Server",
				httputils.BadRequest,
				http.StatusBadRequest, nil)
			return
		}
		// special case: no session active that can be utilized
		if onWhatsapp == nil {
			httputils.RenderErrResponse(w, r,
				"no active session to use",
				httputils.FailedToFetchData,
				http.StatusNoContent, nil)
			return
		}

		// prepares response body
		respBody := httputils.Response{
			Success:     true,
			Data:        *onWhatsapp,
			MessageText: "fetch success",
			Total:       1,
		}

		// renders OK response
		_ = httputils.RenderOKResponse(w, r, respBody)
	}
}
