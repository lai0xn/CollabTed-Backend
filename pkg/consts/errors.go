package consts

import "errors"

var (
	EMAIL_IN_USE        = errors.New("Email already in use")
	INVALID_CREDENTIALS = errors.New("Invalid credentials")
)
