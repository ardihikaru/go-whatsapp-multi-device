package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
)

// MiddlewarePhoneCtx enriches the request with the captured phone on the URL parameter
func MiddlewarePhoneCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// define the URL parameters
		var phoneKey Phone = PhoneKey

		// read the URL parameter
		ctx := context.WithValue(r.Context(), phoneKey, chi.URLParam(r, PhoneKey))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
