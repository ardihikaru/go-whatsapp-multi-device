package app

import (
	"net/http"
)

// BuildHttpClient builds a http client
func BuildHttpClient() *http.Client {
	return &http.Client{}
}
