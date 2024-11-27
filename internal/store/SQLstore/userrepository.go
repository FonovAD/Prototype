package sqlstore

import (
	"context"

	"github.com/FonovAD/Prototype/internal/models"
)

type UserRepository struct {
	store *Store
}

func (u *UserRepository) Create(ctx context.Context) (models.User, error) {
	return models.User{}, nil
}
func (u *UserRepository) GetByUID(ctx context.Context, uid int) (models.User, error) {
	return models.User{}, nil
}
func (u *UserRepository) GetByRole(ctx context.Context, role string) ([]models.User, error) {
	return []models.User{}, nil
}
func (u *UserRepository) CheckByToken(ctx context.Context, token string) (bool, error) {
	return false, nil
}
func (u *UserRepository) Delete(ctx context.Context, uid int) error {
	return nil
}
