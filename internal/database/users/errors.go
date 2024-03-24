package users

import "errors"

var (
	ErrCreateUser  = errors.New("cannot create user")
	ErrNoUserFound = errors.New("no user found")
)
