package signup

import (
	"context"

	"github.com/mroobert/go-tickets/auth/internal/usecase/signup/vstruct"
)

// (Port) Service defines how the interaction between the "core" and the "signup http handler" has to be done.
type Service interface {
	// SignUp returns the newly created user.
	SignUp(context.Context, vstruct.SignUpUser) (user, error)
}

// (Port) AuthnProvider defines how the interaction between the "core" and the "authn provider" has to be done.
type AuthnProvider interface {
	// Create inserts a new user into the authentication provider.
	Create(context.Context, vstruct.SignUpUser) (user, error)
}
