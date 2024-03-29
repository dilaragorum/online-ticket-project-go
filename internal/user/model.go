package user

import (
	"github.com/dilaragorum/online-ticket-project-go/internal/auth"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/mail"
	"time"
)

const (
	MinPasswordLen int = 5
	MaxPasswordLen int = 12
)

type User struct {
	ID        uint          `gorm:"primarykey"`
	UserName  string        `gorm:"not null;unique" json:"user_name"`
	Password  string        `gorm:"not null" json:"password"`
	UserType  auth.UserType `gorm:"check: user_type in('admin','individual','corporate')" json:"user_type"`
	Email     string        `gorm:"unique" json:"email"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) IsNameEmpty() bool {
	return u.UserName == ""
}

func (u *User) IsAuthTypeInvalid() bool {
	return !u.IsUserTypeValid()
}

func (u *User) IsUserTypeValid() bool {
	switch u.UserType {
	case auth.CorporateUser:
		fallthrough
	case auth.Admin:
		fallthrough
	case auth.IndividualUser:
		return true
	default:
		return false
	}
}

func (u *User) IsUserTypeInvalid() bool {
	return !u.IsUserTypeValid()
}

func (u *User) IsEmailValid() bool {
	_, err := mail.ParseAddress(u.Email)
	return err == nil
}

func (u *User) IsEmailInvalid() bool {
	return !u.IsEmailValid()
}

func (u *User) IsPasswordValid() bool {
	passLength := len(u.Password)
	return passLength >= MinPasswordLen && passLength <= MaxPasswordLen
}

func (u *User) IsPasswordInvalid() bool {
	return !u.IsPasswordValid()
}

func (u *User) HashPassword() (string, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	return string(passwordBytes), err
}
