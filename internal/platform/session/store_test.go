package session_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pulse/internal/platform/session"
)

func setupTestStore(t *testing.T) (*session.Store, *miniredis.Miniredis) {
	t.Helper()

	mr, err := miniredis.Run()
	require.NoError(t, err)

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	t.Cleanup(func() {
		_ = rdb.Close()
		mr.Close()
	})

	store := session.NewStore(rdb, 15*time.Minute)
	return store, mr
}

func TestStore_CreateAndGet_Success(t *testing.T) {
	store, _ := setupTestStore(t)
	ctx := context.Background()
	token := "sess_test_123"

	expectedData := session.SessionData{
		UserID:     "user_42",
		ActiveRole: "COACH",
		ActiveClub: "club_rimouski",
		Roles:      []string{"COACH", "GUARDIAN"},
	}

	err := store.Create(ctx, token, expectedData)
	require.NoError(t, err)

	data, err := store.Get(ctx, token)
	require.NoError(t, err)
	assert.Equal(t, expectedData.UserID, data.UserID)
	assert.Equal(t, expectedData.ActiveRole, data.ActiveRole)
	assert.Equal(t, expectedData.ActiveClub, data.ActiveClub)
	assert.Equal(t, expectedData.Roles, data.Roles)
}

func TestStore_Get_NotFound(t *testing.T) {
	store, _ := setupTestStore(t)
	ctx := context.Background()

	_, err := store.Get(ctx, "token_inexistant")
	assert.ErrorIs(t, err, session.ErrSessionNotFound)
}

func TestStore_Delete_Success(t *testing.T) {
	store, _ := setupTestStore(t)
	ctx := context.Background()
	token := "sess_to_delete"

	_ = store.Create(ctx, token, session.SessionData{UserID: "user_delete"})

	err := store.Delete(ctx, token)
	require.NoError(t, err)

	_, err = store.Get(ctx, token)
	assert.ErrorIs(t, err, session.ErrSessionNotFound)
}
