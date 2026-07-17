package model

import "errors"

var (
	ErrEmailDuplicate  = errors.New("email already exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrUserRoleInvalid = errors.New("user role invalid")
	ErrBadPassword     = errors.New("incorrect password")
	ErrBadToken        = errors.New("invalid token")
	ErrExpiredToken    = errors.New("expired token")

	ErrPluginNotFound         = errors.New("plugin not found")
	ErrInvalidUrlPlugin       = errors.New("url not recognized by plugins")
	ErrPluginRateLimit        = errors.New("plugin rate limit reached")
	ErrPluginSearchEmptyQuery = errors.New("search query empty")

	ErrInvalidInput = errors.New("invalid input")
	ErrConflict     = errors.New("conflict")
	ErrNotFound     = errors.New("not found")
	ErrInternal     = errors.New("internal error")

	ErrAdminAlreadyExist = errors.New("admin already exist")

	ErrUnknown = errors.New("unknown error")
)
