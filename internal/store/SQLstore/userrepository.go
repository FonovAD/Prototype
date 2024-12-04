package sqlstore

import (
	"context"
	"database/sql"

	"github.com/FonovAD/Prototype/internal/models"
)

type UserRepository struct {
	store *Store
}

func (u *UserRepository) Create(ctx context.Context) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.store.QueryTimeout)
	defer cancel()
	user := &models.User{
		Token: models.CreateToken(models.MAX_TOKEN_LENGHT),
		Role:  models.ROLE_USER,
	}
	if err := u.store.db.QueryRowContext(ctx,
		"INSERT INTO users(Token, Role) VALUES ($1, $2) RETURNING UID;",
		user.Token,
		user.Role,
	).Scan(
		&user.UID,
	); err != nil {
		return &models.User{}, err
	}
	return user, nil
}

func (u *UserRepository) GetByUID(ctx context.Context, uid int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.store.QueryTimeout)
	defer cancel()
	user := &models.User{
		UID: uid,
	}
	if err := u.store.db.QueryRowContext(ctx,
		"SELECT Token, Role FROM users WHERE uid = $1",
		user.UID,
	).Scan(
		&user.Token,
		&user.Role,
	); err != nil {
		return &models.User{}, err
	}
	return user, nil
}

func (u *UserRepository) GetByToken(ctx context.Context, token string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.store.QueryTimeout)
	defer cancel()
	user := &models.User{
		Token: token,
	}
	err := u.store.db.QueryRowContext(ctx,
		"SELECT UID, Role FROM users WHERE Token = $1",
		user.Token,
	).Scan(
		&user.UID,
		&user.Role,
	)
	if err == sql.ErrNoRows {
		return &models.User{}, nil
	} else if err != nil {
		return &models.User{}, err
	}
	return user, nil
}

func (u *UserRepository) GetByRole(ctx context.Context, role string, limit int, offset int) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.store.QueryTimeout)
	defer cancel()
	var users []models.User = make([]models.User, limit)
	rows, err := u.store.db.QueryContext(ctx,
		"SELECT uid, Token, Role FROM users WHERE role = $1 LIMIT $2 OFFSET $3",
		role, limit, offset,
	)
	if err != nil {
		return []models.User{}, err
	}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.UID,
			&user.Token,
			&user.Role,
		); err != nil {
			return []models.User{}, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *UserRepository) CheckByToken(ctx context.Context, token string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, u.store.QueryTimeout)
	defer cancel()
	var expectedToken string
	if err := u.store.db.QueryRowContext(ctx,
		"SELECT Token FROM users WHERE token = $1",
		token,
	).Scan(&expectedToken); err != nil {
		return false, err
	}
	if token == expectedToken {
		return true, nil
	}
	return false, nil
}

func (u *UserRepository) Delete(ctx context.Context, unique interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, u.store.QueryTimeout)
	defer cancel()
	switch unique.(type) {
	case int:
		if _, err := u.store.db.ExecContext(ctx,
			"DELETE FROM users WHERE uid = $1",
			unique,
		); err != nil {
			return err
		}
	case string:
		if _, err := u.store.db.ExecContext(ctx,
			"DELETE FROM users WHERE Token = $1",
			unique,
		); err != nil {
			return err
		}
	}
	return nil
}
