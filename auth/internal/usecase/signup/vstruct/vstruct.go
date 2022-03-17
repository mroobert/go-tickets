// Package vstruct provides the "value objects" used by the signup.
package vstruct

import (
	"fmt"
	"net/mail"
	"unicode"

	fbauthn "firebase.google.com/go/v4/auth"
)

// SignUpUser reprezents a "value object" inside domain.
type SignUpUser struct {
	email       string
	password    string
	displayName string
}

// NewSignUpUser creates a new SignUpUser that is in a valid state.
func NewSignUpUser(email string, password string, displayName string) (SignUpUser, error) {
	if email == "" {
		return SignUpUser{}, fmt.Errorf("email must be a non-empty string")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return SignUpUser{}, err
	}

	if password == "" {
		return SignUpUser{}, fmt.Errorf("password must be a non-empty string")
	}
	sixOrMore, number, upper, special := parsePassword(password)
	if !(sixOrMore) {
		return SignUpUser{}, fmt.Errorf("password must be 6 or more characters long")
	}
	if !(number) {
		return SignUpUser{}, fmt.Errorf("password must contain a number")
	}
	if !(upper) {
		return SignUpUser{}, fmt.Errorf("password must contain an upper letter")
	}
	if !(special) {
		return SignUpUser{}, fmt.Errorf("password must contain a special character")
	}

	if displayName == "" {
		return SignUpUser{}, fmt.Errorf("display name must be a non-empty string")
	}

	return SignUpUser{
		email:       email,
		password:    password,
		displayName: displayName,
	}, nil
}

// parsePassword takes the incoming password and validates it.
func parsePassword(password string) (sixOrMore, number, upper, special bool) {
	letters := 0
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			letters++
		default:
			return false, false, false, false
		}
	}
	sixOrMore = letters >= 6
	return
}

//
func ToFirebaseUser(su SignUpUser) fbauthn.UserToCreate {
	newUser := fbauthn.UserToCreate{}
	newUser.Email(su.email)
	newUser.Password(su.password)
	newUser.DisplayName(su.displayName)

	return newUser
}
