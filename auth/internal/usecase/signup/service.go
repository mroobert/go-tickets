package signup

import (
	"context"
	"fmt"
)

// Service represents "signup" core service.
type service struct {
	ap AuthnProvider
}

// NewService creates a "signup core service" with the necessary dependencies.
func NewService(ap AuthnProvider) *service {
	return &service{ap: ap}
}

// SignUp creates a new user.
func (s *service) SignUp(ctx context.Context, su SignUpUser) (user, error) {
	u, err := s.ap.Create(ctx, su)
	if err != nil {
		return user{}, fmt.Errorf("signup: %w", err)
	}
	return u, nil
}
