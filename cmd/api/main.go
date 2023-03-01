package main

import (
	"github.com/dilaragorum/online-ticket-project-go/internal/admin"
	"github.com/dilaragorum/online-ticket-project-go/internal/mail"
	"github.com/dilaragorum/online-ticket-project-go/internal/notification"
	"github.com/dilaragorum/online-ticket-project-go/internal/trip"
	"github.com/dilaragorum/online-ticket-project-go/internal/user"
	"github.com/dilaragorum/online-ticket-project-go/pkg/database"
	mail2 "github.com/dilaragorum/online-ticket-project-go/pkg/mail"
	"github.com/labstack/echo/v4"
	"log"
	"os"
)

func main() {
	e := echo.New()

	connectionPool, err := database.Setup()
	if err != nil {
		log.Fatal(err)
	}
	database.Migrate()

	jwtSecretKey := os.Getenv("ONLINE_TICKET_GO_JWTKEY")

	// MAIL
	mailClient := mail2.NewMail()
	mailRepository := notification.NewNotificationRepository(connectionPool)
	mailService := mail.NewService(mailClient, mailRepository)

	// TRİP
	tripRepo := trip.NewTripRepository(connectionPool)
	tripService := trip.NewTripService(tripRepo)
	trip.Handler(e, tripService)

	// USER
	userRepository := user.NewRepository(connectionPool)
	userService := user.NewUserService(userRepository)
	user.NewHandler(e, userService, mailService, jwtSecretKey)

	// ADMIN
	adminRepository := trip.NewTripRepository(connectionPool)
	adminService := admin.NewAdminService(adminRepository)
	admin.NewHandler(e, adminService, jwtSecretKey)

	e.Logger.Fatal(e.Start(":8080"))
}
