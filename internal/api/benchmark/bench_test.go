package benchmark

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime/trace"
	"testing"
	"time"

	"github.com/FonovAD/Prototype/config"
	"github.com/FonovAD/Prototype/internal/api"
	"github.com/FonovAD/Prototype/internal/logger"
	"github.com/FonovAD/Prototype/internal/metric"
	"github.com/FonovAD/Prototype/internal/store"
	sqlstore "github.com/FonovAD/Prototype/internal/store/SQLstore"
)

func BenchmarkServer_Hello(b *testing.B) {
	configPath := "./Benchmark_default.yaml"
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		b.Fatalf("Failed to read config file: %v", err)
	}

	db, f := SetupDB(cfg.SQLite3.Path, cfg.SQLite3.Schema)
	defer f()
	s := api.NewServer(logger.New(cfg.API.LogLevel), metric.New(), db, cfg.API.URL)

	payload := map[string]interface{}{
		"message": "hello!",
	}
	rec := httptest.NewRecorder()
	body, err := json.Marshal(payload)
	if err != nil {
		b.Fatalf(err.Error())
	}

	b.ResetTimer()
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

	configPath := "./Benchmark_default.yaml"
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		b.Fatalf("Failed to read config file: %v", err)
	}

	db, clearDB := SetupDB(cfg.SQLite3.Path, cfg.SQLite3.Schema)
	defer clearDB()
	s := api.NewServer(logger.New(cfg.API.LogLevel), metric.New(), db, cfg.API.URL)

	if err := trace.Start(f); err != nil {
		b.Fatalf("Error starting trace: %v", err)
	}
	defer trace.Stop()

	rec := httptest.NewRecorder()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest(http.MethodPost, "/create_user", nil)
		req.Header.Set("Authorization", "Token test")
		s.ServeHTTP(rec, req)
	}
}

func SetupDB(databasePath, schemaPath string) (store.Store, func()) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_synchronous=NORMAL", databasePath))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.SetMaxOpenConns(80)
	db.SetMaxIdleConns(20)
	schemaSQL, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		log.Fatalf("Failed read db schema: %v", err)
	}
	if _, err := db.Exec(fmt.Sprintf(string(schemaSQL), "test")); err != nil {
		log.Fatalf("Failed to setup test database schema: %v", err)
	}
	return sqlstore.New(db, 1*time.Second), func() {
		db.Close()
		os.Remove(databasePath)
	}
}
