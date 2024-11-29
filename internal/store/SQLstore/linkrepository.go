package sqlstore

import (
	"context"

	"github.com/FonovAD/Prototype/internal/models"
)

type LinkRepository struct {
	store *Store
}

func (l *LinkRepository) Create(ctx context.Context) (*models.Link, error) {
	return &models.Link{}, nil
}
func (l *LinkRepository) GetByUID(ctx context.Context, uid int) ([]models.Link, error) {
	return make([]models.Link, 5), nil
}
func (l *LinkRepository) GetByOriginLink(ctx context.Context, originLink string) ([]models.Link, error) {
	return make([]models.Link, 5), nil
}
func (l *LinkRepository) Delete(ctx context.Context, originLink string) error {
	return nil
}
func (l *LinkRepository) ReActivate(ctx context.Context, originLink string) error {
	return nil
}
