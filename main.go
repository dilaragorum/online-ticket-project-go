package main

import (
	"github.com/dilaragorum/online-ticket-project-go/database"
	"github.com/dilaragorum/online-ticket-project-go/handler"
	"github.com/dilaragorum/online-ticket-project-go/repository"
	"github.com/dilaragorum/online-ticket-project-go/service"
	"github.com/labstack/echo/v4"
	"log"
)

/*var jwtKey = []byte("my_secret_key")

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

type Credentials struct {
	Password string `json:"password"`
	UserName string `json:"user_name"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}*/

func main() {
	e := echo.New()

	connectionPool, err := database.Setup()
	if err != nil {
		log.Fatal(err)
	}
	database.Migrate()

	onlineTicketRepo := repository.NewDefaultRepository(connectionPool)
	onlineTicketService := service.NewDefaultService(onlineTicketRepo)
	handler.NewDefaultOnlineTicketHandler(e, onlineTicketService)

	//e.POST("/logout", logOut)
	e.Logger.Fatal(e.Start(":8080"))

}
