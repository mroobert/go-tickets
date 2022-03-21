package signin

import (
	"context"
	"fmt"
	"net/http"
	"time"

	fbauthn "firebase.google.com/go/v4/auth"
	"github.com/mroobert/go-tickets/auth/internal/foundation/web"
	"github.com/mroobert/go-tickets/auth/internal/usecase/signin/vstruct"
)

// (Adapter) HttpHandler transforms a "signin http request" into a "call on signin core service".
func HttpHandler(s signInService) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		// Expecting: jwt <token>
		authStr := r.Header.Get("authorization")

		// Parse the authorization header.
		token, err := web.ExtractToken(authStr)
		if err != nil {
			return fmt.Errorf("unable to sign in: %w", err)
		}

		// Generate session cookie
		scookie, err := s.SignIn(token)
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
}

// (Adapter) Firebase transforms a "core service call" into a "call on firebase authn provider".
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
		return token{}, fmt.Errorf("failed to verify the token: %w", err)
	}

	t := toToken(decoded)
	return t, nil
}

// SessionCookie creates a new firebase session cookie from the given token and expiry duration.
func (fb Firebase) SessionCookie(tkn string, expiresIn time.Duration) (vstruct.Session, error) {
	// Create the session cookie. This will also verify the ID token in the process.
	// The session cookie will have the same claims as the ID token.
	value, err := fb.client.SessionCookie(context.Background(), tkn, expiresIn)
	if err != nil {
		return vstruct.Session{}, fmt.Errorf("failed to create a session cookie on firebase: %w", err)
	}

	session, err := vstruct.NewSession(value, expiresIn)
	if err != nil {
		return vstruct.Session{}, fmt.Errorf("failed to create a session cookie: %w", err)
	}
	return session, nil
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
