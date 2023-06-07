// Package router provides the base configurations to build a router.
// This package is open for an extension whether new routes need to be added.
package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	"github.com/ardihikaru/go-whatsapp-multi-device/internal/app"
	"github.com/ardihikaru/go-whatsapp-multi-device/internal/logger"
	h "github.com/ardihikaru/go-whatsapp-multi-device/internal/router/handlers"
)

// GetRouter configures a chi router and starts the http server
// @title          Template API Service
// @description    Template API Service implements sample RESTApi
// @contact.name   Developer SatuMedis
// @contact.email  dev@satumedis.com
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

func buildTree(r *chi.Mux, deps *app.Dependencies) {
	// handles auth related route(s)
	r.Mount("/auth", h.AuthMainHandler(deps.Config, deps.DB, deps.Log, deps.TokenAuth))

	// handles users route(s)
	r.Mount("/users", h.UserMainHandler(deps.DB, deps.Log, deps.TokenAuth))

	// handles accounts route(s)
	r.Mount("/accounts", h.AccountMainHandler(deps.DB, deps.Log, deps.TokenAuth))

	// handles roles route(s)
	r.Mount("/roles", h.RoleMainHandler(deps.DB, deps.Log, deps.TokenAuth))

}
