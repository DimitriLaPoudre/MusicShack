package model

import "errors"

var (
	ErrEmailDuplicate = errors.New("email already exists")
	ErrUserNotFound   = errors.New("user not found")
	ErrBadPassword    = errors.New("incorrect password")
	ErrBadToken       = errors.New("invalid token")
	ErrUnknown        = errors.New("unknown error")
)
