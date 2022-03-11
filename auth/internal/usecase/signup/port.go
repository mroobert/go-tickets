package signup

import "context"

// (Port) Service defines how the interaction between the "core" and the "signup http handler" has to be done.
type Service interface {
	// SignUp returns the newly created user.
	SignUp(context.Context, newUser) (user, error)
}

// (Port) AuthnProvider defines how the interaction between the "core" and the "authn provider" has to be done.
type AuthnProvider interface {
	// Create inserts a new user into the authentication provider.
	Create(context.Context, newUser) (user, error)
}
