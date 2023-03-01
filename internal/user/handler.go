package user

import (
	"errors"
	"github.com/dilaragorum/online-ticket-project-go/internal/aut"
	"github.com/dilaragorum/online-ticket-project-go/internal/mail"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

var (
	WarnWhenEmailOrUsernameIsNotUnique = "This email or username has already been registered"

	WarnInternalServerError = "an error occurred please try again later"
	WarnEmptyUserName       = "Username cannot be empty"
	WarnNonValidEmail       = "Please enter valid email address"
	WarnPasswordLength      = "Password should be between 5 and 12 characters"

	WarnNonValidCredentials  = "Please enter valid username or password"
	WarnWhenUsernameNotFound = "Invalid username, please enter valid user name"
	WarnEmailCannotSend      = "Email could not be sent"
	SuccessLoginMessage      = "Congratulations, you have successfully logged into the system."
)

type handler struct {
	userService  Service
	mailService  mail.Service
	JwtSecretKey string
}

func NewHandler(e *echo.Echo, userService Service, mailService mail.Service, jwtSecretKey string) *handler {
	h := handler{
		userService:  userService,
		mailService:  mailService,
		JwtSecretKey: jwtSecretKey,
	}

	//user := e.Group("/user", h.userMiddleware)

	e.POST("/register", h.Register)
	e.POST("/login", h.Login)
	e.GET("/logout", h.Logout)

	return &h
}

func (h *handler) Register(c echo.Context) error {
	user := new(User)

	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if user.IsNameEmpty() {
		return c.String(http.StatusBadRequest, WarnEmptyUserName)
	}

	if user.IsEmailInvalid() {
		return c.String(http.StatusBadRequest, WarnNonValidEmail)
	}

	if user.IsPasswordInvalid() {
		return c.String(http.StatusBadRequest, WarnPasswordLength)
	}

	password, err := user.HashPassword()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	user.Password = password

	requestCtx := c.Request().Context()

	err = h.userService.Register(requestCtx, user)
	if err != nil {
		switch {
		case errors.Is(err, ErrDuplicatedValue):
			return c.String(http.StatusBadRequest, WarnWhenEmailOrUsernameIsNotUnique)
		default:
			return c.String(http.StatusInternalServerError, WarnInternalServerError)
		}
	}

	err = h.mailService.SendWelcomeMail(requestCtx, user.Email)
	if errors.Is(err, mail.MailCanNotSent) {
		return c.String(http.StatusInternalServerError, WarnEmailCannotSend)
	}

	return c.NoContent(http.StatusCreated)
}

func (h *handler) Login(c echo.Context) error {
	var credentials aut.Credentials
	if err := c.Bind(&credentials); err != nil {
		return c.String(http.StatusBadRequest, WarnNonValidCredentials)
	}

	user, err := h.userService.Login(c.Request().Context(), credentials)
	if err != nil {
		switch {
		case errors.Is(err, ErrUsernameNotFound):
			return c.String(http.StatusNotFound, WarnWhenUsernameNotFound)
		case errors.Is(err, ErrUsernameOrPasswordInvalid):
			return c.String(http.StatusUnauthorized, WarnNonValidCredentials)
		default:
			return c.String(http.StatusInternalServerError, WarnInternalServerError)
		}
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := &jwt.NumericDate{Time: time.Now().Add(time.Hour)}
	claims := aut.Claims{
		Username:          user.UserName,
		AuthorizationType: user.AuthorizationType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expirationTime,
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	// Header(algorithm + JWT) + Payload(Claim)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string = secretkey + Header + Claim

	tokenString, err := token.SignedString([]byte(h.JwtSecretKey))
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = expirationTime.Time
	c.SetCookie(cookie)

	return c.String(http.StatusOK, SuccessLoginMessage)
}

func (h *handler) Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Expires = time.Now()
	return c.String(http.StatusOK, "You have successfully logout")
}

