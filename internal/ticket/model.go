package ticket

import (
	"gorm.io/gorm"
	"net/mail"
	"regexp"
	"time"
)

type Gender string

const (
	Male   Gender = "Male"
	Female Gender = "Female"
)

type Ticket struct {
	ID     int  `gorm:"primaryKey" json:"id"`
	TripID int  `gorm:"not null" json:"trip_id"`
	UserID uint `gorm:"not null" json:"user_id"`
	Passenger
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Passenger struct {
	Gender   Gender `gorm:"not null" json:"gender"`
	FullName string `gorm:"not null" json:"full_name"`
	Email    string `gorm:"not null" json:"email"`
	Phone    string `gorm:"not null" json:"phone"`
}

func (t *Ticket) CheckFieldsEmpty() bool {
	return t.isGenderEmpty() || t.isFullNameEmpty() || t.isEmailEmpty() || t.isPhoneEmpty()
}

func (t *Ticket) isTripIDEmpty() bool {
	return t.TripID == 0
}

func (p *Passenger) isGenderEmpty() bool {
	return p.Gender == ""
}

func (p *Passenger) isFullNameEmpty() bool {
	return p.FullName == ""
}

func (p *Passenger) isEmailEmpty() bool {
	return p.Email == ""
}

func (p *Passenger) isPhoneEmpty() bool {
	return p.Phone == ""
}

func (p *Passenger) IsEmailValid() bool {
	_, err := mail.ParseAddress(p.Email)
	return err == nil
}

func (p *Passenger) IsEmailInvalid() bool {
	return !p.IsEmailValid()
}

func (p *Passenger) IsPhoneNumberValid() bool {
	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	return re.MatchString(p.Phone)
}

func (p *Passenger) IsPhoneNumberInvalid() bool {
	return !p.IsPhoneNumberValid()
}
