package benchmark

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime/trace"
	"testing"
	"time"

	"github.com/FonovAD/Prototype/internal/api"
	"github.com/FonovAD/Prototype/internal/logger"
	"github.com/FonovAD/Prototype/internal/metric"
	"github.com/FonovAD/Prototype/internal/store"
	sqlstore "github.com/FonovAD/Prototype/internal/store/SQLstore"
)

func BenchmarkServer_Hello(b *testing.B) {
	s := api.NewServer(logger.New("debug"), metric.NewTest(), SetupDB(), "127.0.0.1:80")
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

func BenchmarkServer_CreateUser(b *testing.B) {
	f, err := os.Create("trace_benchmark.out")
	if err != nil {
		b.Fatalf("Error creating trace file: %v", err)
	}
	defer f.Close()

	if err := trace.Start(f); err != nil {
		b.Fatalf("Error starting trace: %v", err)
	}
	defer trace.Stop()

	s := api.NewServer(logger.New("debug"), metric.NewTest(), SetupDB(), "127.0.0.1:80")
	rec := httptest.NewRecorder()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest(http.MethodPost, "/create_user", nil)
		req.Header.Set("Authorization", "token test")
		s.ServeHTTP(rec, req)
	}
}

func SetupDB() store.Store {
	databasePath := "./LinkShortener"
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if _, err := db.Exec(sqlstore.Schema); err != nil {
		log.Fatalf("Failed to setup test database schema: %v", err)
	}
	return sqlstore.New(db, 1*time.Second)
}
