package database

import (
	"fmt"
	"github.com/dilaragorum/online-ticket-project-go/internal/notification"
	"github.com/dilaragorum/online-ticket-project-go/internal/ticket"
	model "github.com/dilaragorum/online-ticket-project-go/internal/trip"
	"github.com/dilaragorum/online-ticket-project-go/internal/user"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Setup() (*gorm.DB, error) {
	connstr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("POSTGRES_HOST"),
		viper.GetString("POSTGRES_PORT"),
		viper.GetString("POSTGRES_USER"),
		viper.GetString("POSTGRES_PASSWORD"),
		viper.GetString("POSTGRES_DB"))

	var err error
	db, err = gorm.Open(postgres.Open(connstr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate() {
	if err := db.AutoMigrate(&user.User{}, &notification.Log{}, &model.Trip{}, &ticket.Ticket{}); err != nil {
		panic(err)
	}
}
