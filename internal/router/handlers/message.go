package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ardihikaru/go-modules/pkg/logger"
	"github.com/ardihikaru/go-modules/pkg/utils/httputils"
	botHook "github.com/ardihikaru/go-modules/pkg/whatsappbot/wawebhook"
	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/ardihikaru/go-whatsapp-multi-device/internal/config"
	deviceSvc "github.com/ardihikaru/go-whatsapp-multi-device/internal/service/device"
	sessionSvc "github.com/ardihikaru/go-whatsapp-multi-device/internal/service/session"
	"github.com/ardihikaru/go-whatsapp-multi-device/internal/storage"
)

// MessageMainHandler handles all whatsapp message related routes
func MessageMainHandler(cfg *config.Config, db *storage.DataStoreMongo, log *logger.Logger,
	whatsAppBot *botHook.WaManager, httpClient *http.Client, bcList *botHook.BotClientList) http.Handler {
	r := chi.NewRouter()

	// initializes services
	deviceService := deviceSvc.NewService(db, log)
	sessionService := sessionSvc.NewService(deviceService, log, whatsAppBot, httpClient,
		cfg.WhatsappImageDir, cfg.WhatsappQrCodeDir, cfg.WhatsappWebhookEcho, cfg.WhatsappWebhookEnabled,
		cfg.WhatsappQrToTerminal, bcList)

	r.Route("/", func(r chi.Router) {
		r.Post("/text", postMessage(sessionService, log))
		r.Post("/image", postImageMessage(sessionService, log))
	})

	return r
}

// postMessage processes the request to send a whatsapp message
func postMessage(sessionService *sessionSvc.Service, log *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload botHook.MessagePayload

		// extracts request body
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Debug(httputils.ResponseText("", httputils.InvalidRequestJSON), zap.Error(err))
			httputils.RenderErrResponse(w, r, httputils.ResponseText("", httputils.InvalidRequestJSON),
				httputils.InvalidRequestJSON, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		// read JSON body from the request
		err = json.Unmarshal(b, &payload)
		if err != nil {
			log.Debug(httputils.ResponseText("", httputils.InvalidRequestJSON), zap.Error(err))
			httputils.RenderErrResponse(w, r,
				httputils.ResponseText("", httputils.InvalidRequestJSON),
				httputils.InvalidRequestJSON,
				http.StatusBadRequest, err)
			return
		}

		// submits new message
		err = sessionService.SendTextMessage(payload)
		if err != nil {
			log.Debug(httputils.ResponseText("", httputils.CreateDataFailed), zap.Error(err))
			httputils.RenderErrResponse(w, r,
				err.Error(),
				httputils.CreateDataFailed,
				http.StatusBadRequest, nil)
			return
		}

		// prepares response body
		respBody := httputils.Response{
			Success:     true,
			Data:        nil,
			MessageText: "message has been sent",
			Total:       1,
		}

		// renders OK response
		_ = httputils.RenderOKResponse(w, r, respBody)
	}
}

// postImageMessage processes the request to send a whatsapp image-based message
func postImageMessage(sessionService *sessionSvc.Service, log *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload botHook.MessagePayload

		// extracts request body
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Debug(httputils.ResponseText("", httputils.InvalidRequestJSON), zap.Error(err))
			httputils.RenderErrResponse(w, r, httputils.ResponseText("", httputils.InvalidRequestJSON),
				httputils.InvalidRequestJSON, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		// read JSON body from the request
		err = json.Unmarshal(b, &payload)
		if err != nil {
			log.Debug(httputils.ResponseText("", httputils.InvalidRequestJSON), zap.Error(err))
			httputils.RenderErrResponse(w, r,
				httputils.ResponseText("", httputils.InvalidRequestJSON),
				httputils.InvalidRequestJSON,
				http.StatusBadRequest, err)
			return
		}

		// submits new message
		err = sessionService.SendImageMessage(payload)
		if err != nil {
			log.Debug(httputils.ResponseText("", httputils.CreateDataFailed), zap.Error(err))
			httputils.RenderErrResponse(w, r,
				err.Error(),
				httputils.CreateDataFailed,
				http.StatusBadRequest, nil)
			return
		}

		// prepares response body
		respBody := httputils.Response{
			Success:     true,
			Data:        nil,
			MessageText: "image message has been sent",
			Total:       1,
		}

		// renders OK response
		_ = httputils.RenderOKResponse(w, r, respBody)
	}
}
