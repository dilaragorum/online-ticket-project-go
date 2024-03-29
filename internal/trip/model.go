package trip

import (
	"gorm.io/gorm"
	"time"
)

type Vehicle string

const (
	VehicleBus    Vehicle = "Bus"
	VehicleFlight Vehicle = "Flight"

	CapacityOfBus    = 45
	CapacityOfFlight = 189
	DefaultCapacity  = 0
)

type Filter struct {
	TripID  int       `json:"trip_id"`
	From    string    `json:"from"`
	To      string    `json:"to"`
	Vehicle Vehicle   `json:"vehicle"`
	Date    time.Time `json:"date"`
}

type Trip struct {
	ID              int       `gorm:"primaryKey" json:"id"`
	From            string    `gorm:"not null;index:,unique,composite:idx_member" json:"from"`
	To              string    `gorm:"not null;index:,unique,composite:idx_member" json:"to"`
	Vehicle         Vehicle   `gorm:"not null;index:,unique,composite:idx_member" json:"vehicle"`
	Date            time.Time `gorm:"not null;index:,unique,composite:idx_member" json:"date"`
	ArrivalDuration string    `json:"arrival_duration"`
	Capacity        uint      `gorm:"not null;check:capacity>0" json:"capacity"`
	AvailableSeat   uint      `gorm:"not null;check:available_seat>=0" json:"available_seat"`
	Price           float64   `gorm:"not null;check:price>0" json:"price"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (t *Trip) CheckAvailableSeat(numberOfTickets int) bool {
	if t.AvailableSeat == 0 || t.AvailableSeat < uint(numberOfTickets) {
		return false
	}
	return true
}

func (t *Trip) BeforeCreate(tx *gorm.DB) error {
	switch t.Vehicle {
	case VehicleFlight:
		t.Capacity = CapacityOfFlight
	case VehicleBus:
		t.Capacity = CapacityOfBus
	default:
		t.Capacity = DefaultCapacity
	}

	t.AvailableSeat = t.Capacity

	return nil
}

func (t *Trip) CheckFieldsEmpty() bool {
	return t.IsStartingPlaceEmpty() || t.IsDestinationPlaceEmpty() || t.IsDateEmpty()
}

func (t *Trip) IsStartingPlaceEmpty() bool {
	return t.From == ""
}

func (t *Trip) IsDestinationPlaceEmpty() bool {
	return t.To == ""
}

func (t *Trip) IsInvalidVehicle() bool {
	return !t.IsValidVehicle()
}

func (t *Trip) IsValidVehicle() bool {
	return t.Vehicle == VehicleFlight || t.Vehicle == VehicleBus
}

func (t *Trip) IsDateEmpty() bool {
	return t.Date.IsZero()
}

func (t *Trip) IsInvalidPrice() bool {
	return !t.IsValidPrice()
}

func (t *Trip) IsValidPrice() bool {
	return t.Price >= 0
}

func IsInvalidID(id int) bool {
	return !IsValidID(id)
}

func IsValidID(id int) bool {
	return id > 0
}
