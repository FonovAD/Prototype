package store

import (
	"context"

	"github.com/FonovAD/Prototype/internal/models"
)

type UserRepository interface {
	Create(context.Context) (*models.User, error)
	GetByUID(context.Context, int) (*models.User, error)
	GetByRole(context.Context, string, int, int) ([]models.User, error)
	GetByToken(context.Context, string) (*models.User, error)
	CheckByToken(context.Context, string) (bool, error)
	// TODO: add user deletion by any uniq parameter: Users have a unique uid and token.
	Delete(context.Context, interface{}) error
}

type LinkRepository interface {
	Create(context.Context, int, string, string) (*models.Link, error)
	GetByUID(context.Context, int) ([]*models.Link, error)
	GetByOriginLink(context.Context, string) (*models.Link, error)
	ShortLinkExist(context.Context, string) (bool, error)
	Delete(context.Context, string) error
	ReActivate(context.Context, string) error
}
