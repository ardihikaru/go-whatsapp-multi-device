package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
)

// MiddlewareIDCtx enriches the request with the captured id on the URL parameter
func MiddlewareIDCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// define the URL parameters
		var idKey ID = IDKey

		// read the URL parameter
		ctx := context.WithValue(r.Context(), idKey, chi.URLParam(r, IDKey))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
