package customErrors

import "errors"

var ErrSessionInvalid error = errors.New("session id is invalid")
var ErrUserExists error = errors.New("user already exists")
var ErrInvalidPassword error = errors.New("invalid email or password")
