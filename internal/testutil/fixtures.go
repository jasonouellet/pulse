package testutil

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/redis/go-redis/v9"
)

// SetupTestDB crée un pool Postgres mocké
func SetupTestDB(t *testing.T) pgxmock.PgxPoolIface {
	t.Helper()
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Erreur pgxmock: %v", err)
	}
	t.Cleanup(func() {
		mock.Close()
	})
	return mock
}

// SetupTestRedis crée un vrai serveur Redis in-memory
func SetupTestRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	t.Helper()
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Erreur miniredis: %v", err)
	}
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	t.Cleanup(func() {
		client.Close()
		mr.Close()
	})

	return client, mr
}
