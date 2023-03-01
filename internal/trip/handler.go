package trip

import (
	"errors"
	"github.com/dilaragorum/online-ticket-project-go/internal/user"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	WarnNoTripMeetConditions = "There is no trip which meet your conditions."
)

type handler struct {
	tripService Service
}

func Handler(e *echo.Echo, tripService Service) *handler {
	h := handler{
		tripService: tripService,
	}

	e.GET("/trips", h.FilterTrips)

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
