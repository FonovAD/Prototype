package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/FonovAD/Prototype/internal/models"
)

func (s *server) HandleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			s.logger.Info(r.Method, r.RemoteAddr, " Unexpected HTTP Method")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		s.logger.Info(r.Method, r.RemoteAddr)
		w.WriteHeader(http.StatusOK)
	}
}

func (s *server) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	s.logger.Error(r.Method, r.URL.Path, err.Error())
	s.metricMonitor.IncErrorCount(r.Method, r.URL.Path, http.StatusInternalServerError)
	return
}

func (s *server) CreateUser() http.HandlerFunc {
	type response struct {
		Token string `json:"token"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			s.logger.Info(r.Method, r.RemoteAddr, "Unexpected HTTP Method")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		auth := strings.Split(r.Header.Get("Authorization"), " ")
		if len(auth) < 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := auth[1]
		fmt.Println(token)
		user, err := s.store.User().GetByToken(r.Context(), token)
		if err != nil {
			s.ServerError(w, r, err)
			return
		}
		fmt.Println(user)
		if user == nil || user.Role != models.ROLE_ADMIN {
			w.WriteHeader(http.StatusForbidden)
			s.logger.Info(r.Method, r.URL.Path, http.StatusForbidden)
			s.metricMonitor.IncErrorCount(r.Method, r.URL.Path, http.StatusForbidden)
			return
		}
		NewUser, err := s.store.User().Create(r.Context())
		if err != nil {
			s.ServerError(w, r, err)
			return
		}
		resp := response{
			Token: NewUser.Token,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			s.ServerError(w, r, err)
			return
		}
		s.logger.Info(r.Method, r.URL.Path, http.StatusOK)
		w.WriteHeader(http.StatusOK)
	}
}

func (s *server) Link() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Query().Get("path")
		if path == "" {
			s.logger.Info(r.Method, r.URL.Path, "Path is empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		linkModel, err := s.store.Link().GetByShortLink(r.Context(), path)
		if err != nil {
			s.ServerError(w, r, err)
		}
		s.logger.Info(r.Method, r.URL.Path, http.StatusOK)
		http.Redirect(w, r, linkModel.OriginLink, http.StatusFound)
	}
}
