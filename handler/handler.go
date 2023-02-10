package handler

import (
	"encoding/json"
	"github.com/dilaragorum/online-ticket-project-go/model"
	"github.com/dilaragorum/online-ticket-project-go/service"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"time"
)

var (
	WarnMessageWhenUserNameIsNotUnique = "This username was taken before"
	WarnMessageWhenEmailIsNotUnique    = "This email has already been registered"
	WarnInternalServerError            = "an error occurred please try again later"
	WarnEmptyUserName                  = "Username cannot be empty"
	WarnNonValidEmail                  = "Please enter valid email address"
	WarnPasswordLength                 = "Password should be eight or more characters"
	WarnNonValidCredentials            = "Please enter valid username or password"
	WarnWhenUsernameNotFound           = "Invalid username, please enter valid user name"
	SuccessLoginMessage                = "Congratulations, you have successfully logged into the system."
)

type DefaultHandler struct {
	service service.Service
}

func NewDefaultOnlineTicketHandler(e *echo.Echo, service service.Service) *DefaultHandler {
	ot := DefaultHandler{service: service}
	e.POST("/register", ot.Register)
	e.POST("/login", ot.LogIn)
	e.GET("logout", ot.LogOut)
	return &DefaultHandler{}
}

func (ot *DefaultHandler) Register(c echo.Context) error {
	user := new(model.User)

	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if len(user.UserName) == 0 {
		return c.String(http.StatusBadRequest, WarnEmptyUserName)
	}

	if !strings.Contains(user.Email, "@") {
		return c.String(http.StatusBadRequest, WarnNonValidEmail)
	}

	if len(user.Password) < 8 {
		return c.String(http.StatusBadRequest, WarnPasswordLength)
	}

	register, err := ot.service.Register(c.Request().Context(), user)
	if err != nil {
		switch err.Error() {
		case service.ErrDuplicatedEmail.Error():
			return c.String(http.StatusBadRequest, WarnMessageWhenEmailIsNotUnique)
		case service.ErrDuplicatedUserName.Error():
			return c.String(http.StatusBadRequest, WarnMessageWhenUserNameIsNotUnique)
		default:
			return c.String(http.StatusInternalServerError, WarnInternalServerError)
		}
	}

	return c.JSON(http.StatusCreated, register)
}
func (ot *DefaultHandler) LogIn(c echo.Context) error {
	var credentials model.Credentials
	err := json.NewDecoder(c.Request().Body).Decode(&credentials)
	if err != nil {
		return c.String(http.StatusBadRequest, WarnNonValidCredentials)
	}

	user, err := ot.service.LogIn(c.Request().Context(), credentials)
	if err != nil {
		if err.Error() == service.ErrUsernameNotFound.Error() {
			return c.String(http.StatusNotFound, WarnWhenUsernameNotFound)
		} else if err.Error() == service.ErrUsernameOrPasswordInvalid.Error() {
			return c.String(http.StatusUnauthorized, WarnNonValidCredentials)
		}
		return c.String(http.StatusInternalServerError, WarnInternalServerError)
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expitationTime := &jwt.NumericDate{Time: time.Now().Add(5 * time.Minute)}
	claims := model.Claims{
		Username: user.UserName,
		UserType: user.UserType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expitationTime,
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	// Header(algorithm + JWT) + Payload(Claim)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string = secretkey + Header + Claim
	tokenString, err := token.SignedString([]byte("jwtKey")) //ToDo: Burada key farklı şekilde ver.
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = expitationTime.Time
	c.SetCookie(cookie)

	return c.String(http.StatusOK, SuccessLoginMessage)
}

func (ot *DefaultHandler) LogOut(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Expires = time.Now()
	return c.String(http.StatusOK, "You have successfully logout")
}
