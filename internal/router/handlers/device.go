package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ardihikaru/go-modules/pkg/logger"
	"github.com/ardihikaru/go-modules/pkg/utils/httputils"
	"github.com/ardihikaru/go-modules/pkg/utils/query"
	"github.com/go-chi/chi"
	"go.uber.org/zap"

	m "github.com/ardihikaru/go-whatsapp-multi-device/internal/middleware"
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

		r.Post("/", devicePost(deviceService, log))
		r.Get("/", deviceList(deviceService, log))

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

		r.Route("/{id}", func(r chi.Router) {
			// extracts the phone as id on the URL parameter
			r.Use(m.MiddlewareIDCtx)

			r.Get("/", getDeviceByPhone(deviceService, log))
		})
	})

	return r
}

// getDeviceByPhone processes the request to verify if the designated phone on whatsapp or not
func getDeviceByPhone(svc *deviceSvc.Service, log *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// extracts phone from the context and cast them into a string
		var idKey m.ID = m.IDKey
		phone := r.Context().Value(idKey).(string)

		// gets device document
		device, err := svc.GetDeviceByPhone(r.Context(), phone)
		if err != nil {
			httputils.RenderErrResponse(w, r,
				"device not found",
				httputils.BadRequest,
				http.StatusNoContent, nil)
			return
		}

		// prepares response body
		respBody := httputils.Response{
			Success:     true,
			Data:        device,
			MessageText: "fetch success",
			Total:       1,
		}

		// renders OK response
		_ = httputils.RenderOKResponse(w, r, respBody)
	}
}

// deviceList processes the request to list all devices
func deviceList(svc *deviceSvc.Service, log *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		// extracts filter from the context and cast them into a string
		var filterKey m.QueryFilter = m.QueryFilterKey
		filter := r.Context().Value(filterKey).(string)

		var filterParams query.FilterQueryParams

		err = json.Unmarshal([]byte(filter), &filterParams)
		if err != nil {
			httputils.RenderErrResponse(w, r,
				httputils.ResponseText("", httputils.RequestJSONExtractionFailed),
				httputils.RequestJSONExtractionFailed,
				http.StatusBadRequest, err)
			return
		}

		// extracts limit from the context and cast them into a string
		var limitKey m.QueryLimit = m.QueryLimitKey
		limit := r.Context().Value(limitKey).(int64)

		// extracts offset from the context and cast them into a string
		var offsetKey m.QueryOffset = m.QueryOffsetKey
		offset := r.Context().Value(offsetKey).(int64)

		// extracts order from the context and cast them into a string
		var orderKey m.QueryOrder = m.QueryOrderKey
		order := r.Context().Value(orderKey).(string)

		// extracts sort from the context and cast them into a string
		var sortKey m.QuerySort = m.QuerySortKey
		sort := r.Context().Value(sortKey).(string)

		err = json.Unmarshal([]byte(filter), &filterParams)
		if err != nil {
			httputils.RenderErrResponse(w, r,
				httputils.ResponseText("", httputils.RequestJSONExtractionFailed),
				httputils.RequestJSONExtractionFailed,
				http.StatusBadRequest, err)
			return
		}

		// builds query parameters
		params := httputils.GetQueryParams{
			Limit:  limit,
			Offset: offset,
			Order:  order,
			Sort:   sort,
			Search: filterParams.Keyword,
		}

		// list all device data
		total, devices, err := svc.GetDevices(r.Context(), params)
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
			Data:        devices,
			MessageText: "fetch devices success",
			Total:       total,
		}

		// renders OK response
		_ = httputils.RenderOKResponse(w, r, respBody)
	}
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
		err = svc.UpdateWebhook(context.Background(), deviceId, reqPayload.WebhookUrl)
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
