package middleware

import (
	"context"
	"github.com/ardihikaru/go-modules/pkg/utils/httputils"
	"net/http"

	"github.com/go-chi/chi"
)

// WhatsappCtx validates the request related with whatsapp bot
func (rs Resource) WhatsappCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// define the URL parameters
		var phoneKey Phone = PhoneKey
		phone := chi.URLParam(r, PhoneKey)

		// validates
		if _, ok := (*rs.BotClients)[phone]; ok {
			httputils.RenderErrResponse(w, r,
				"session for this device has been logged in",
				304,
				http.StatusNotModified, nil)
			return
		}

		// read the URL parameter
		ctx := context.WithValue(r.Context(), phoneKey, chi.URLParam(r, PhoneKey))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
