package sqlstore_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	sqlstore "github.com/FonovAD/Prototype/internal/store/SQLstore"
	"github.com/stretchr/testify/assert"
)

func TestSessionRepository_CreateSession(t *testing.T) {
	// databasePath := "./test"
	db, teardown := sqlstore.SetupTestDB(t)
	defer teardown("users", "links")
	s := sqlstore.New(db, time.Millisecond*100)
	ctxb := context.Background()

	u, err := s.User().Create(ctxb)
	fmt.Println(u)
	assert.NoError(t, err)
}
