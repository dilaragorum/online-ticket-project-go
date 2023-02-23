package handler

import (
	"errors"
	"github.com/dilaragorum/online-ticket-project-go/model"
	"github.com/dilaragorum/online-ticket-project-go/service"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"time"
)

var (
	WarnMessageWhenUserNameIsNotUnique = "This username was taken before"
	WarnMessageWhenEmailIsNotUnique    = "This email has already been registered"
	WarnInternalServerError            = "an error occurred please try again later"
	WarnEmptyUserName                  = "Username cannot be empty"
	WarnNonValidEmail                  = "Please enter valid email address"
	WarnPasswordLength                 = "Password should be between 5 and 12 characters"
	WarnNonValidCredentials            = "Please enter valid username or password"
	WarnWhenUsernameNotFound           = "Invalid username, please enter valid user name"
	WarnEmailCannotSend                = "Email could not be sent"
	SuccessLoginMessage                = "Congratulations, you have successfully logged into the system."
)

type userHandler struct {
	userService service.UserService
	mailService service.MailService
	JwtKey      string
}

func NewUserHandler(e *echo.Echo, userService service.UserService, mailService service.MailService) *userHandler {
	h := userHandler{
		userService: userService,
		mailService: mailService,
		JwtKey:      os.Getenv("ONLINE_TICKET_GO_JWTKEY"),
	}

	e.POST("/register", h.Register)
	e.POST("/login", h.Login)
	e.GET("/logout", h.Logout)

	return &h
}

func (h *userHandler) Register(c echo.Context) error {
	user := new(model.User)

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
		case errors.Is(err, service.ErrDuplicatedEmail):
			return c.String(http.StatusBadRequest, WarnMessageWhenEmailIsNotUnique)
		case errors.Is(err, service.ErrDuplicatedUserName):
			return c.String(http.StatusBadRequest, WarnMessageWhenUserNameIsNotUnique)
		default:
			return c.String(http.StatusInternalServerError, WarnInternalServerError)
		}
	}

	err = h.mailService.SendWelcomeMail(requestCtx, user.Email)
	if errors.Is(err, service.MailCanNotSent) {
		return c.String(http.StatusInternalServerError, WarnEmailCannotSend)
	}

	return c.NoContent(http.StatusCreated)
}

func (h *userHandler) Login(c echo.Context) error {
	var credentials model.Credentials
	if err := c.Bind(&credentials); err != nil {
		return c.String(http.StatusBadRequest, WarnNonValidCredentials)
	}

	user, err := h.userService.Login(c.Request().Context(), credentials)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUsernameNotFound):
			return c.String(http.StatusNotFound, WarnWhenUsernameNotFound)
		case errors.Is(err, service.ErrUsernameOrPasswordInvalid):
			return c.String(http.StatusUnauthorized, WarnNonValidCredentials)
		default:
			return c.String(http.StatusInternalServerError, WarnInternalServerError)
		}
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

	tokenString, err := token.SignedString([]byte(h.JwtKey))
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

func (h *userHandler) Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Expires = time.Now()
	return c.String(http.StatusOK, "You have successfully logout")
}
