package signup

import "errors"

// ErrDuplicate is used when a user already exists.
var ErrDuplicate = errors.New("user already exists")
