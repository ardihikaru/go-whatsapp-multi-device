package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
	"go.uber.org/zap"

	"github.com/satumedishub/sea-cucumber-api-service/internal/logger"
	m "github.com/satumedishub/sea-cucumber-api-service/internal/middleware"
	svc "github.com/satumedishub/sea-cucumber-api-service/internal/service/account"
	userSvc "github.com/satumedishub/sea-cucumber-api-service/internal/service/user"
	"github.com/satumedishub/sea-cucumber-api-service/internal/storage"

	"github.com/satumedishub/sea-cucumber-api-service/pkg/utils/httputils"
	"github.com/satumedishub/sea-cucumber-api-service/pkg/utils/query"
)

// AccountMainHandler handles all account related routes
func AccountMainHandler(db *storage.DataStoreMongo, log *logger.Logger, tokenAuth *jwtauth.JWTAuth) http.Handler {
	r := chi.NewRouter()

	// initializes user service
	userService := userSvc.NewService(db, log)
	accountSvc := svc.NewService(userService, db, log)

	// initializes session middleware resource
	sessionMiddleware := m.UserResource{
		Svc: userService,
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

		r.Route("/{id}", func(r chi.Router) {
			// Seeks, verifies and validates JWT tokens
			r.Use(jwtauth.Verifier(tokenAuth))

			// validates token. Got invalids if (expired, missing)
			r.Use(jwtauth.Authenticator)

			// extracts the id on the URL parameter
			r.Use(m.MiddlewareIDCtx)

			// extracts the id on the URL parameter
			r.Use(sessionMiddleware.SessionCtx)

			r.Get("/", getAccountById(accountSvc, log)) // GET /accounts/{id} - Read a single account by :id.
			r.Put("/", putAccount(accountSvc, log))     // PUT /accounts - Edit an existing account
		})

		r.HandleFunc("/", listAccount(accountSvc, log)) // GET /accounts - Read a list of accounts.
		r.Post("/", postAccount(accountSvc, log))       // POST /accounts - Create a new account
	})

	return r
}

// getAccountById processes the request to collect and response account data by {id}
func getAccountById(accountSvc *svc.Service, log *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		// extracts userID from the context and cast them into a string
		var idKey m.ID = m.IDKey
		userId := r.Context().Value(idKey).(string)

		// gets user document
		var user svc.Account
		user, err = accountSvc.GetAccountByID(r.Context(), userId, true)
		if err != nil {
			log.Warn("user not found")
			httputils.RenderErrResponse(w, r,
				"user not found",
				httputils.FailedToFetchData,
				http.StatusBadRequest, err)
			return
		}

		// prepares response body
		respBody := httputils.Response{
			Data:        user,
			MessageText: "user found",
		}

		// renders OK response
		_ = httputils.RenderOKResponse(w, r, respBody)
	}
}

// listAccount processes the request to collect and response all accounts
// the optional query parameters may filter out the results
func listAccount(accountSvc *svc.Service, log *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		// extracts filter from the context and cast them into a string
		var filterKey m.QueryFilter = m.QueryFilterKey
		filter := r.Context().Value(filterKey).(string)

		var filterParams query.FilterQueryParams

		err = json.Unmarshal([]byte(filter), &filterParams)
		if err != nil {
			httputils.RenderErrResponse(w, r,
				httputils.ResponseText("", httputils.RequestJSONExtractionFailed),
				httputils.RequestJSONExtractionFailed,
				http.StatusBadRequest, err)
			return
		}

		// extracts limit from the context and cast them into a string
		var limitKey m.QueryLimit = m.QueryLimitKey
		limit := r.Context().Value(limitKey).(int64)

		// extracts offset from the context and cast them into a string
		var offsetKey m.QueryOffset = m.QueryOffsetKey
		offset := r.Context().Value(offsetKey).(int64)

		// extracts order from the context and cast them into a string
		var orderKey m.QueryOrder = m.QueryOrderKey
		order := r.Context().Value(orderKey).(string)

		// extracts sort from the context and cast them into a string
		var sortKey m.QuerySort = m.QuerySortKey
		sort := r.Context().Value(sortKey).(string)

		// var data map[string]interface{}
		err = json.Unmarshal([]byte(filter), &filterParams)
		if err != nil {
			httputils.RenderErrResponse(w, r,
				httputils.ResponseText("", httputils.RequestJSONExtractionFailed),
				httputils.RequestJSONExtractionFailed,
				http.StatusBadRequest, err)
			return
		}

		// builds query parameters
		params := httputils.GetQueryParams{
			Limit:  limit,
			Offset: offset,
			Order:  order,
			Sort:   sort,
			Search: filterParams.Keyword,
		}

		total, users, err := accountSvc.GetAccounts(r.Context(), params, true)
		if err != nil {
			log.Debug(httputils.ResponseText("", httputils.FailedToFetchData), zap.Error(err))
			httputils.RenderErrResponse(w, r,
				httputils.ResponseText("", httputils.FailedToFetchData),
				httputils.FailedToFetchData,
				http.StatusBadRequest, err)
			return
		}

		// prepares response body
		respBody := httputils.Response{
			Data:        users,
			MessageText: "account data fetched successfully",
			Total:       total,
		}

		// renders OK response
		_ = httputils.RenderOKResponse(w, r, respBody)
	}
}

