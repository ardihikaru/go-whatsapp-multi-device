// Package httputils provides utilities for HTTP related operations
package httputils

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

// the constants below describes the available response codes and error codes in this system.
// the first digit (x0000) defines the category of the error
// the second digit (0x000) defines the type of code, where 0 for response code and 0 for error code
// the last three digit (00xxx) defines the incremental value of the same category and type of code
const (
	// RenderFailed is an application error code for a failed rendering.
	// a failed rendering may happen when the fetched response body is invalid
	RenderFailed = 11001

	// InvalidRequestJSON is an application error code where system is unable to extract JSON on the request body
	InvalidRequestJSON = 11002

	// RequestJSONExtractionFailed is an application error code where system is unable to unmarshal captured JSON data
	RequestJSONExtractionFailed = 11003

	// InputValidationError is an application error code where the captured input data got validation error
	InputValidationError = 11004

	// UnauthorizedAccess is an application error code where the identity has no authorize to access
	UnauthorizedAccess = 11005

	// LoginFailed is an application error code where user failed to log in
	LoginFailed = 11006

	// InsertFailed is an application error code where user failed to insert a new record
	InsertFailed = 11007

	// FailedToFetchData is an application error code where user failed to fetch data from the database
	FailedToFetchData = 11008

	// MapToJSONConversionFailed is an application error code where system failed to convert Map to JSON format
	MapToJSONConversionFailed = 11009

	// MissingPhoneInHeader is an application error code where the identity is missing in the header
	MissingPhoneInHeader = 11010

	// InvalidOrderQuery is an application error code where the order is invalid
	InvalidOrderQuery = 11011

	// InvalidSortQuery is an application error code where the sort is invalid
	InvalidSortQuery = 11012

	// UpdateFailed is an application error code where user failed to update an existing record
	UpdateFailed = 11013
)

// responseText is list of error test for each application-level error code
var responseText = map[int]string{
	RenderFailed: "failed to render a valid response body",

	InvalidRequestJSON:          "failed to extract request body",
	RequestJSONExtractionFailed: "failed to read JSON body from the request",
	InputValidationError:        "got input validation error",
	UnauthorizedAccess:          "identity is unauthorized to access this API",
	LoginFailed:                 "invalid login data",
	InsertFailed:                "failed to insert a new record",
	UpdateFailed:                "failed to update an existing record",
	FailedToFetchData:           "failed to fetch data from the database",
	MapToJSONConversionFailed:   "failed to convert Map to JSON format",
	MissingPhoneInHeader:        "phone is missing from the header",
	InvalidOrderQuery:           "invalid order in URL parameters",
	InvalidSortQuery:            "invalid sort in URL parameters",
}

// ResponseText returns a text for the HTTP status code in the application level.
// It returns the empty string if the code is unknown.
func ResponseText(identifier string, code int) string {
	if identifier != "" {
		return fmt.Sprintf("[%s] %s", identifier, responseText[code])
	} else {
		return responseText[code]
	}
}

// Response renderer type for handling all sorts of http response
type Response struct {
	Err            error  `json:"-"`                                                                            // low-level runtime error
	HTTPStatusCode int    `json:"-"`                                                                            // response response status code
	Data           any    `json:"data"`                                                                         // always set as empty
	Total          int64  `json:"total"`                                                                        // total fetched records
	Success        bool   `json:"success,omitempty"`                                                            // success status
	MessageText    string `json:"message,omitempty" example:"Resource not found."`                              // user-level status message
	AppErrCode     int64  `json:"code,omitempty" example:"404"`                                                 // application-specific error code
	ErrorText      string `json:"error,omitempty" example:"The requested resource was not found on the server"` // application-level error message, for debugging
} // @name  Response

// Render implements the github.com/go-chi/render.Renderer interface for Response
func (e *Response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// httpErrPayload returns a structured error response
func httpErrPayload(errText string, appErrCode int64, httpStatusCode int, err error) render.Renderer {
	errorTxt := errText
	if err != nil {
		errorTxt = err.Error()
	}
	return &Response{
		HTTPStatusCode: httpStatusCode,
		AppErrCode:     appErrCode,
		Err:            err,
		ErrorText:      errorTxt,
		MessageText:    errText,
	}
}

// httpOKPayload returns a structured OK response
func httpOKPayload(respBody Response) render.Renderer {
	return &Response{
		HTTPStatusCode: http.StatusOK,
		Success:        true,
		Data:           respBody.Data,
		Total:          respBody.Total,
		MessageText:    respBody.MessageText,
	}
}

// RenderErrResponse renders the error http response
func RenderErrResponse(w http.ResponseWriter, r *http.Request, errText string, appErrCode int64,
	httpStatusCode int, err error) {
	_ = render.Render(w, r, httpErrPayload(errText, appErrCode, httpStatusCode, err))
}

// RenderOKResponse returns a rendered http response
// rendering may fails and returns an error, otherwise it returns nil value
func RenderOKResponse(w http.ResponseWriter, r *http.Request, respBody Response) error {
	return render.Render(w, r, httpOKPayload(respBody))
}
