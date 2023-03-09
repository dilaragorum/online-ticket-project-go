package trip

import (
	"errors"
	"github.com/dilaragorum/online-ticket-project-go/internal/auth"
	"github.com/dilaragorum/online-ticket-project-go/internal/user"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

const (
	WarnNoTripMeetConditions = "There is no trip which meet your conditions."
	WarnInternalError        = "Somethings go wrong. Please try later again"

	WarnAlreadyCreatedTrip            = "This trip is already created. Please create another trip."
	WarnMessageWhenThereAreEmptyBlank = "Please fill required area"
	WarnMessageWhenInvalidVehicle     = "Please enter valid Vehicle Type"
	WarnMessageWhenInvalidPrice       = "Please enter valid price"

	WarnMessageWhenInvalidID             = "Please enter valid ID"
	WarnMessageWhenTripNotExistForDelete = "This trip does not exist or it is deleted already. "
)

type handler struct {
	tripService Service
}

func Handler(e *echo.Echo, tripService Service) *handler {
	h := handler{
		tripService: tripService,
	}

	e.GET("/trips", h.FilterTrips)
	e.POST("/trips", h.CreateTrip, auth.AdminMiddleware)
	e.DELETE("/trips/:id", h.CancelTrip, auth.AdminMiddleware)
	e.GET("/trips/sold/:id", h.GetSoldTicketNumber, auth.AdminMiddleware)

	return &h
}

func (t *handler) FilterTrips(c echo.Context) error {
	filter := Filter{}
	if err := c.Bind(&filter); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	trips, err := t.tripService.FilterTrips(c.Request().Context(), &filter)
	if err != nil {
		if errors.Is(err, user.ErrThereIsNoTrip) {
			return c.String(http.StatusBadRequest, WarnNoTripMeetConditions)
		}
		return c.String(http.StatusInternalServerError, user.WarnInternalServerError)
	}

	return c.JSON(http.StatusOK, trips)
}

func (t *handler) CreateTrip(c echo.Context) error {
	trip := new(Trip)
	if err := c.Bind(&trip); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if trip.CheckFieldsEmpty() {
		return c.String(http.StatusBadRequest, WarnMessageWhenThereAreEmptyBlank)
	}

	if trip.IsInvalidVehicle() {
		return c.String(http.StatusBadRequest, WarnMessageWhenInvalidVehicle)
	}

	if trip.IsInvalidPrice() {
		return c.String(http.StatusBadRequest, WarnMessageWhenInvalidPrice)
	}

	requestCtx := c.Request().Context()

	if err := t.tripService.CreateTrip(requestCtx, trip); err != nil {
		if errors.Is(err, ErrAlreadyCreatedTrip) {
			return c.String(http.StatusBadRequest, WarnAlreadyCreatedTrip)
		}
		return c.String(http.StatusInternalServerError, WarnInternalError)
	}

	return c.NoContent(http.StatusCreated)
}

func (t *handler) CancelTrip(c echo.Context) error {
	tripIDStr := c.Param("id")
	tripID, _ := strconv.Atoi(tripIDStr)

	if IsInvalidID(tripID) {
		return c.String(http.StatusBadRequest, WarnMessageWhenInvalidID)
	}

	requestCtx := c.Request().Context()

	if err := t.tripService.CancelTrip(requestCtx, tripID); err != nil {
		switch {
		case errors.Is(err, ErrTripNotExist):
			return c.String(http.StatusBadRequest, WarnMessageWhenTripNotExistForDelete)
		}
		return c.String(http.StatusInternalServerError, WarnInternalError)
	}

	return c.NoContent(http.StatusNoContent)
}

func (t *handler) GetSoldTicketNumber(c echo.Context) error {
	p := c.Param("id")
	param, err := strconv.Atoi(p)
	if err != nil {
		return c.String(http.StatusBadRequest, WarnMessageWhenInvalidID)
	}

	number, err := t.tripService.GetSoldTicketNumber(c.Request().Context(), param)
	if err != nil {
		return c.String(http.StatusInternalServerError, WarnInternalError)
	}

	return c.JSON(http.StatusOK, number)
}

func (t *handler) GetTotalRevenueForSpecificTrip(c echo.Context) error {

}
