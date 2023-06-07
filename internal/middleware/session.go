package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth/v5"

	"github.com/ardihikaru/go-whatsapp-multi-device/internal/storage"

	auth "github.com/ardihikaru/go-whatsapp-multi-device/pkg/authenticator"
)

// SessionCtx enriches the request with the captured JWT private claims
func (rs UserResource) SessionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sessionKey JWTSession = SessionKey

		// extracts token from the header
		token, _, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		// token is authenticated, extracts the private claims
		privateClaims := token.PrivateClaims()

		// extracts
		session := auth.Session{
			AccountId: privateClaims[string(storage.FnAccountsId)].(string),
			UserId:    privateClaims[string(storage.FnAccountsUserId)].(string),
			Username:  privateClaims[string(storage.FnAccountsUsername)].(string),
		}

		// token is authenticated, enrich token to the request parameter
		ctx := context.WithValue(r.Context(), sessionKey, session)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
