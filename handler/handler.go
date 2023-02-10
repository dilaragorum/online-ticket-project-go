package handler

import (
	"github.com/dilaragorum/online-ticket-project-go/model"
	"github.com/dilaragorum/online-ticket-project-go/service"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

var (
	WarnMessageWhenUserNameIsNotUnique = "This username was taken before"
	WarnMessageWhenEmailIsNotUnique    = "This email has already been registered"
	WarnInternalServerError            = "an error occurred please try again later"
)

type DefaultHandler struct {
	service service.Service
}

func NewDefaultOnlineTicketHandler(e *echo.Echo, service service.Service) *DefaultHandler {
	ot := DefaultHandler{service: service}
	e.POST("/register", ot.Register)
	return &DefaultHandler{}
}

func (ot *DefaultHandler) Register(c echo.Context) error {
	user := new(model.User)

	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if len(user.UserName) == 0 {
		return c.String(http.StatusBadRequest, "Username cannot be empty ")
	}

	if !strings.Contains(user.Email, "@") {
		return c.String(http.StatusBadRequest, "Please enter valid email address")
	}

	if len(user.Password) < 8 {
		return c.String(http.StatusBadRequest, "Password should be eight or more characters ")
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
