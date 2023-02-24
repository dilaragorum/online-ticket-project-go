package model

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	Username          string            `json:"username"`
	AuthorizationType AuthorizationType `json:"authorization_type"`
	jwt.RegisteredClaims
}

func (c *Claims) IsAdmin() bool {
	return c.AuthorizationType == AuthAdmin
}
