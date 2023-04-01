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
		viper.Get("POSTGRES_HOST").(string),
		viper.Get("POSTGRES_PORT").(string),
		viper.Get("POSTGRES_USER").(string),
		viper.Get("POSTGRES_PASSWORD").(string),
		viper.Get("POSTGRES_DB").(string))

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
