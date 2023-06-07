package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"net/http"

	"github.com/ardihikaru/go-modules/pkg/logger"
	"github.com/ardihikaru/go-modules/pkg/utils/httputils"
	"github.com/go-chi/chi"

	deviceSvc "github.com/ardihikaru/go-whatsapp-multi-device/internal/service/device"
	svc "github.com/ardihikaru/go-whatsapp-multi-device/internal/service/device"
	"github.com/ardihikaru/go-whatsapp-multi-device/internal/storage"
)

// AuthMainHandler handles all authentication related routes
func AuthMainHandler(db *storage.DataStoreMongo, log *logger.Logger) http.Handler {
	r := chi.NewRouter()

	// Initialize services
	deviceService := svc.NewService(db, log)

	r.Route("/", func(r chi.Router) {
		r.Post("/", devicePost(deviceService, log)) // POST /device/register - register a new WhatsApp account
	})

	return r
}

// devicePost processes the request to create new device data
func devicePost(svc *deviceSvc.Service, log *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqPayload deviceSvc.RegisterPayload

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
		err = json.Unmarshal(b, &reqPayload)
		if err != nil {
			log.Debug(httputils.ResponseText("", httputils.InvalidRequestJSON), zap.Error(err))
			httputils.RenderErrResponse(w, r,
				httputils.ResponseText("", httputils.InvalidRequestJSON),
				httputils.InvalidRequestJSON,
				http.StatusBadRequest, err)
			return
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
