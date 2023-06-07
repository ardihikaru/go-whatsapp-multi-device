package middleware

import (
	"github.com/satumedishub/sea-cucumber-api-service/internal/logger"
	svcUser "github.com/satumedishub/sea-cucumber-api-service/internal/service/user"
)

type ID string
type Phone string
type QueryLimit string
type QueryOffset string
type QueryOrder string
type QuerySort string
type QueryFilter string
type JWTSession string

const (
	// DefaultLimitQuery is the default value of the query limit
	DefaultLimitQuery = int64(10)
	// DefaultOffsetQuery is the default value of the query offset
	DefaultOffsetQuery = int64(0)
	// DefaultOrderQuery is the default value of the query offset
	DefaultOrderQuery = "ASC"

	// QueryLimitKey is the limit key to store query limit which is captured from the request Query parameters
	QueryLimitKey = "limit"
	// QueryOffsetKey is the offset key to store query offset which is captured from the request Query parameters
	QueryOffsetKey = "offset"
	// QueryOffsetKey is the order key to store query offset which is captured from the request Query parameters
	QueryOrderKey = "order"
	// QueryOffsetKey is the sort key to store query offset which is captured from the request Query parameters
	QuerySortKey = "soft"

	// QueryOffsetKey is the filter key to store query offset which is captured from the request Query parameters
	QueryFilterKey = "filter"

	// IDKey is the identifier key to store ID which is captured from the request URL parameters
	IDKey = "id"

	// SessionKey is the context key to store JWT private claims which is captured from the request
	SessionKey = "session"
)

// UserResource is a middleware resource for user
type UserResource struct {
	Svc *svcUser.Service
	Log *logger.Logger
}
