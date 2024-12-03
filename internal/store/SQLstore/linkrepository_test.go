package sqlstore_test

import (
	"context"
	"testing"
	"time"

	sqlstore "github.com/FonovAD/Prototype/internal/store/SQLstore"
	"github.com/stretchr/testify/assert"
)

func TestLinkRepository_Create(t *testing.T) {
	// databasePath := "./test"
	db, teardown := sqlstore.SetupTestDB(t)
	defer teardown("users", "links")
	s := sqlstore.New(db, time.Millisecond*100)
	ctxb := context.Background()

	u, err := s.User().Create(ctxb)
	assert.NoError(t, err)

	originLink := "http://hello.world"
	shortLink := "http://world.hello"
	l, err := s.Link().Create(ctxb, u.UID, originLink, shortLink)
	assert.NoError(t, err)
	assert.Equal(t, originLink, l.OriginLink)
	assert.Equal(t, shortLink, l.ShortLink)

	l, err = s.Link().Create(ctxb, u.UID, originLink, "")
	assert.Error(t, err)
}

func TestLinkRepository_GetByUID(t *testing.T) {
	// databasePath := "./test"
	db, teardown := sqlstore.SetupTestDB(t)
	defer teardown("users", "links")
	s := sqlstore.New(db, time.Millisecond*100)
	ctxb := context.Background()

	u, err := s.User().Create(ctxb)
	assert.NoError(t, err)

	originLink := "http://hello.world"
	shortLink := "http://world.hello"
	l, err := s.Link().Create(ctxb, u.UID, originLink, shortLink)
	assert.NoError(t, err)

	links, err := s.Link().GetByUID(ctxb, u.UID)
	assert.NoError(t, err)
	assert.Len(t, links, 1)
	assert.Equal(t, l, links[0])
}

func TestLinkRepository_GetByOriginLink(t *testing.T) {
	// databasePath := "./test"
	db, teardown := sqlstore.SetupTestDB(t)
	defer teardown("users", "links")
	s := sqlstore.New(db, time.Millisecond*100)
	ctxb := context.Background()

	u, err := s.User().Create(ctxb)
	assert.NoError(t, err)

	originLink := "http://hello.world"
	shortLink := "http://world.hello"
	expectedl, err := s.Link().Create(ctxb, u.UID, originLink, shortLink)
	assert.NoError(t, err)

	l, err := s.Link().GetByOriginLink(ctxb, expectedl.OriginLink)
	assert.NoError(t, err)
	assert.EqualValues(t, expectedl, l)
}

func TestLinkRepository_ShortLinkExist(t *testing.T) {
	// databasePath := "./test"
	db, teardown := sqlstore.SetupTestDB(t)
	defer teardown("users", "links")
	s := sqlstore.New(db, time.Millisecond*100)
	ctxb := context.Background()

	u, err := s.User().Create(ctxb)
	assert.NoError(t, err)

	originLink := "http://hello.world"
	shortLink := "http://world.hello"
	l, err := s.Link().Create(ctxb, u.UID, originLink, shortLink)
	assert.NoError(t, err)

	exist, err := s.Link().ShortLinkExist(ctxb, l.ShortLink)
	assert.NoError(t, err)
	assert.True(t, exist)
}

func TestLinkRepository_Delete(t *testing.T) {
	// databasePath := "./test"
	db, teardown := sqlstore.SetupTestDB(t)
	defer teardown("users", "links")
	s := sqlstore.New(db, time.Millisecond*100)
	ctxb := context.Background()

	u, err := s.User().Create(ctxb)
	assert.NoError(t, err)

	originLink := "http://hello.world"
	shortLink := "http://world.hello"
	_, err = s.Link().Create(ctxb, u.UID, originLink, shortLink)
	assert.NoError(t, err)

	err = s.Link().Delete(ctxb, originLink)
	assert.NoError(t, err)
}
