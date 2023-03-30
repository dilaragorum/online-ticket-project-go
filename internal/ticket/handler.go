package ticket

import (
	"github.com/dilaragorum/online-ticket-project-go/internal/auth"
	"github.com/labstack/echo/v4"
	"net/http"
)

var (
	WarnWhenEmptyFields  = "You should fill empty fields"
	WarnWhenEmailInvalid = "Please enter valid email"
	WarnWhenPhoneInvalid = "Please enter valid phone number"

	WarnWhenExceedAllowedTicketToPurchaseForTwenty = "You are not allowed to purchase ticket more than 20"
	WarnWhenExceedAllowedTicketToPurchaseForFive   = "You are not allowed to purchase ticket more than 5"

	WarnWhenExceedMaleTicketNumber = "You are not allowed to purchase ticket for male more than 2"

	WarnWhenCapacityFull     = "Capacity is full. Please search another trip"
	WarnWhenTripDoesNotExist = "This trip does not exist. Please check trip information."

	WarnSystemFailureMessage = "There is something wrong. Please try again later"
	SuccessPurchasedMessage  = "Ticket was successfully purchased"
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
	claim := c.Get("claim").(auth.Claims)

	// Tickets olması gerekecek.
	tickets := new([]Ticket)

	if err := c.Bind(&tickets); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	for _, ticket := range *tickets {
		if ticket.CheckFieldsEmpty() {
			return c.String(http.StatusBadRequest, WarnWhenEmptyFields)
		}

		if ticket.IsEmailInvalid() {
			return c.String(http.StatusBadRequest, WarnWhenEmailInvalid)
		}

		if ticket.IsPhoneNumberInvalid() {
			return c.String(http.StatusBadRequest, WarnWhenPhoneInvalid)
		}
	}

	if err := ti.service.Purchase(c.Request().Context(), *tickets, claim); err != nil {
		switch err {
		case ErrExceedAllowedTicketToPurchaseForTwenty:
			return c.String(http.StatusBadRequest, WarnWhenExceedAllowedTicketToPurchaseForTwenty)
		case ErrExceedAllowedTicketToPurchaseForFive:
			return c.String(http.StatusBadRequest, WarnWhenExceedAllowedTicketToPurchaseForFive)
		case ErrExceedMaleTicketNumber:
			return c.String(http.StatusBadRequest, WarnWhenExceedMaleTicketNumber)
		case ErrNoCapacity:
			return c.String(http.StatusBadRequest, WarnWhenCapacityFull)
		case ErrTripNotFound:
			return c.String(http.StatusBadRequest, WarnWhenTripDoesNotExist)
		default:
			return c.String(http.StatusInternalServerError, WarnSystemFailureMessage)
		}
	}

	return c.String(http.StatusOK, SuccessPurchasedMessage)
}
