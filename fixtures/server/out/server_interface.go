package http

import "net/http"

var _ ServerInterface = (*Server)(nil)

type ServerInterface interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}
