package api

import (
	"net/http"

	"github.com/FonovAD/Prototype/internal/logger"
	"github.com/FonovAD/Prototype/internal/metric"
)

type server struct {
	router        *http.ServeMux
	logger        *logger.Logger
	metricMonitor *metric.MetricMonitor
}

func NewServer(logger *logger.Logger, metricMonitor *metric.MetricMonitor) *server {
	s := &server{
		router:        http.NewServeMux(),
		logger:        logger,
		metricMonitor: metricMonitor,
	}
	s.ConfigureRouter()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) ConfigureRouter() {
	s.router.HandleFunc("/hello", s.HandleHello())
}

func API(logLevel, serverAddr string) error {
	serv := NewServer(logger.New(logLevel), metric.New())
	return http.ListenAndServe(serverAddr, serv)
}
