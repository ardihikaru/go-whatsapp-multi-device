package app

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"

	"github.com/ardihikaru/go-whatsapp-multi-device/internal/config"
	"github.com/ardihikaru/go-whatsapp-multi-device/internal/logger"
	"github.com/ardihikaru/go-whatsapp-multi-device/internal/storage"
)

// Dependencies holds the primitives and structs and/or interfaces that are required
// for the application's business logic.
type Dependencies struct {
	Config     *config.Config
	DB         *storage.DataStoreMongo
	Log        *logger.Logger
	TokenAuth  *jwtauth.JWTAuth
	HttpClient *http.Client
}
