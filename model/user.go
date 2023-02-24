package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/mail"
	"time"
)

type Channel string

const (
	ChannelSMS   Channel = "SMS"
	ChannelEMAIL Channel = "EMAIL"
)

type AuthorizationType string

const (
	AuthAdmin AuthorizationType = "admin"
	AuthUser  AuthorizationType = "user"
)

const (
	MinPasswordLen int = 5
	MaxPasswordLen int = 12
)

type NotificationLog struct {
	ID        uint    `gorm:"primarykey"`
	Channel   Channel `gorm:"not null" json:"channel"` // SMS, EMAIL
	Log       string  `gorm:"not null" json:"log"`     // Yeni kayıt olan dilaragorum@gmail.com'a hoşgeldin maili gönderildi
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	ID                uint              `gorm:"primarykey"`
	UserName          string            `gorm:"not null;unique" json:"user_name"`
	Password          string            `gorm:"not null" json:"password"`
	AuthorizationType AuthorizationType `gorm:"check: authorization_type in('admin','user')" json:"authorization_type"`
	Email             string            `gorm:"unique" json:"email"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

func (u *User) IsNameEmpty() bool {
	return u.UserName == ""
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
