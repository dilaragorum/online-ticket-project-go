package handler

import (
	"errors"
	"github.com/dilaragorum/online-ticket-project-go/model"
	"github.com/dilaragorum/online-ticket-project-go/service"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

const (
	WarnInternalError = "Somethings go wrong. Please try later again"

	WarnAlreadyCreatedTrip            = "This trip is already created. Please create another trip."
	WarnMessageWhenThereAreEmptyBlank = "Please fill required area"
	WarnMessageWhenInvalidVehicle     = "Please enter valid Vehicle Type"
	WarnMessageWhenInvalidPrice       = "Please enter valid price"

	WarnMessageWhenInvalidID             = "Please enter valid ID"
	WarnMessageWhenTripNotExistForDelete = "This trip does not exist or it is deleted already. "
)

type adminHandler struct {
	adminService service.AdminService
	jwtSecretKey string
}

func (ah *adminHandler) adminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		token := cookie.Value

		claim := model.Claims{}
		parsedTokenInfo, err := jwt.ParseWithClaims(token, &claim, func(token *jwt.Token) (interface{}, error) {
			return []byte(ah.jwtSecretKey), nil
		})
		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				return c.String(http.StatusUnauthorized, "Please login again")
			}

			return c.String(http.StatusForbidden, "Please login again")
		}

		if !parsedTokenInfo.Valid {
			return c.String(http.StatusUnauthorized, "Invalid token")
		}

		if !claim.IsAdmin() {
			return c.String(http.StatusForbidden, "You have no authority")
		}

		return next(c)
	}
}

func NewAdminHandler(e *echo.Echo, adminService service.AdminService, jwtSecretKey string) *adminHandler {
	ah := adminHandler{adminService: adminService, jwtSecretKey: jwtSecretKey}

	admin := e.Group("/admin", ah.adminMiddleware)

	admin.POST("/trips", ah.CreateTrip)
	admin.DELETE("/trips/:id", ah.CancelTrip)

	return &ah
}

func (ah *adminHandler) CreateTrip(c echo.Context) error {
	trip := new(model.Trip)
	if err := c.Bind(&trip); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	cookie, err := c.Cookie("token")
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	claims := model.Claims{}
	jwt.ParseWithClaims(cookie.Value, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(ah.jwtSecretKey), nil
	})

	if claims.AuthorizationType != model.AuthAdmin {
		return c.String(http.StatusForbidden, err.Error())
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

	if err := ah.adminService.CreateTrip(requestCtx, trip); err != nil {
		if errors.Is(err, service.ErrAlreadyCreatedTrip) {
			return c.String(http.StatusBadRequest, WarnAlreadyCreatedTrip)
		}
		return c.String(http.StatusInternalServerError, WarnInternalError)
	}

	return c.NoContent(http.StatusCreated)
}

func (ah *adminHandler) CancelTrip(c echo.Context) error {
	tripIDStr := c.Param("id")
	tripID, _ := strconv.Atoi(tripIDStr)

	if model.IsInvalidID(tripID) {
		return c.String(http.StatusBadRequest, WarnMessageWhenInvalidID)
	}

	requestCtx := c.Request().Context()

	if err := ah.adminService.CancelTrip(requestCtx, tripID); err != nil {
		switch {
		case errors.Is(err, service.ErrTripNotExist):
			return c.String(http.StatusBadRequest, WarnMessageWhenTripNotExistForDelete)
		}
		return c.String(http.StatusInternalServerError, WarnInternalError)
	}

	return c.NoContent(http.StatusNoContent)
}
