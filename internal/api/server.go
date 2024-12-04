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
)

type server struct {
	router        *http.ServeMux
	logger        *logger.Logger
	metricMonitor metric.Monitor
	store         store.Store
}

func NewServer(logger *logger.Logger, metricMonitor metric.Monitor, store store.Store) *server {
	s := &server{
		router:        http.NewServeMux(),
		logger:        logger,
		metricMonitor: metricMonitor,
		store:         store,
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
}

func Start(logLevel, serverAddr string) error {
	serv := NewServer(logger.New(logLevel), metric.New(), SetupDB())
	return http.ListenAndServe(serverAddr, serv)
}

func SetupDB() store.Store {
	databasePath := "./LinkShortener"
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	schema := `
PRAGMA foreign_keys = ON;
CREATE TABLE IF NOT EXISTS users(
UID INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
Token TEXT NOT NULL,
Role varchar(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS links(
UID INTEGER REFERENCES users(UID) ON DELETE CASCADE,
OriginLink TEXT UNIQUE NOT NULL,
ShortLink TEXT UNIQUE NOT NULL,
CreatedAt integer,
ExpirationTime integer NOT NULL,
Status varchar(10) NOT NULL,
ScheduledDeletionTime integer NOT NULL
);`
	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("Failed to setup test database schema: %v", err)
	}
	return sqlstore.New(db, 1*time.Second)
}
