package api

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/FonovAD/Prototype/config"
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
	s.router.HandleFunc("/short/{path}", s.Link())
	s.router.HandleFunc("/", s.OutputHtml())
	s.router.Handle("/metrics", promhttp.Handler())
}

func Start(cfg *config.Config, UseSQLite3 bool) error {
	ctxb := context.Background()
	var db store.Store
	if UseSQLite3 {
		db = SetupDB(
			cfg.SQLite3.Path,
			cfg.SQLite3.Schema,
		)
	} else {
		db = InitPostgres(
			ctxb,
			cfg.Postgres.User,
			cfg.Postgres.Password,
			cfg.Postgres.Host,
			cfg.Postgres.Port,
			cfg.Postgres.Database,
		)
	}
	serv := NewServer(logger.New(cfg.API.LogLevel), metric.New(), db, cfg.API.URL)
	// http.Handle("/metrics", promhttp.Handler())
	servWithMiddleware := serv.WriteMetric(serv)
	log.Print("Server started with param: ", cfg)
	return http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.API.Host, cfg.API.Port), servWithMiddleware)
}

func SetupDB(databasePath, schemaPath string) store.Store {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	schemaSQL, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		log.Fatalf("Couldn't read the database schema: %v", err)
	}
	if _, err := db.Exec(string(schemaSQL)); err != nil {
		log.Fatalf("Failed to setup test database schema: %v", err)
	}
	return sqlstore.New(db, 1*time.Second)
}

func InitPostgres(ctxb context.Context, dbUser, dbPass, dbAddr, dbPort, dbName string) store.Store {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbAddr, dbPort, dbName)
	db, err := sql.Open("postgres", databaseURL)
	db.SetMaxOpenConns(80)
	db.SetMaxIdleConns(20)
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
