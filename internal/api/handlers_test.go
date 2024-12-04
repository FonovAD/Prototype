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
	s := NewServer(logger.New("debug"), metric.NewTest(), sqlstore.New(db, 5*time.Second))
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
	s := NewServer(logger.New("debug"), metric.NewTest(), sqlstore.New(db, 5*time.Second))
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			ctxb := context.Background()
			token := tc.prepare(ctxb, s.store)

			req, _ := http.NewRequest(tc.httpMethod, "/create_user", nil)
			req.Header.Set("Authorization", "token "+token)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
