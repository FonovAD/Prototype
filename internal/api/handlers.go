package api

import "net/http"

func (s *server) HandleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			s.logger.Info(r.Method, r.RemoteAddr, "Unexpected HTTP Method")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		s.logger.Info(r.Method, r.RemoteAddr)
		w.WriteHeader(http.StatusOK)
	}
}
