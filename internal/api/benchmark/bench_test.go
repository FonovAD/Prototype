package benchmark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FonovAD/Prototype/internal/api"
)

func BenchmarkServer_Hello(b *testing.B) {
	s := api.NewServer("test")
	payload := map[string]interface{}{
		"message": "hello!",
	}
	rec := httptest.NewRecorder()
	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest(http.MethodGet, "/hello", bytes.NewReader(body))
		s.ServeHTTP(rec, req)
	}
}
