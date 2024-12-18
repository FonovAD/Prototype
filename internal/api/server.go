package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/FonovAD/Prototype/internal/logger"
	"github.com/FonovAD/Prototype/internal/metric"
	"github.com/FonovAD/Prototype/internal/store"
	sqlstore "github.com/FonovAD/Prototype/internal/store/SQLstore"
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

func Start(logLevel, serverAddr, URL string) error {
	serv := NewServer(logger.New(logLevel), metric.New(), SetupDB(), URL)
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
