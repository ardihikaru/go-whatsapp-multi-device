// Package router provides the base configurations to build a router.
// This package is open for an extension whether new routes need to be added.
package router

import (
	"github.com/ardihikaru/go-modules/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	"github.com/ardihikaru/go-whatsapp-multi-device/internal/app"
	h "github.com/ardihikaru/go-whatsapp-multi-device/internal/router/handlers"
)

// GetRouter configures a chi router and starts the http server
// @title          Go WhatsApp Multi-device API Service
// @description    Go WhatsApp Multi-device API Service implements sample RESTApi
// @contact.name   Muhammad Febrian Ardiansyah
// @contact.email  mfardiansyah@outlook.com
// @BasePath       /
func GetRouter(deps *app.Dependencies) *chi.Mux {
	r := chi.NewRouter()

	if deps.Log != nil {
		r.Use(logger.SetLogger(deps.Log))
	}

	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   deps.Config.CORSAllowOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   deps.Config.CORSAllowHeaders,
		ExposedHeaders:   deps.Config.CORSExposedHeaders,
		AllowCredentials: true,
		MaxAge:           600, // Maximum value not ignored by any of major browsers
	}))

	buildTree(r, deps)

	return r
}

// buildTree builds routes
func buildTree(r *chi.Mux, deps *app.Dependencies) {
	// handles device related route(s)
	r.Mount("/api/device", h.AuthMainHandler(deps.DB, deps.Log))

	// handles session related route(s)
	r.Mount("/api/session", h.SessionMainHandler(deps.Config, deps.DB, deps.Log, deps.WhatsAppBot,
		deps.HttpClient, deps.BotClients))

	// handles whatsapp message related route(s)
	r.Mount("/api/message", h.MessageMainHandler(deps.Config, deps.DB, deps.Log, deps.WhatsAppBot,
		deps.HttpClient, deps.BotClients))
}
