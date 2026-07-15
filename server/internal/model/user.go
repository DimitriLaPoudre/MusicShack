package model

import (
	"fmt"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

func (r UserRole) IsValid() error {
	switch r {
	case UserRoleAdmin, UserRoleUser:
		return nil
	default:
		return fmt.Errorf("UserRole.IsValid: %s: %w", r, ErrUserRoleInvalid)
	}
}

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

type UserFilter struct {
	ID       *uuid.UUID
	Username *string
	Password *string
	HiRes    *bool
	Role     *UserRole
}
