package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
	"go.uber.org/zap"

	"github.com/ardihikaru/go-whatsapp-multi-device/internal/config"
	"github.com/ardihikaru/go-whatsapp-multi-device/internal/logger"
	svc "github.com/ardihikaru/go-whatsapp-multi-device/internal/service/account"
	authSvc "github.com/ardihikaru/go-whatsapp-multi-device/internal/service/auth"
	userSvc "github.com/ardihikaru/go-whatsapp-multi-device/internal/service/user"
	"github.com/ardihikaru/go-whatsapp-multi-device/internal/storage"

	"github.com/ardihikaru/go-whatsapp-multi-device/pkg/utils/httputils"
)

// AuthMainHandler handles all authentication related routes
func AuthMainHandler(cfg *config.Config, db *storage.DataStoreMongo, log *logger.Logger, tokenAuth *jwtauth.JWTAuth) http.Handler {
	r := chi.NewRouter()

	// Initialize services
	userService := userSvc.NewService(db, log)
	accountSvc := svc.NewService(userService, db, log)
	authService := authSvc.NewService(accountSvc, cfg.JWTExpiredInSec, tokenAuth)

	r.Route("/login", func(r chi.Router) {
		r.Post("/", authLogin(authService, log)) // POST /auth/login - authorize login user
	})

	return r
}

// authLogin processes the request to create access token
func authLogin(svc *authSvc.Service, log *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginData authSvc.LoginData
		var errText string

		// extracts request body
		b, err := io.ReadAll(r.Body)
		if err != nil {
			errText = "failed to extract request body"
			log.Debug(errText, zap.Error(err))
			httputils.RenderErrResponse(w, r, errText, 400, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		// read JSON body from the request
		err = json.Unmarshal(b, &loginData)
		if err != nil {
			log.Debug(httputils.ResponseText("", httputils.InvalidRequestJSON), zap.Error(err))
			httputils.RenderErrResponse(w, r,
				httputils.ResponseText("", httputils.InvalidRequestJSON),
				httputils.InvalidRequestJSON,
				http.StatusBadRequest, err)
			return
		}

		// gets JWT payload if login authorized
		payload, err := svc.Authorize(r.Context(), loginData)
		if err != nil {
			log.Debug(httputils.ResponseText("", httputils.LoginFailed), zap.Error(err))
			httputils.RenderErrResponse(w, r,
				httputils.ResponseText("", httputils.LoginFailed),
				httputils.LoginFailed,
				http.StatusBadRequest, err)
			return
		}

		// prepares response body
		respBody := httputils.Response{
			Success:     true,
			Data:        payload,
			MessageText: "login success",
		}

		// renders OK response
		_ = httputils.RenderOKResponse(w, r, respBody)
	}
}
