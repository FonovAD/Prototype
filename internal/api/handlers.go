package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/FonovAD/Prototype/internal/models"
	sqlstore "github.com/FonovAD/Prototype/internal/store/SQLstore"
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
		user, err := s.store.User().GetByToken(r.Context(), token)
		if err != nil {
			s.ServerError(w, r, err)
			return
		}
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

func (s *server) CreateLink() http.HandlerFunc {
	type request struct {
		Link string `json:"origin_link"`
	}
	type response struct {
		ShortLink string `json:"short_link"`
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

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.logger.Info(r.Method, r.RemoteAddr, http.StatusUnprocessableEntity, err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		token := auth[1]
		user, err := s.store.User().GetByToken(r.Context(), token)
		if err != nil {
			s.ServerError(w, r, err)
			return
		}
		if user == nil {
			w.WriteHeader(http.StatusForbidden)
			s.logger.Info(r.Method, r.URL.Path, http.StatusForbidden)
			s.metricMonitor.IncErrorCount(r.Method, r.URL.Path, http.StatusForbidden)
			return
		}
		NewLink, err := s.store.Link().Create(r.Context(), user.UID, req.Link, "")
		if err == sqlstore.InvalidLinkError {
			w.WriteHeader(http.StatusUnprocessableEntity)
			s.logger.Info(r.Method, r.URL.Path, http.StatusUnprocessableEntity)
			return
		} else if err != nil {
			s.ServerError(w, r, err)
			return
		}
		resp := response{
			ShortLink: fmt.Sprintf("http://%s/short/%s", s.url, NewLink.ShortLink),
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

		path := r.PathValue("path")
		if path == "" {
			s.logger.Info(r.Method, r.URL.Path, "Path is empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		linkModel, err := s.store.Link().GetByShortLink(r.Context(), path)
		if err != nil {
			s.ServerError(w, r, err)
			return
		}
		s.logger.Info(r.Method, r.URL.Path, http.StatusOK)
		http.Redirect(w, r, linkModel.OriginLink, http.StatusFound)
	}
}

func (s *server) OutputHtml() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tmpl, err := template.ParseFiles("./internal/api/UI/firstUi.html")
		if err != nil {
			s.ServerError(w, r, err)
			return
		}

		err = tmpl.Execute(w, r)
		if err != nil {
			s.ServerError(w, r, err)
		}

		s.logger.Info(r.Method, r.URL.Path, http.StatusOK)
	}
}
