package http

import (
	"log"
	"net/http"
)

type Server struct {
	logger *log.Logger
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
}
