package sqlstore

import (
	"database/sql"

	"github.com/FonovAD/Prototype/internal/store"
)

type Store struct {
	db       *sql.DB
	userRepo *UserRepository
	linkRepo *LinkRepository
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
