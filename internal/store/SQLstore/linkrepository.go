package sqlstore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/FonovAD/Prototype/internal/models"
)

type LinkRepository struct {
	store       *Store
	ExpDuration time.Duration
	DelDuration time.Duration
}

var (
	InvalidLinkError = errors.New("The passed link is not valid")
	ExistLinkError   = errors.New("The desired link already exists")
)

// Заменить вызов двух функций на транзакцию
func (l *LinkRepository) Create(ctx context.Context, UID int, originLink string, preferredLink string) (*models.Link, error) {
	ctx, cancel := context.WithTimeout(ctx, l.store.QueryTimeout)
	defer cancel()
	if models.IsValidURL(originLink) {
		if preferredLink == "" {
			preferredLink = models.GenerateShortLink(10)
		}
		exist, err := l.ShortLinkExist(ctx, preferredLink)
		if err != nil {
			return nil, err
		}
		if !exist {
			newLink := &models.Link{
				UID:                   UID,
				OriginLink:            originLink,
				ShortLink:             preferredLink,
				CreateTime:            time.Now().Unix(),
				ExpireTime:            time.Now().Add(l.ExpDuration).Unix(),
				Status:                models.STATUS_ACTIVE,
				ScheduledDeletionTime: time.Now().Add(l.DelDuration).Unix(),
			}
			CreatedLink := &models.Link{}
			err := l.store.db.QueryRowContext(ctx,
				`INSERT INTO links (UID, OriginLink, ShortLink, CreatedAt, ExpirationTime, Status, ScheduledDeletionTime)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				ON CONFLICT(UID, OriginLink) DO UPDATE 
				SET 
    				ShortLink = EXCLUDED.ShortLink,
    				CreatedAt = EXCLUDED.CreatedAt,
    				ExpirationTime = EXCLUDED.ExpirationTime,
    				Status = EXCLUDED.Status,
    				ScheduledDeletionTime = EXCLUDED.ScheduledDeletionTime
				RETURNING UID, OriginLink, ShortLink, CreatedAt, ExpirationTime, Status, ScheduledDeletionTime;`,
				newLink.UID,
				newLink.OriginLink,
				newLink.ShortLink,
				newLink.CreateTime,
				newLink.ExpireTime,
				newLink.Status,
				newLink.ScheduledDeletionTime,
			).Scan(
				&CreatedLink.UID,
				&CreatedLink.OriginLink,
				&CreatedLink.ShortLink,
				&CreatedLink.CreateTime,
				&CreatedLink.ExpireTime,
				&CreatedLink.Status,
				&CreatedLink.ScheduledDeletionTime,
			)
			if err != nil {
				return nil, fmt.Errorf("error inserting record: %w", err)
			}
			return CreatedLink, nil
		} else {
			return nil, ExistLinkError
		}
	} else {
		return nil, InvalidLinkError
	}
}

func (l *LinkRepository) GetByUID(ctx context.Context, uid int) ([]*models.Link, error) {
	ctx, cancel := context.WithTimeout(ctx, l.store.QueryTimeout)
	defer cancel()
	var links []*models.Link
	rows, err := l.store.db.QueryContext(ctx,
		"SELECT uid, OriginLink, ShortLink, CreatedAt, ExpirationTime, Status, ScheduledDeletionTime FROM links WHERE uid = $1;",
		uid,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		link := &models.Link{}
		fmt.Println("ROW!")
		if err := rows.Scan(
			&link.UID,
			&link.OriginLink,
			&link.ShortLink,
			&link.CreateTime,
			&link.ExpireTime,
			&link.Status,
			&link.ScheduledDeletionTime,
		); err != nil {
			return []*models.Link{}, err
		}
		links = append(links, link)
	}
	return links, nil
}

func (l *LinkRepository) GetByOriginLink(ctx context.Context, originLink string) (*models.Link, error) {
	ctx, cancel := context.WithTimeout(ctx, l.store.QueryTimeout)
	defer cancel()
	link := &models.Link{}
	if err := l.store.db.QueryRowContext(ctx,
		"SELECT * FROM links WHERE OriginLink = $1;",
		originLink,
	).Scan(
		&link.UID,
		&link.OriginLink,
		&link.ShortLink,
		&link.CreateTime,
		&link.ExpireTime,
		&link.Status,
		&link.ScheduledDeletionTime,
	); err != nil {
		return nil, err
	}
	return link, nil
}

func (l *LinkRepository) GetByShortLink(ctx context.Context, shortLink string) (*models.Link, error) {
	ctx, cancel := context.WithTimeout(ctx, l.store.QueryTimeout)
	defer cancel()
	link := &models.Link{}
	if err := l.store.db.QueryRowContext(ctx,
		"SELECT * FROM links WHERE ShortLink = $1;",
		shortLink,
	).Scan(
		&link.UID,
		&link.OriginLink,
		&link.ShortLink,
		&link.CreateTime,
		&link.ExpireTime,
		&link.Status,
		&link.ScheduledDeletionTime,
	); err != nil {
		return nil, err
	}
	return link, nil
}

func (l *LinkRepository) ShortLinkExist(ctx context.Context, shortLink string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, l.store.QueryTimeout)
	defer cancel()
	var count int
	if err := l.store.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM links WHERE ShortLink = $1;",
		shortLink,
	).Scan(&count); err != nil {
		return true, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (l *LinkRepository) Delete(ctx context.Context, originLink string) error {
	ctx, cancel := context.WithTimeout(ctx, l.store.QueryTimeout)
	defer cancel()
	if err := l.store.db.QueryRowContext(ctx,
		"DELETE FROM links WHERE OriginLink = $1;",
		originLink,
	).Err(); err != nil {
		return err
	}
	return nil
}

func (l *LinkRepository) ReActivate(ctx context.Context, originLink string) error {
	ctx, cancel := context.WithTimeout(ctx, l.store.QueryTimeout)
	defer cancel()
	if err := l.store.db.QueryRowContext(ctx,
		"UPDATE links SET Status = $1 WHERE OriginLink = $2;",
		models.STATUS_ACTIVE,
		originLink,
	).Err(); err != nil {
		return err
	}
	return nil
}
