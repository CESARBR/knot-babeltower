package cache

import (
	"time"

	"github.com/CESARBR/knot-babeltower/pkg/network"
)

// SessionStore abstracts the operations for storing session related data.
// It provides basic `Get/Save` methods without any explicit dependency with
// the underlying database technology.
type SessionStore interface {
	Get(email string) (string, error)
	Save(email string, id string) error
}

type sessionStore struct {
	redis          *network.Redis
	expirationTime string
}

// NewSessionStore creates a new sessionStore instance that implements SessionStore interface.
// It takes a Redis instance as dependency, which is the database abstracted by SessionStore.
func NewSessionStore(redis *network.Redis, expirationTime string) SessionStore {
	return &sessionStore{redis, expirationTime}
}

// Get retrieves a session from the database based on the user email.
func (ss *sessionStore) Get(email string) (string, error) {
	session, err := ss.redis.Get(email)
	if err != nil {
		return "", err
	}

	return session, nil
}

// Save stores a new session to the database, which is represented by a user email and the
// respective session ID.
func (ss *sessionStore) Save(email, id string) error {
	duration, err := time.ParseDuration(ss.expirationTime)
	if err != nil {
		return err
	}

	err = ss.redis.Set(email, id, duration)
	if err != nil {
		return err
	}

	return nil
}
