package sqlstore_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/FonovAD/Prototype/internal/models"
	sqlstore "github.com/FonovAD/Prototype/internal/store/SQLstore"
	"github.com/stretchr/testify/assert"
)

// INTEGRATION TESTS
func TestUserRepository_Create(t *testing.T) {
	// databasePath := "./test"
	db, teardown := sqlstore.SetupTestDB(t)
	defer teardown("users", "links")
	s := sqlstore.New(db, time.Millisecond*100)
	ctxb := context.Background()

	u, err := s.User().Create(ctxb)
	fmt.Println(u)
	assert.NoError(t, err)
}

func TestUserRepository_GetByUID(t *testing.T) {
	db, teardown := sqlstore.SetupTestDB(t)
	defer teardown("users", "links")
	s := sqlstore.New(db, time.Millisecond*100)
	ctxb := context.Background()

	expectedU, err := s.User().Create(ctxb)
	// fmt.Println(u)
	assert.NoError(t, err)

	u, err := s.User().GetByUID(ctxb, expectedU.UID)
	assert.NoError(t, err)
	assert.Equal(t, expectedU, u)
}

func TestUserRepository_GetByRole(t *testing.T) {
	db, teardown := sqlstore.SetupTestDB(t)
	defer teardown("users", "links")
	s := sqlstore.New(db, time.Millisecond*100)
	ctxb := context.Background()

	expectedU, err := s.User().Create(ctxb)
	assert.NoError(t, err)

	users, err := s.User().GetByRole(ctxb, expectedU.Role, 10, 0)
	flag := false
	for _, user := range users {
		if user == *expectedU {
			flag = true
		}
	}
	assert.NoError(t, err)
	assert.Equal(t, true, flag)
}

func TestUserRepository_GetByToken(t *testing.T) {
	db, teardown := sqlstore.SetupTestDB(t)
	defer teardown("users", "links")
	s := sqlstore.New(db, time.Millisecond*100)
	ctxb := context.Background()

	expectedU, err := s.User().Create(ctxb)
	assert.NoError(t, err)

	u, err := s.User().GetByToken(ctxb, expectedU.Token)
	assert.NoError(t, err)
	assert.Equal(t, expectedU, u)
}

func TestUserRepository_CheckByToken(t *testing.T) {
	db, teardown := sqlstore.SetupTestDB(t)
	defer teardown("users", "links")
	s := sqlstore.New(db, time.Millisecond*100)
	ctxb := context.Background()

	expectedU, err := s.User().Create(ctxb)
	assert.NoError(t, err)

	exist, err := s.User().CheckByToken(ctxb, expectedU.Token)
	assert.NoError(t, err)
	assert.True(t, exist)
}

func TestUserRepository_Delete(t *testing.T) {
	db, teardown := sqlstore.SetupTestDB(t)
	defer teardown("users", "links")
	s := sqlstore.New(db, time.Millisecond*100)
	ctxb := context.Background()

	user1, err := s.User().Create(ctxb)
	assert.NoError(t, err)

	err = s.User().Delete(ctxb, user1.Token)
	assert.NoError(t, err)

	deletedU, err := s.User().GetByToken(ctxb, user1.Token)
	assert.Equal(t, models.User{}, *deletedU)

	user2, err := s.User().Create(ctxb)
	assert.NoError(t, err)

	err = s.User().Delete(ctxb, user2.UID)
	assert.NoError(t, err)

	deletedU2, err := s.User().GetByToken(ctxb, user1.Token)
	assert.Equal(t, models.User{}, *deletedU2)

}
