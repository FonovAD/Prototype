package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/FonovAD/Prototype/internal/logger"
	"github.com/FonovAD/Prototype/internal/metric"
	"github.com/FonovAD/Prototype/internal/models"
	"github.com/FonovAD/Prototype/internal/store"
	sqlstore "github.com/FonovAD/Prototype/internal/store/SQLstore"
	"github.com/stretchr/testify/assert"
)

func TestServer_Hello(t *testing.T) {
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
		httpMethod   string
	}{
		{
			name: "Base case",
			payload: map[string]interface{}{
				"message": "hello!",
			},
			expectedCode: http.StatusOK,
			httpMethod:   http.MethodGet,
		},
		{
			name: "Unexpected http method",
			payload: map[string]interface{}{
				"message": "hello!",
			},
			expectedCode: http.StatusMethodNotAllowed,
			httpMethod:   http.MethodPost,
		},
	}
	db, f := sqlstore.SetupTestDB(t, "test")
	defer f()
	s := NewServer(logger.New("debug"), metric.NewTest(), sqlstore.New(db, 5*time.Second), "127.0.0.1:80")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			body, err := json.Marshal(tc.payload)
			if err != nil {
				assert.NoError(t, err)
				return
			}
			req, _ := http.NewRequest(tc.httpMethod, "/hello", bytes.NewReader(body))
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_CreateUser(t *testing.T) {
	testCases := []struct {
		name         string
		expectedCode int
		httpMethod   string
		prepare      func(context.Context, store.Store) string
	}{
		{
			name:         "Base case",
			expectedCode: http.StatusOK,
			httpMethod:   http.MethodPost,
			prepare: func(ctx context.Context, s store.Store) string {
				return "test"
			},
		},
		{
			name:         "Unexpected http method",
			expectedCode: http.StatusMethodNotAllowed,
			httpMethod:   http.MethodGet,
			prepare: func(ctx context.Context, s store.Store) string {
				return "test"
			},
		},
		{
			name:         "Not Admin",
			expectedCode: http.StatusForbidden,
			httpMethod:   http.MethodPost,
			prepare: func(ctx context.Context, s store.Store) string {
				u, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}
				return u.Token
			},
		},
	}
	db, f := sqlstore.SetupTestDB(t, "test")
	defer f()
	s := NewServer(logger.New("debug"), metric.NewTest(), sqlstore.New(db, 5*time.Second), "127.0.0.1:80")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			ctxb := context.Background()
			token := tc.prepare(ctxb, s.store)

			req, _ := http.NewRequest(tc.httpMethod, "/user/create", nil)
			req.Header.Set("Authorization", "token "+token)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_CreateLink(t *testing.T) {
	testCases := []struct {
		name          string
		expectedCode  int
		httpMethod    string
		payload       interface{}
		preferredLink string
		prepare       func(context.Context, store.Store) *models.User
	}{
		{
			name:         "Base case",
			expectedCode: http.StatusOK,
			httpMethod:   http.MethodPost,
			payload: map[string]interface{}{
				"origin_link": "http://validLink.ru",
			},
			prepare: func(ctx context.Context, s store.Store) *models.User {
				u, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}
				return u
			},
		},
		{
			name:         "Unexpected http method",
			expectedCode: http.StatusMethodNotAllowed,
			httpMethod:   http.MethodGet,
			payload: map[string]interface{}{
				"origin_link": "http://validLink.ru",
			},
			prepare: func(ctx context.Context, s store.Store) *models.User {
				u, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}
				return u
			},
		},
		{
			name:         "InvalidLink",
			expectedCode: http.StatusUnprocessableEntity,
			httpMethod:   http.MethodPost,
			payload: map[string]interface{}{
				"origin_link": "invalidLink",
			},
			prepare: func(ctx context.Context, s store.Store) *models.User {
				u, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}
				return u
			},
		},
	}
	db, f := sqlstore.SetupTestDB(t, "test")
	defer f()
	s := NewServer(logger.New("debug"), metric.NewTest(), sqlstore.New(db, 5*time.Second), "127.0.0.1:80")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			ctxb := context.Background()
			u := tc.prepare(ctxb, s.store)

			body, err := json.Marshal(tc.payload)
			if err != nil {
				assert.NoError(t, err)
				return
			}
			req, _ := http.NewRequest(tc.httpMethod, "/link/create", bytes.NewReader(body))
			req.Header.Set("Authorization", "token "+u.Token)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_DeliteUserByToken(t *testing.T) {
	testCases := []struct {
		name         string
		expectedCode int
		httpMethod   string
		prepare      func(context.Context, store.Store) (string, string)
	}{
		{
			name:         "Admin removes User",
			expectedCode: http.StatusOK,
			httpMethod:   http.MethodPost,
			prepare: func(ctx context.Context, s store.Store) (string, string) {

				u2, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}
				return "test", u2.Token
			},
		},
		{
			name:         "User removes Admin",
			expectedCode: http.StatusForbidden,
			httpMethod:   http.MethodPost,
			prepare: func(ctx context.Context, s store.Store) (string, string) {

				u1, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}
				return u1.Token, "test"
			},
		},
		{
			name:         "Admin removes Admin",
			expectedCode: http.StatusForbidden,
			httpMethod:   http.MethodPost,
			prepare: func(ctx context.Context, s store.Store) (string, string) {
				return "test", "test"
			},
		},
		{
			name:         "User1 removes User1",
			expectedCode: http.StatusOK,
			httpMethod:   http.MethodPost,
			prepare: func(ctx context.Context, s store.Store) (string, string) {
				u1, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}

				return u1.Token, u1.Token
			},
		},
		{
			name:         "User1 removes User2",
			expectedCode: http.StatusForbidden,
			httpMethod:   http.MethodPost,
			prepare: func(ctx context.Context, s store.Store) (string, string) {
				u1, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}
				u2, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}
				return u1.Token, u2.Token
			},
		},
		{
			name:         "Unexpected http method",
			expectedCode: http.StatusMethodNotAllowed,
			httpMethod:   http.MethodGet,
			prepare: func(ctx context.Context, s store.Store) (string, string) {
				u1, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}
				u2, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}
				return u1.Token, u2.Token
			},
		},
	}
	db, f := sqlstore.SetupTestDB(t, "test")
	defer f()
	s := NewServer(logger.New("debug"), metric.NewTest(), sqlstore.New(db, 5*time.Second), "127.0.0.1:80")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			ctxb := context.Background()
			uToken1, uToken2 := tc.prepare(ctxb, s.store)
			body, err := json.Marshal(map[string]string{"user_token": uToken2})
			if err != nil {
				assert.NoError(t, err)
				return
			}
			req, _ := http.NewRequest(tc.httpMethod, "/user/delete/token", bytes.NewReader(body))
			req.Header.Set("Authorization", "token "+uToken1)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_DeliteUserByUID(t *testing.T) {
	testCases := []struct {
		name         string
		expectedCode int
		httpMethod   string
		prepare      func(context.Context, store.Store) (string, int)
	}{
		{
			name:         "Admin removes User",
			expectedCode: http.StatusOK,
			httpMethod:   http.MethodPost,
			prepare: func(ctx context.Context, s store.Store) (string, int) {
				u1, err := s.User().GetByToken(ctx, "test")
				if err != nil {
					t.Fatal(err)
				}
				u2, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}
				return u1.Token, u2.UID
			},
		},
		{
			name:         "User removes Admin",
			expectedCode: http.StatusForbidden,
			httpMethod:   http.MethodPost,
			prepare: func(ctx context.Context, s store.Store) (string, int) {
				u1, err := s.User().GetByToken(ctx, "test")
				if err != nil {
					t.Fatal(err)
				}
				u2, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}
				return u2.Token, u1.UID
			},
		},
		{
			name:         "Admin removes Admin",
			expectedCode: http.StatusForbidden,
			httpMethod:   http.MethodPost,
			prepare: func(ctx context.Context, s store.Store) (string, int) {
				u1, err := s.User().GetByToken(ctx, "test")
				if err != nil {
					t.Fatal(err)
				}
				return "test", u1.UID
			},
		},
		{
			name:         "User1 removes User1",
			expectedCode: http.StatusOK,
			httpMethod:   http.MethodPost,
			prepare: func(ctx context.Context, s store.Store) (string, int) {
				u1, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}

				return u1.Token, u1.UID
			},
		},
		{
			name:         "User1 removes User2",
			expectedCode: http.StatusForbidden,
			httpMethod:   http.MethodPost,
			prepare: func(ctx context.Context, s store.Store) (string, int) {
				u1, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}
				u2, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}
				return u1.Token, u2.UID
			},
		},
		{
			name:         "Unexpected http method",
			expectedCode: http.StatusMethodNotAllowed,
			httpMethod:   http.MethodGet,
			prepare: func(ctx context.Context, s store.Store) (string, int) {
				u1, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}
				u2, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}
				return u1.Token, u2.UID
			},
		},
	}
	db, f := sqlstore.SetupTestDB(t, "test")
	defer f()
	s := NewServer(logger.New("debug"), metric.NewTest(), sqlstore.New(db, 5*time.Second), "127.0.0.1:80")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			ctxb := context.Background()
			uToken1, uID2 := tc.prepare(ctxb, s.store)
			body, err := json.Marshal(map[string]int{"user_ID": uID2})
			if err != nil {
				assert.NoError(t, err)
				return
			}
			req, _ := http.NewRequest(tc.httpMethod, "/user/delete/uid", bytes.NewReader(body))
			req.Header.Set("Authorization", "token "+uToken1)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_DeleteLink(t *testing.T) {
	testCases := []struct {
		name          string
		expectedCode  int
		httpMethod    string
		payload       interface{}
		preferredLink string
		prepare       func(context.Context, store.Store) string
	}{
		{
			name:         "User deletes own link",
			expectedCode: http.StatusOK,
			httpMethod:   http.MethodPost,
			payload: map[string]interface{}{
				"origin_link": "http://validLink.ru",
			},
			prepare: func(ctx context.Context, s store.Store) string {
				u, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}
				_, err = s.Link().Create(ctx, u.UID, "http://validLink.ru", "")
				if err != nil {
					t.Fatal(err)
				}
				return u.Token
			},
		},
		{
			name:         "Admin deletes User link",
			expectedCode: http.StatusOK,
			httpMethod:   http.MethodPost,
			payload: map[string]interface{}{
				"origin_link": "http://validLink.ru",
			},
			prepare: func(ctx context.Context, s store.Store) string {
				u, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}
				_, err = s.Link().Create(ctx, u.UID, "http://validLink.ru", "")

				if err != nil {
					t.Fatal(err)
				}
				return "test"
			},
		},
		{
			name:         "User1 deletes User2 link",
			expectedCode: http.StatusForbidden,
			httpMethod:   http.MethodPost,
			payload: map[string]interface{}{
				"origin_link": "http://validLink.ru",
			},
			prepare: func(ctx context.Context, s store.Store) string {
				u1, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}
				u2, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}
				_, err = s.Link().Create(ctx, u2.UID, "http://validLink.ru", "")
				if err != nil {
					t.Fatal(err)
				}
				return u1.Token
			},
		},
		{
			name:         "User deletes Admin link",
			expectedCode: http.StatusForbidden,
			httpMethod:   http.MethodPost,
			payload: map[string]interface{}{
				"origin_link": "http://validLink.ru",
			},
			prepare: func(ctx context.Context, s store.Store) string {
				u, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}
				adm, err := s.User().GetByToken(ctx, "test")
				if err != nil {
					t.Fatal(err)
				}
				_, err = s.Link().Create(ctx, adm.UID, "http://validLink.ru", "")
				if err != nil {
					t.Fatal(err)
				}
				return u.Token
			},
		},
		{
			name:         "InvalidLink",
			expectedCode: http.StatusUnprocessableEntity,
			httpMethod:   http.MethodPost,
			payload: map[string]interface{}{
				"origin_link": "invalidLink",
			},
			prepare: func(ctx context.Context, s store.Store) string {
				u, err := s.User().Create(ctx)
				if err != nil {
					t.Fatal(err)
				}
				return u.Token
			},
		},
	}
	db, f := sqlstore.SetupTestDB(t, "test")
	defer f()
	s := NewServer(logger.New("debug"), metric.NewTest(), sqlstore.New(db, 5*time.Second), "127.0.0.1:80")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			ctxb := context.Background()
			token := tc.prepare(ctxb, s.store)

			body, err := json.Marshal(tc.payload)
			if err != nil {
				assert.NoError(t, err)
				return
			}
			req, _ := http.NewRequest(tc.httpMethod, "/link/delete", bytes.NewReader(body))
			req.Header.Set("Authorization", "token "+token)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
