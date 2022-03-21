package signin

import (
	"fmt"
	"time"
)

type Session struct {
	Value     string
	ExpiresIn time.Duration
}

// NewSession creates a new Session that is in a valid state.
func NewSession(value string, expires time.Duration) (Session, error) {
	if value == "" {
		return Session{}, fmt.Errorf("value must be a non-empty string")
	}
	if expires < time.Hour*24 {
		return Session{}, fmt.Errorf("the session must not expires in less than 24h")
	}

	return Session{
		Value:     value,
		ExpiresIn: expires,
	}, nil
}
