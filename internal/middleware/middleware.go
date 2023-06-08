package middleware

import (
	"github.com/ardihikaru/go-modules/pkg/logger"
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
	// QueryOrderKey is the order key to store query offset which is captured from the request Query parameters
	QueryOrderKey = "order"
	// QuerySortKey is the sort key to store query offset which is captured from the request Query parameters
	QuerySortKey = "soft"

	// QueryFilterKey is the filter key to store query offset which is captured from the request Query parameters
	QueryFilterKey = "filter"

	// IDKey is the identifier key to store ID which is captured from the request URL parameters
	IDKey = "id"

	// PhoneKey is the identifier key to store phone which is captured from the request URL parameters
	PhoneKey = "phone"
)

// Resource is a middleware resource
type Resource struct {
	Log *logger.Logger
}
