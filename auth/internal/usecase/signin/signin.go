// Package signin contains all the components needed to
// fulfill the signin use case.
package signin

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	fbauthn "firebase.google.com/go/v4/auth"
	"github.com/mroobert/go-tickets/auth/internal/foundation/web"
)

//? ==============================================================================
//? Ports

// (Port) Service defines how the interaction between the "core" and the "signin http handler" has to be done.
type SignInService interface {
	// Signin returns the session cookie.
	SignIn(tkn string) (session, error)
}

// (Port) AuthnProvider defines how the interaction between the "core" and the "authn provider" has to be done.
type SignInProvider interface {
	// VerifyToken verifies the signature and payload of the provided token.
	VerifyToken(tkn string) (token, error)
	// SessionCookie creates a new session cookie from the given token and expiry duration.
	SessionCookie(tkn string, expiresIn time.Duration) (session, error)
}

//? ==============================================================================
//? Core models
type session struct {
	Value     string
	ExpiresIn time.Duration
}

type token struct {
	AuthTime int64
	Issuer   string
	Audience string
	Expires  int64
	IssuedAt int64
	Subject  string
	UID      string
}

func toToken(fbToken *fbauthn.Token) token {
	t := token{
		AuthTime: fbToken.AuthTime,
		Issuer:   fbToken.Issuer,
		Audience: fbToken.Audience,
		Expires:  fbToken.Expires,
		IssuedAt: fbToken.IssuedAt,
		Subject:  fbToken.Subject,
		UID:      fbToken.UID,
	}

	return t
}

//? ==============================================================================
//? Core services

// Service represents "signin" core service.
type service struct {
	p SignInProvider
}

// NewService creates a "signin" core service with the necessary dependencies.
func NewService(p SignInProvider) SignInService {
	return &service{p: p}
}

// SignIn returns the session cookie.
func (s *service) SignIn(token string) (session, error) {
	decoded, err := s.p.VerifyToken(token)
	if err != nil {
		//return v1Web.NewRequestError(auth.ErrInvalidToken, http.StatusUnauthorized)
		return session{}, err
	}
	// Return error if the sign-in is older than 5 minutes.
	signInTime := time.Now().Unix() - decoded.AuthTime
	if signInTime > 5*60 {
		return session{}, fmt.Errorf("recent sign-in required: %w", err)
	}

	// Set session expiration to 2 days.
	expiresIn := time.Hour * 24 * 2

	// Create the session cookie. This will also verify the ID token in the process.
	// The session cookie will have the same claims as the ID token.
	ses, err := s.p.SessionCookie(token, expiresIn)
	if err != nil {
		return session{}, fmt.Errorf("failed to create a session cookie: %w", err)
	}

	return ses, nil
}

//? ==============================================================================
//? Adapters

// (Adapter) Handler transforms a "signin http request" into a "call on core service".
type Handler struct {
	service SignInService
}

// NewHandler construct a handler for signin.
func NewHandler(service SignInService) *Handler {
	return &Handler{
		service: service,
	}
}

// SignIn returns the session cookie as http response.
func (h *Handler) SignIn(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Expecting: jwt <token>
	authStr := r.Header.Get("authorization")
	// Parse the authorization header.
	parts := strings.Split(authStr, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "jwt" {
		err := errors.New("expected authorization header format: jwt <token>")
		//return v1Web.NewRequestError(err, http.StatusUnauthorized)
		return err
	}
	scookie, err := h.service.SignIn(parts[1])
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    scookie.Value,
		MaxAge:   int(scookie.ExpiresIn.Seconds()),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	if err != nil {
		return err
	}

	status := struct {
		Status string
	}{
		Status: "Success",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}

// (Adapter) FirebaseAuthn transforms a "core service call" into a "call on firebase authn provider".
type Firebase struct {
	client *fbauthn.Client
}

// NewFirebase sets a firebase authentication client for signin use case.
func NewFirebase(client *fbauthn.Client) *Firebase {
	return &Firebase{
		client: client,
	}
}

// VerifyToken verifies the signature and payload of the provided firebase token.
func (fb Firebase) VerifyToken(tkn string) (token, error) {
	decoded, err := fb.client.VerifyIDToken(context.Background(), tkn)
	if err != nil {
		//return v1Web.NewRequestError(auth.ErrInvalidToken, http.StatusUnauthorized)
		return token{}, err
	}

	t := toToken(decoded)
	return t, nil
}

// SessionCookie creates a new firebase session cookie from the given token and expiry duration.
func (fb Firebase) SessionCookie(tkn string, expiresIn time.Duration) (session, error) {
	// Create the session cookie. This will also verify the ID token in the process.
	// The session cookie will have the same claims as the ID token.
	value, err := fb.client.SessionCookie(context.Background(), tkn, expiresIn)
	if err != nil {
		//return fmt.Errorf("failed to create a session cookie: %w", err)
		return session{}, err
	}

	return session{
		Value:     value,
		ExpiresIn: expiresIn,
	}, nil
}
