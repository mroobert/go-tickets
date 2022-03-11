package signup

import (
	"fmt"
	"net/mail"
	"unicode"

	fbauthn "firebase.google.com/go/v4/auth"
)

type newUser struct {
	Email       string
	Password    string
	DisplayName string
}

func (nu newUser) Validate() error {
	if nu.Email == "" {
		return fmt.Errorf("email must be a non-empty string")
	}
	_, err := mail.ParseAddress(nu.Email)
	if err != nil {
		return err
	}

	if nu.Password == "" {
		return fmt.Errorf("password must be a non-empty string")
	}
	sixOrMore, number, upper, special := nu.parsePassword()
	if !(sixOrMore) {
		return fmt.Errorf("password must be 6 or more characters long")
	}
	if !(number) {
		return fmt.Errorf("password must contain a number")
	}
	if !(upper) {
		return fmt.Errorf("password must contain an upper letter")
	}
	if !(special) {
		return fmt.Errorf("password must contain a special character")
	}

	if nu.DisplayName == "" {
		return fmt.Errorf("display name must be a non-empty string")
	}

	return nil
}

func (nu newUser) parsePassword() (sixOrMore, number, upper, special bool) {
	letters := 0
	for _, c := range nu.Password {
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

func toFirebaseUser(nu newUser) fbauthn.UserToCreate {
	newUser := fbauthn.UserToCreate{}
	newUser.Email(nu.Email)
	newUser.Password(nu.Password)
	newUser.DisplayName(nu.DisplayName)

	return newUser
}
