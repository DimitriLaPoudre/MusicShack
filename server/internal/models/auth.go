package models

import "github.com/golang-jwt/jwt/v5"

type JwtClaims struct {
	Id uint `json:"username"`
	jwt.RegisteredClaims
}

type Signup struct {
	Username string
	Password string
}

type Login struct {
	Username string
	Password string
}
