package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

	s := NewServer("Debug")
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
