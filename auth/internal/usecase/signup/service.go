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
func NewService(ap AuthnProvider) Service {
	return &service{ap: ap}
}

// SignUp creates a new user.
func (s *service) SignUp(ctx context.Context, nu newUser) (user, error) {
	err := nu.Validate()
	if err != nil {
		return user{}, err
	}
	u, err := s.ap.Create(ctx, nu)
	if err != nil {
		return user{}, fmt.Errorf("signup service: %w", err)
	}
	return u, nil
}
