package errorz

import "errors"

var (
	EmailAlreadyExists = errors.New("email already exists")
	AuthHeaderIsEmpty  = errors.New("auth header is empty")
	Forbidden          = errors.New("forbidden")
)
