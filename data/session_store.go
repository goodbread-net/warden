package data

import (
	"time"

	"github.com/matthiase/warden/data/redis"
	db "github.com/redis/go-redis/v9"
)

type SessionStore interface {
	Create(string) (string, error)
	Find(string) (string, error)
	//Touch(t models.SessionToken, userID int) error
	//FindAll(userID int) ([]models.SessionToken, error)
	Revoke(string) error
}

func NewSessionStore(client *db.Client, ttl time.Duration) (SessionStore, error) {
	return &redis.SessionStore{
		Client:    client,
		Namespace: "warden:session:",
		TTL:       ttl,
	}, nil
}
