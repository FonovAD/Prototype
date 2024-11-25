package api

import (
	"net/http"

	"github.com/FonovAD/Prototype/internal/logger"
)

type server struct {
	router *http.ServeMux
	logger *logger.Logger
}

func newServer(LogLevel string) *server {
	s := &server{
		router: http.NewServeMux(),
		logger: logger.NewLogger(LogLevel),
	}
	s.ConfigureRouter()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) ConfigureRouter() {}
