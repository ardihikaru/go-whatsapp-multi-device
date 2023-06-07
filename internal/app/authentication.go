package app

import (
	"github.com/go-chi/jwtauth/v5"

	"github.com/satumedishub/sea-cucumber-api-service/internal/config"
	"github.com/satumedishub/sea-cucumber-api-service/internal/logger"

	"github.com/satumedishub/sea-cucumber-api-service/pkg/authenticator"
)

// GetTokenAuthentication creates an autehntication token from the authenticator
func GetTokenAuthentication(cfg *config.Config, log *logger.Logger) *jwtauth.JWTAuth {
	tokenAuth, err := authenticator.MakeTokenAuth(cfg.JWTAlgorithm, cfg.JWTSecret)
	if err != nil {
		FatalOnError(err, "failed to create a JWT authenticator")
	}

	return tokenAuth
}
