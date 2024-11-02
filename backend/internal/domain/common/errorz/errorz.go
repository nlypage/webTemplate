package errorz

import "errors"

var (
	EmailAlreadyExists = errors.New("email already exists")
)
