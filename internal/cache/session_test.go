package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestSessionCache(t *testing.T) {
	// 1. Démarrer le serveur Redis In-Memory local
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Impossible de démarrer miniredis : %v", err)
	}
	defer mr.Close()

	// 2. Initialiser le client go-redis pointant sur l'adresse de miniredis
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	ctx := context.Background()

	// 3. Exécuter ton code métier
	err = rdb.Set(ctx, "session:user_123", "active", 10*time.Minute).Err()
	if err != nil {
		t.Fatalf("Échec de l'écriture Redis : %v", err)
	}

	// 4. Assertions
	val, err := rdb.Get(ctx, "session:user_123").Result()
	if err != nil || val != "active" {
		t.Errorf("Attendu 'active', obtenu '%s'", val)
	}

	// Optionnel : tu peux manipuler le temps écoulé pour tester le TTL !
	mr.FastForward(11 * time.Minute)

	if mr.Exists("session:user_123") {
		t.Errorf("La clé aurait dû expirer après 10 minutes")
	}
}
