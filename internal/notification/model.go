package notification

import (
	"gorm.io/gorm"
	"time"
)

type Channel string

const (
	SMS   Channel = "sms"
	Email Channel = "mail"
)

type Param struct {
	Channel     Channel
	To          string
	From        string
	Title       string
	Description string
	LogMsg      string
}

type Log struct {
	ID        uint    `gorm:"primarykey"`
	Channel   Channel `gorm:"not null" json:"channel"` // SMS, EMAIL
	Log       string  `gorm:"not null" json:"log"`     // Yeni kayıt olan dilaragorum@gmail.com'a hoşgeldin maili gönderildi
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
