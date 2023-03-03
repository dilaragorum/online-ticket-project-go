package ticket

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

var (
	WarnWhenEmptyFields  = "You should fill empty fields"
	WarnWhenEmailInvalid = "Please enter valid email"
	WarnWhenPhoneInvalid = "Please enter valid phone number"

	SuccessPurchasedMessage = "Ticket was successfully purchased"
)

type handler struct {
	service Service
}

func NewHandler(e *echo.Echo, service Service) *handler {
	h := handler{service: service}

	e.POST("/purchase/:id", h.Purchase)

	return &h
}

func (ti *handler) Purchase(c echo.Context) error {
	ticket := new(Ticket)

	if err := c.Bind(&ticket); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if ticket.CheckFieldsEmpty() {
		return c.String(http.StatusBadRequest, WarnWhenEmptyFields)
	}

	if ticket.IsEmailInvalid() {
		return c.String(http.StatusBadRequest, WarnWhenEmailInvalid)
	}

	if ticket.IsPhoneNumberInvalid() {
		return c.String(http.StatusBadRequest, WarnWhenPhoneInvalid)
	}

	if err := ti.service.Purchase(c.Request().Context(), ticket); err != nil {
		return c.String(http.StatusInternalServerError, "There is something wrong")
	}

	return c.String(http.StatusOK, SuccessPurchasedMessage)
}