// postAccount processes the request to create a new user
func postAccount(accountSvc *svc.Service, log *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var userData svc.PostData

		// extracts request body
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Debug(httputils.ResponseText("", httputils.InvalidRequestJSON), zap.Error(err))
			httputils.RenderErrResponse(w, r,
				httputils.ResponseText("", httputils.InvalidRequestJSON),
				httputils.InvalidRequestJSON,
				http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		// read JSON body from the request
		err = json.Unmarshal(b, &userData)
		if err != nil {
			log.Debug(httputils.ResponseText("", httputils.RequestJSONExtractionFailed), zap.Error(err))
			httputils.RenderErrResponse(w, r,
				httputils.ResponseText("", httputils.RequestJSONExtractionFailed),
				httputils.RequestJSONExtractionFailed,
				http.StatusBadRequest, err)
			return
		}

		// executes post process
		user, err := accountSvc.Post(r.Context(), userData)
		if err != nil {
			log.Debug(httputils.ResponseText("", httputils.InsertFailed), zap.Error(err))
			httputils.RenderErrResponse(w, r,
				"",
				httputils.InsertFailed,
				http.StatusBadRequest, err)
			return
		}

		// prepares response body
		respBody := httputils.Response{
			Data:        user,
			MessageText: "insert new account success",
			Total:       1,
		}

		// renders OK response
		_ = httputils.RenderOKResponse(w, r, respBody)
	}
}

// putAccount processes the request to update an existing account
func putAccount(accountSvc *svc.Service, log *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var userData svc.PutData

		// extracts userID from the context and cast them into a string
		var idKey m.ID = m.IDKey
		userId := r.Context().Value(idKey).(string)

		// extracts request body
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Debug(httputils.ResponseText("", httputils.InvalidRequestJSON), zap.Error(err))
			httputils.RenderErrResponse(w, r,
				httputils.ResponseText("", httputils.InvalidRequestJSON),
				httputils.InvalidRequestJSON,
				http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		// read JSON body from the request
		err = json.Unmarshal(b, &userData)
		if err != nil {
			log.Debug(httputils.ResponseText("", httputils.RequestJSONExtractionFailed), zap.Error(err))
			httputils.RenderErrResponse(w, r,
				httputils.ResponseText("", httputils.RequestJSONExtractionFailed),
				httputils.RequestJSONExtractionFailed,
				http.StatusBadRequest, err)
			return
		}

		// executes put process
		user, err := accountSvc.Put(r.Context(), userData, userId)
		if err != nil {
			log.Debug(httputils.ResponseText("", httputils.UpdateFailed), zap.Error(err))
			httputils.RenderErrResponse(w, r,
				err.Error(),
				httputils.UpdateFailed,
				http.StatusBadRequest, err)
			return
		}

		// prepares response body
		respBody := httputils.Response{
			Data:        user,
			MessageText: "update password success",
		}

		// renders OK response
		_ = httputils.RenderOKResponse(w, r, respBody)
	}
}
