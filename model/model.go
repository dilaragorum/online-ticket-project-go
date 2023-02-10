package model

import "gorm.io/gorm"

type User struct {
	UserName string `gorm:"not null;unique" json:"user_name"`
	Password string `gorm:"not null;check: length(password) > 8" json:"password"`
	UserType string `gorm:"check: user_type in('admin','user')" json:"user_type""`
	Email    string `gorm:"unique" json:"email"`
	gorm.Model
}
