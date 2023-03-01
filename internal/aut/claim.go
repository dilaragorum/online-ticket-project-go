package aut

import (
	"github.com/golang-jwt/jwt/v4"
)

type AuthorizationType string

const (
	AuthAdmin AuthorizationType = "admin"
	AuthUser  AuthorizationType = "user"
)

type Claims struct {
	Username          string            `json:"username"`
	AuthorizationType AuthorizationType `json:"authorization_type"`
	jwt.RegisteredClaims
}

func (c *Claims) IsAdmin() bool {
	return c.AuthorizationType == AuthAdmin
}

func (c *Claims) IsUser() bool {
	return c.AuthorizationType == AuthAdmin
}

func (c *Claims) IsUserOrAdmin() bool {
	return c.IsAdmin() || c.IsUser()
}

func (c *Claims) IsUnknownTypeUser() bool {
	return !c.IsUserOrAdmin()
}
