package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/exp/slices"

	"github.com/satumedishub/sea-cucumber-api-service/internal/logger"
	m "github.com/satumedishub/sea-cucumber-api-service/internal/middleware"
	svc "github.com/satumedishub/sea-cucumber-api-service/internal/service/user"
	uRole "github.com/satumedishub/sea-cucumber-api-service/internal/service/user/role"
	"github.com/satumedishub/sea-cucumber-api-service/internal/storage"

	"github.com/satumedishub/sea-cucumber-api-service/pkg/utils/httputils"
	"github.com/satumedishub/sea-cucumber-api-service/pkg/utils/query"
)

// RoleMainHandler handles all role related routes
func RoleMainHandler(db *storage.DataStoreMongo, log *logger.Logger, tokenAuth *jwtauth.JWTAuth) http.Handler {
	r := chi.NewRouter()

	// initializes user service
	userSvc := svc.NewService(db, log)

	// initializes session middleware resource
	sessionMiddleware := m.UserResource{
		Svc: userSvc,
		Log: log,
	}

	r.Route("/", func(r chi.Router) {
		// Seeks, verifies and validates JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// validates token. Got invalids if (expired, missing)
		r.Use(jwtauth.Authenticator)

		// extracts the query parameters (maybe empty)
		r.Use(m.URLQueryCtx)

		// extracts the id on the URL parameter
		r.Use(sessionMiddleware.SessionCtx)

		r.HandleFunc("/", listRole(userSvc, log)) // GET /users - Read a list of users.
	})

	return r
}

// listRole fetches valid roles
func listRole(userSvc *svc.Service, log *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		// extracts userID from the context and cast them into a string
		var filterKey m.QueryFilter = m.QueryFilterKey
		filter := r.Context().Value(filterKey).(string)

		var filterParams query.FilterListParams

		err = json.Unmarshal([]byte(filter), &filterParams)
		if err != nil {
			panic(err)
		}

		roles := []uRole.Role{}
		for _, role := range uRole.GetRoleList() {
			if len(filterParams.Ids) > 0 && !slices.Contains(filterParams.Ids, role) {
				continue
			}

			roles = append(roles, uRole.Role{
				RoleId:   role,
				RoleName: role,
			})
		}

		// prepares response body
		respBody := httputils.Response{
			Data:        roles,
			MessageText: "fetching roles success",
			Total:       int64(len(roles)),
		}

		// renders OK response
		_ = httputils.RenderOKResponse(w, r, respBody)
	}
}
