package ticket

import (
	"gorm.io/gorm"
	"time"
)

type Ticket struct {
	ID                 int
	TripID             int
	NumberOfSoldTicket int
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}
