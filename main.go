package main

import (
	"encoding/json"
	"github.com/dilaragorum/online-ticket-project-go/database"
	"github.com/dilaragorum/online-ticket-project-go/handler"
	"github.com/dilaragorum/online-ticket-project-go/repository"
	"github.com/dilaragorum/online-ticket-project-go/service"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

var jwtKey = []byte("my_secret_key")

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
}

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

	//e.POST("/signin", signIn)
	//e.POST("/logout", logOut)
	e.Logger.Fatal(e.Start(":8080"))

}

func signIn(c echo.Context) error {
	var credentials Credentials
	err := json.NewDecoder(c.Request().Body).Decode(&credentials)
	if err != nil {
		c.NoContent(http.StatusBadRequest)
		return err
	}

	// Map tuttuğumuz için buradan kontrol ediyoruz. Diğer durumda Database'e sormak gerekiyor.
	expectedPassword, ok := users[credentials.UserName]

	// Bu user var mı var ise girdiği password doğru mu
	if !ok || expectedPassword != credentials.Password {
		c.NoContent(http.StatusUnauthorized)
		return nil
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes

	expirationTime := &jwt.NumericDate{Time: time.Now().Add(5 * time.Minute)}
	claims := Claims{
		Username: credentials.UserName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expirationTime,
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	// Header(algorithm + JWT) + Payload(Claim)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string = secretkey + Header + Claim
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.NoContent(http.StatusInternalServerError)
		return err
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself

	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = expirationTime.Time
	c.SetCookie(cookie)

	c.String(http.StatusOK, "Başarılı bir şekilde giriş yapıldı")
	return nil
}

func logOut(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Expires = time.Now()
	c.String(http.StatusOK, "Başarılı bir şekilde çıkış yapıldı")
	return nil
}
