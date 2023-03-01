package notification

import (
	"gorm.io/gorm"
	"time"
)

type Channel string

const (
	ChannelSMS   Channel = "SMS"
	ChannelEMAIL Channel = "EMAIL"
)

type Log struct {
	ID        uint    `gorm:"primarykey"`
	Channel   Channel `gorm:"not null" json:"channel"` // SMS, EMAIL
	Log       string  `gorm:"not null" json:"log"`     // Yeni kayıt olan dilaragorum@gmail.com'a hoşgeldin maili gönderildi
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
