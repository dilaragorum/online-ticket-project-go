package handler

import (
	"errors"
	"github.com/dilaragorum/online-ticket-project-go/model"
	"github.com/dilaragorum/online-ticket-project-go/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	WarnAlreadyCreatedTrip            = "This trip is already created. Please create another trip."
	WarnInternalError                 = "Somethings go wrong. Please try later again"
	WarnMessageWhenThereAreEmptyBlank = "Please fill required area"
	WarnMessageWhenInvalidVehicle     = "Please enter valid Vehicle Type"
	WarnMessageWhenInvalidPrice       = "Please enter valid price"
)

type adminHandler struct {
	adminService service.AdminService
}

func NewAdminHandler(e *echo.Echo, adminService service.AdminService) *adminHandler {
	ah := adminHandler{adminService: adminService}

	e.POST("/admin/trips", ah.CreateTrip)
	return &ah
}

func (ah *adminHandler) CreateTrip(c echo.Context) error {
	trip := new(model.Trip)
	if err := c.Bind(&trip); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if trip.IsStartingPlaceEmpty() || trip.IsDestinationPlaceEmpty() || trip.IsDateEmpty() {
		return c.String(http.StatusBadRequest, WarnMessageWhenThereAreEmptyBlank)
	}

	if trip.IsInvalidVehicle() {
		return c.String(http.StatusBadRequest, WarnMessageWhenInvalidVehicle)
	}

	if trip.IsNotValidPrice() {
		return c.String(http.StatusBadRequest, WarnMessageWhenInvalidPrice)
	}

	requestCtx := c.Request().Context()

	if err := ah.adminService.CreateTrip(requestCtx, trip); err != nil {
		if errors.Is(err, service.ErrAlreadyCreatedTrip) {
			return c.String(http.StatusBadRequest, WarnAlreadyCreatedTrip)
		}
		return c.String(http.StatusInternalServerError, WarnInternalError)
	}

	return c.NoContent(http.StatusCreated)
}
