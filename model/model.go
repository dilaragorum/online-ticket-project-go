package model

import (
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type User struct {
	UserName string `gorm:"not null;unique" json:"user_name"`
	Password string `gorm:"not null;check: length(password) > 8" json:"password"`
	UserType string `gorm:"check: user_type in('admin','user')" json:"user_type""`
	Email    string `gorm:"unique" json:"email"`
	gorm.Model
}

type Credentials struct {
	Password string `json:"password"`
	UserName string `json:"user_name"`
}

type Claims struct {
	Username string `json:"username"`
	UserType string `json:"user_type"`
	jwt.RegisteredClaims
}
