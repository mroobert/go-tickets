package signup

import (
	fbauthn "firebase.google.com/go/v4/auth"
)

type user struct {
	UID         string
	Email       string
	DisplayName string
}

func fbToUser(u *fbauthn.UserRecord) user {
	user := user{
		UID:         u.UID,
		Email:       u.Email,
		DisplayName: u.DisplayName,
	}
	return user
}
