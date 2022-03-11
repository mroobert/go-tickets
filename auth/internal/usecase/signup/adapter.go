package signup

import (
	"context"
	"fmt"
	"net/http"

	fbauthn "firebase.google.com/go/v4/auth"
	fberrors "firebase.google.com/go/v4/errorutils"
	"github.com/mroobert/go-tickets/auth/internal/foundation/web"
)

// (Adapter) HttpHandler transforms a "signup http request" into a "call on signup core service".
func HttpHandler(s Service) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var nuDto newUserDto
		if err := web.Decode(r, &nuDto); err != nil {
			return fmt.Errorf("unable to decode payload: %w", err)
		}

		nu := dtoToNewUser(nuDto)
		user, err := s.SignUp(ctx, nu)
		if err != nil {
			return fmt.Errorf("signup handler[%+v]: %w", &nu, err)
		}
		dto := userToDto(user)
		return web.Respond(ctx, w, dto, http.StatusCreated)
	}
}

// (Adapter) Firebase transforms a "signup core service call" into a "call on firebase".
type Firebase struct {
	client *fbauthn.Client
}

// NewFirebase sets a firebase authentication client for signup use case.
func NewFirebase(client *fbauthn.Client) *Firebase {
	return &Firebase{
		client: client,
	}
}

// Create adds a new user in firebase with the specified properties.
func (fb Firebase) Create(ctx context.Context, nu newUser) (user, error) {
	fbUser := toFirebaseUser(nu)
	u, err := fb.client.CreateUser(ctx, &fbUser)
	if err != nil {
		if fberrors.IsAlreadyExists(err) {
			return user{}, ErrDuplicate
		}
		return user{}, fmt.Errorf("firebase: %w", err)
	}

	return fbToUser(u), nil
}
