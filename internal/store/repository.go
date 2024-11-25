package store

import (
	"github.com/FonovAD/Proto/internal/models"
)

type UserRepository interface {
	Create() (models.User, error)
	GetByUID(int) (models.User, error)
	GetByRole(string) ([]models.User, error)
	CheckByToken(string) (bool, error)
	Delete(int) error
}

type LinkRepository interface {
	Create() (models.Link, error)
	GetByUID(int) ([]models.Link, error)
	GetByOriginLink(string) (models.Link, error)
	Delete(string) error
	ReActivate(string) error
}
