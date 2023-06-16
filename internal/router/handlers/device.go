package handlers

import (
	"context"
	m "github.com/ardihikaru/go-whatsapp-multi-device/internal/middleware"
	"net/http"

	"github.com/ardihikaru/go-modules/pkg/logger"
	"github.com/ardihikaru/go-modules/pkg/utils/httputils"
	"github.com/go-chi/chi"
	"go.uber.org/zap"

	deviceSvc "github.com/ardihikaru/go-whatsapp-multi-device/internal/service/device"
	svc "github.com/ardihikaru/go-whatsapp-multi-device/internal/service/device"
	"github.com/ardihikaru/go-whatsapp-multi-device/internal/storage"
)

// AuthMainHandler handles all device related routes
func AuthMainHandler(db *storage.DataStoreMongo, log *logger.Logger) http.Handler {
	r := chi.NewRouter()

	// Initialize services
	deviceService := svc.NewService(db, log)

	r.Route("/", func(r chi.Router) {

		r.Post("/", devicePost(deviceService, log)) // POST /api/device - register a new WhatsApp account

		r.Route("/name/{id}", func(r chi.Router) {
			// extracts the id on the URL parameter
			r.Use(m.MiddlewareIDCtx)

			r.Put("/", deviceNamePut(deviceService, log))
		})

		r.Route("/webhook/{id}", func(r chi.Router) {
			// extracts the id on the URL parameter
			r.Use(m.MiddlewareIDCtx)

			r.Put("/", deviceWebhook(deviceService, log))
		})
	})

	return r
}

// devicePost processes the request to create new device data
func devicePost(svc *deviceSvc.Service, log *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqPayload deviceSvc.RegisterPayload

		// extracts request body
		eCode, httpCode, err := httputils.GetJsonBody(r.Body, &reqPayload)
		if err != nil {
			log.Debug(httputils.ResponseText("", eCode), zap.Error(err))
			httputils.RenderErrResponse(w, r, httputils.ResponseText("", eCode), int64(eCode), httpCode, err)
		}

		// submits new device data
		device, err := svc.Register(r.Context(), reqPayload)
		if err != nil {
			log.Debug(httputils.ResponseText("", httputils.CreateDataFailed), zap.Error(err))
			httputils.RenderErrResponse(w, r,
				httputils.ResponseText("", httputils.CreateDataFailed),
				httputils.CreateDataFailed,
				http.StatusBadRequest, err)
			return
		}

		// prepares response body
		respBody := httputils.Response{
			Success:     true,
			Data:        device,
			MessageText: "new device has been registered",
			Total:       1,
		}

		// renders OK response
		_ = httputils.RenderOKResponse(w, r, respBody)
	}
}

// deviceNamePut processes the request to update device name
func deviceNamePut(svc *deviceSvc.Service, log *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqPayload deviceSvc.RegisterPayload

		// extracts userID from the context and cast them into a string
		var idKey m.ID = m.IDKey
		deviceId := r.Context().Value(idKey).(string)

		// extracts request body
		eCode, httpCode, err := httputils.GetJsonBody(r.Body, &reqPayload)
		if err != nil {
			log.Debug(httputils.ResponseText("", eCode), zap.Error(err))
			httputils.RenderErrResponse(w, r, httputils.ResponseText("", eCode), int64(eCode), httpCode, err)
		}

		// update device name now
		err = svc.UpdateDeviceName(context.Background(), deviceId, reqPayload.Name)
		if err != nil {
			log.Warn("failed to update device name information")
			return
		}

		// prepares response body
		respBody := httputils.Response{
			Success:     true,
			Data:        nil,
			MessageText: "device name has been updated",
			Total:       1,
		}

		// renders OK response
		_ = httputils.RenderOKResponse(w, r, respBody)
	}
}

// deviceWebhook processes the request to update webhook URL
func deviceWebhook(svc *deviceSvc.Service, log *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqPayload deviceSvc.RegisterPayload

		// extracts userID from the context and cast them into a string
		var idKey m.ID = m.IDKey
		deviceId := r.Context().Value(idKey).(string)

		// extracts request body
		eCode, httpCode, err := httputils.GetJsonBody(r.Body, &reqPayload)
		if err != nil {
			log.Debug(httputils.ResponseText("", eCode), zap.Error(err))
			httputils.RenderErrResponse(w, r, httputils.ResponseText("", eCode), int64(eCode), httpCode, err)
		}

		// update webhook URL now
		err = svc.UpdateWebhook(context.Background(), deviceId, reqPayload.Webhook)
		if err != nil {
			log.Warn("failed to update device name information")
			return
		}

		// prepares response body
		respBody := httputils.Response{
			Success:     true,
			Data:        nil,
			MessageText: "webhook URL has been updated",
			Total:       1,
		}

		// renders OK response
		_ = httputils.RenderOKResponse(w, r, respBody)
	}
}
