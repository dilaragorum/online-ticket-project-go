package model

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

type Trip struct {
	From            string    `gorm:"not null;index:,unique,composite:idx_member" json:"from"`
	To              string    `gorm:"not null;index:,unique,composite:idx_member" json:"to"`
	Vehicle         Vehicle   `gorm:"not null;index:,unique,composite:idx_member" json:"vehicle"`
	Date            time.Time `gorm:"not null;index:,unique,composite:idx_member" json:"date"`
	ArrivalDuration string    `json:"arrival_duration"`
	Capacity        int       `gorm:"not null" json:"capacity"`
	Price           float64   `gorm:"not null;check:price>0" json:"price"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
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

	return nil
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

func (t *Trip) IsNotValidPrice() bool {
	return !t.IsValidPrice()
}

func (t *Trip) IsValidPrice() bool {
	return t.Price >= 0
}