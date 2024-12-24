package api

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FonovAD/Prototype/internal/logger"
	"github.com/FonovAD/Prototype/internal/metric"
	"github.com/FonovAD/Prototype/internal/store"
	sqlstore "github.com/FonovAD/Prototype/internal/store/SQLstore"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type server struct {
	router        *http.ServeMux
	logger        *logger.Logger
	metricMonitor metric.Monitor
	store         store.Store
	url           string
}

func NewServer(logger *logger.Logger, metricMonitor metric.Monitor, store store.Store, URL string) *server {
	s := &server{
		router:        http.NewServeMux(),
		logger:        logger,
		metricMonitor: metricMonitor,
		store:         store,
		url:           URL,
	}
	s.ConfigureRouter()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) ConfigureRouter() {
	s.router.HandleFunc("/hello", s.HandleHello())
	s.router.HandleFunc("/create_user", s.CreateUser())
	s.router.HandleFunc("/create_link", s.CreateLink())
	s.router.HandleFunc("/{path}", s.Link())

	s.router.Handle("/metrics", promhttp.Handler())
}

type PostgresConfig struct {
	DBUser string
	DBPass string
	DBAddr string
	DBPort string
	DBName string
}

func Start(logLevel, serverAddr, URL string, postgresParam PostgresConfig) error {
	ctxb := context.Background()
	serv := NewServer(logger.New(logLevel), metric.New(), InitPostgres(
		ctxb,
		postgresParam.DBUser,
		postgresParam.DBPass,
		postgresParam.DBAddr,
		postgresParam.DBPort,
		postgresParam.DBName,
	), URL)
	http.Handle("/metrics", promhttp.Handler())
	servWithMiddleware := serv.WriteMetric(serv)
	return http.ListenAndServe(serverAddr, servWithMiddleware)
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

func InitPostgres(ctxb context.Context, dbUser, dbPass, dbAddr, dbPort, dbName string) store.Store {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbAddr, dbPort, dbName)
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	ctx, cancel := context.WithTimeout(ctxb, 1*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping to database: %v", err)
	}
	return sqlstore.New(db, 1*time.Second)
}
