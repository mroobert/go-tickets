package signup

import (
	"context"
	"fmt"
	"net/http"

	fbauthn "firebase.google.com/go/v4/auth"
	fberrors "firebase.google.com/go/v4/errorutils"
	"github.com/mroobert/go-tickets/auth/internal/foundation/web"
	"github.com/mroobert/go-tickets/auth/internal/usecase/signup/vstruct"
)

// (Adapter) HttpHandler transforms a "signup http request" into a "call on signup core service".
func HttpHandler(s Service) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		// decode payload
		var reqDto signUpRequestDto
		if err := web.Decode(r, &reqDto); err != nil {
			return fmt.Errorf("unable to decode payload: %w", err)
		}

		// business logic
		su, err := dtoToSignUpUser(reqDto)
		if err != nil {
			return fmt.Errorf("invalid payload: %w", err)
		}
		usr, err := s.SignUp(ctx, su)
		if err != nil {
			return fmt.Errorf("unable to signup %w", err)
		}

		// send response
		resp := userToSignUpResponseDto(usr)
		return web.Respond(ctx, w, resp, http.StatusCreated)
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
func (fb Firebase) Create(ctx context.Context, su vstruct.SignUpUser) (user, error) {
	fbUser := vstruct.ToFirebaseUser(su)
	u, err := fb.client.CreateUser(ctx, &fbUser)
	if err != nil {
		if fberrors.IsAlreadyExists(err) {
			return user{}, ErrDuplicate
		}
		return user{}, fmt.Errorf("firebase creating user: %w", err)
	}

	return fbToUser(u), nil
}

func fbToUser(u *fbauthn.UserRecord) user {
	user := user{
		UID:         u.UID,
		Email:       u.Email,
		DisplayName: u.DisplayName,
	}
	return user
}
