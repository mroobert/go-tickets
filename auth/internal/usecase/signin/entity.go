package signin

import "time"

type token struct {
	AuthTime int64
	Issuer   string
	Audience string
	Expires  int64
	IssuedAt int64
	Subject  string
	UID      string
}

func (t token) isOld() bool {
	// Return error if the sign-in is older than 5 minutes.
	signInTime := time.Now().Unix() - t.AuthTime
	return signInTime > 5*60
}
