package model

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	Username string `json:"username"`
	UserType string `json:"user_type"`
	jwt.RegisteredClaims
}
