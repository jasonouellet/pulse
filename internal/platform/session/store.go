package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrSessionNotFound = errors.New("session non trouvée ou expirée")

// SessionData représente l'état utilisateur conservé en mémoire
type SessionData struct {
	UserID     string   `json:"user_id"`
	ActiveRole string   `json:"active_role"`
	ActiveClub string   `json:"active_club_id"`
	Roles      []string `json:"roles"`
}

type Store struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewStore(rdb *redis.Client, ttl time.Duration) *Store {
	return &Store{
		rdb: rdb,
		ttl: ttl,
	}
}

func (s *Store) formatKey(token string) string {
	return fmt.Sprintf("session:%s", token)
}

// Create enregistre une nouvelle session dans Redis avec un TTL
func (s *Store) Create(ctx context.Context, token string, data SessionData) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("erreur de sérialisation session: %w", err)
	}

	return s.rdb.Set(ctx, s.formatKey(token), payload, s.ttl).Err()
}

// Get récupère et désérialise les données de session
func (s *Store) Get(ctx context.Context, token string) (*SessionData, error) {
	payload, err := s.rdb.Get(ctx, s.formatKey(token)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, ErrSessionNotFound
	} else if err != nil {
		return nil, err
	}

	var data SessionData
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

// Delete détruit la session (Logout)
func (s *Store) Delete(ctx context.Context, token string) error {
	return s.rdb.Del(ctx, s.formatKey(token)).Err()
}
