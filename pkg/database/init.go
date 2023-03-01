package database

import (
	"fmt"
	"github.com/dilaragorum/online-ticket-project-go/internal/notification"
	"github.com/dilaragorum/online-ticket-project-go/internal/ticket"
	model2 "github.com/dilaragorum/online-ticket-project-go/internal/trip"
	"github.com/dilaragorum/online-ticket-project-go/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var db *gorm.DB

func Setup() (*gorm.DB, error) {
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_PORT", "5432")
	os.Setenv("POSTGRES_USER", "online_ticket_user")
	os.Setenv("POSTGRES_PASSWORD", "postgres")
	os.Setenv("POSTGRES_DB", "online_ticket_api")

	connstr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))

	var err error
	db, err = gorm.Open(postgres.Open(connstr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate() {
	if err := db.AutoMigrate(&user.User{}, &notification.Log{}, &model2.Trip{}, &ticket.Ticket{}); err != nil {
		panic(err)
	}
}
