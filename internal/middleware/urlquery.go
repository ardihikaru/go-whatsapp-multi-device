package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/satumedishub/sea-cucumber-api-service/pkg/utils/httputils"
	"github.com/satumedishub/sea-cucumber-api-service/pkg/utils/query"
)

// URLQueryCtx enriches the request with the captured id on the URL query parameters
func URLQueryCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// initializes with the default values
		limit := DefaultLimitQuery
		offset := DefaultOffsetQuery
		ctx := r.Context()

		// extracts limit query parameters
		limitQuery, err := strconv.Atoi(r.URL.Query().Get(QueryLimitKey))
		if err == nil {
			// limit query is not provided, set the default value
			limit = int64(limitQuery)
		}

		// extracts offset query parameters
		offsetQuery, err := strconv.Atoi(r.URL.Query().Get(QueryOffsetKey))
		if err == nil {
			// limit query is not provided, set the default value
			offset = int64(offsetQuery)
		}

		// extracts order query parameters
		order := r.URL.Query().Get(QueryOrderKey)
		if order != "" {
			// var orderParsed string
			err = json.Unmarshal([]byte(order), &order)
			if err != nil || !query.GetOrderMap()[order] {
				httputils.RenderErrResponse(w, r,
					httputils.ResponseText("", httputils.InvalidOrderQuery),
					httputils.InvalidOrderQuery,
					http.StatusBadRequest, err)
				return
			}
		} else {
			order = DefaultOrderQuery
		}

		// extracts sort query parameters
		sort := r.URL.Query().Get(QuerySortKey)
		if sort != "" {
			err = json.Unmarshal([]byte(sort), &sort)
			if err != nil {
				httputils.RenderErrResponse(w, r,
					httputils.ResponseText("", httputils.InvalidSortQuery),
					httputils.InvalidSortQuery,
					http.StatusBadRequest, err)
				return
			}
		}

		// extracts filter query parameters
		// data type for filter will be mapped into a particular struct when consumed by a specific API
		//  since it may has more than one data type, e.g., a string, or an array of string
		filter := r.URL.Query().Get(QueryFilterKey)

		// defines the URL query parameters
		var qLimitKey QueryLimit = QueryLimitKey
		var qOffsetKey QueryOffset = QueryOffsetKey
		var qOrderKey QueryOrder = QueryOrderKey
		var qSortKey QuerySort = QuerySortKey
		var qFilterKey QueryFilter = QueryFilterKey

		// read the URL parameter
		ctx = context.WithValue(ctx, qLimitKey, limit)
		ctx = context.WithValue(ctx, qOffsetKey, offset)
		ctx = context.WithValue(ctx, qOrderKey, order)
		ctx = context.WithValue(ctx, qSortKey, sort)
		ctx = context.WithValue(ctx, qFilterKey, filter)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
