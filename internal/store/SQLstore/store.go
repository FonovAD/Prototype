package sqlstore

import (
	"database/sql"
	"time"

	"github.com/FonovAD/Prototype/internal/store"
)

const Schema = `
PRAGMA foreign_keys = ON;
CREATE TABLE IF NOT EXISTS users(
UID INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
Token TEXT NOT NULL,
Role varchar(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS links(
UID INTEGER REFERENCES users(UID) ON DELETE CASCADE,
OriginLink TEXT NOT NULL,
ShortLink TEXT UNIQUE NOT NULL,
CreatedAt integer,
ExpirationTime integer NOT NULL,
Status varchar(10) NOT NULL,
ScheduledDeletionTime integer NOT NULL,
PRIMARY KEY (UID, OriginLink)
);

INSERT INTO users(Token, Role) VALUES("%s", "admin");
`

type Store struct {
	db           *sql.DB
	userRepo     *UserRepository
	linkRepo     *LinkRepository
	QueryTimeout time.Duration
}

func New(db *sql.DB, QueryTimeout time.Duration) *Store {
	return &Store{
		db:           db,
		QueryTimeout: QueryTimeout,
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepo != nil {
		return s.userRepo
	}

	s.userRepo = &UserRepository{
		store: s,
	}

	return s.userRepo
}

func (s *Store) Link() store.LinkRepository {
	if s.linkRepo != nil {
		return s.linkRepo
	}

	s.linkRepo = &LinkRepository{
		store: s,
	}

	return s.linkRepo
}
