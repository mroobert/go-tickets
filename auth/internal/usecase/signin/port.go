package signin

import (
	"time"
)

// (Port) Service defines how the interaction between the "core" and the "signin http handler" has to be done.
type signInService interface {
	// Signin returns the session cookie.
	SignIn(tkn string) (Session, error)
}

// (Port) AuthnProvider defines how the interaction between the "core" and the "authn provider" has to be done.
type authnProvider interface {
	// VerifyToken verifies the signature and payload of the provided token.
	VerifyToken(tkn string) (token, error)
	// SessionCookie creates a new session cookie from the given token and expiry duration.
	SessionCookie(tkn string, expiresIn time.Duration) (Session, error)
}
