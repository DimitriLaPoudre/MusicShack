package model

import "github.com/google/uuid"

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type User struct {
	ID       uuid.UUID
	Username string
	Password string
	HiRes    bool
	Role     UserRole
}

type PartialUser struct {
	ID       uuid.UUID
	Username *string
	Password *string
	HiRes    *bool
	Role     *UserRole
}

type FilterUser struct {
	ID       *uuid.UUID
	Username *string
	Password *string
	HiRes    *bool
	Role     *UserRole
}
