package errorz

import "errors"

var (
	EmailAlreadyTaken = errors.New("email already taken")
	AuthHeaderIsEmpty = errors.New("auth header is empty")
	Forbidden         = errors.New("forbidden")
)
