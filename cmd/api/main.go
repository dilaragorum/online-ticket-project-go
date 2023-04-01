package main

import (
	"fmt"
	"github.com/dilaragorum/online-ticket-project-go/internal/auth"
	"github.com/dilaragorum/online-ticket-project-go/internal/notification"
	"github.com/dilaragorum/online-ticket-project-go/internal/payment"
	"github.com/dilaragorum/online-ticket-project-go/internal/ticket"
	"github.com/dilaragorum/online-ticket-project-go/internal/trip"
	"github.com/dilaragorum/online-ticket-project-go/internal/user"
	"github.com/dilaragorum/online-ticket-project-go/pkg/database"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"log"
)

func main() {
	e := echo.New()
	e.Use(auth.TokenMiddleware)

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	fmt.Println(viper.Get("POSTGRES_DB"))

	jwtSecretKey := viper.GetString("ONLINE_TICKET_GO_JWTKEY")

	connectionPool, err := database.Setup()
	if err != nil {
		log.Fatal(err)
	}
	database.Migrate()

	notificationRepository := notification.NewNotificationRepository(connectionPool)
	notificationService := notification.NewService(notificationRepository)

	// TRİP
	tripRepo := trip.NewTripRepository(connectionPool)
	tripService := trip.NewTripService(tripRepo)
	trip.Handler(e, tripService)

	// USER
	userRepository := user.NewRepository(connectionPool)
	userService := user.NewUserService(userRepository)
	user.NewHandler(e, userService, notificationService, jwtSecretKey)

	//PAYMENT
	paymentClient := payment.NewClient()

	// TICKET
	ticketRepo := ticket.NewTicketRepository(connectionPool)
	service := ticket.NewService(ticketRepo, notificationService, tripRepo, paymentClient)
	ticket.NewHandler(e, service)

	e.Logger.Fatal(e.Start(":8080"))
}
