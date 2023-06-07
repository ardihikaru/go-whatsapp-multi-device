package app

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"

	"github.com/satumedishub/sea-cucumber-api-service/internal/config"
	"github.com/satumedishub/sea-cucumber-api-service/internal/logger"
	"github.com/satumedishub/sea-cucumber-api-service/internal/storage"
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
