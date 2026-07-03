package model

import "errors"

var (
	ErrEmailDuplicate = errors.New("email already exists")
	ErrUserNotFound   = errors.New("user not found")
	ErrBadPassword    = errors.New("incorrect password")
	ErrBadToken       = errors.New("invalid token")

	ErrPluginNotFound   = errors.New("plugin not found")
	ErrInvalidUrlPlugin = errors.New("url not recognized by plugins")

	ErrInvalidInput = errors.New("invalid input")
	ErrConflict     = errors.New("conflict")
	ErrNotFound     = errors.New("not found")
	ErrInternal     = errors.New("internal_error")

	ErrUnknown = errors.New("unknown error")
)
