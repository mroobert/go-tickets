package signin

import (
	"fmt"
	"time"
)

// Service represents "signin" core service.
type service struct {
	p authnProvider
}

// NewService creates a "signin" core service with the necessary dependencies.
func NewService(p authnProvider) *service {
	return &service{p: p}
}

// SignIn returns the session cookie.
func (s *service) SignIn(token string) (Session, error) {
	decoded, err := s.p.VerifyToken(token)
	if err != nil {
		//return v1Web.NewRequestError(auth.ErrInvalidToken, http.StatusUnauthorized)
		return Session{}, err
	}

	// Return error if the sign-in is older than 5 minutes.
	if decoded.isOld() {
		return Session{}, fmt.Errorf("recent sign-in required")
	}

	// Set session expiration to 2 days.
	expiresIn := time.Hour * 24 * 2

	// Create the session cookie. This will also verify the ID token in the process.
	// The session cookie will have the same claims as the ID token.
	ses, err := s.p.SessionCookie(token, expiresIn)
	if err != nil {
		return Session{}, fmt.Errorf("failed to create a session cookie: %w", err)
	}

	return ses, nil
}
