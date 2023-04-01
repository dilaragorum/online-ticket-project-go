package auth

import (
	"github.com/golang-jwt/jwt/v4"
)

type AuthorizationType string

const (
	AuthAdmin AuthorizationType = "admin"
	AuthUser  AuthorizationType = "user"
)

type UserType string

const (
	IndividualUser UserType = "individual"
	CorporateUser  UserType = "corporate"
)

type Claims struct {
	Username          string            `json:"username"`
	AuthorizationType AuthorizationType `json:"authorization_type"`
	UserType          UserType          `json:"user_type"`
	UserID            uint
	jwt.RegisteredClaims
}

func (c *Claims) IsIndividualUser() bool {
	return c.UserType == IndividualUser
}

func (c *Claims) IsCorporatedUser() bool {
	return c.UserType == CorporateUser
}

func (c *Claims) IsAdmin() bool {
	return c.AuthorizationType == AuthAdmin
}

func (c *Claims) IsNotAdmin() bool {
	return !c.IsAdmin()
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
