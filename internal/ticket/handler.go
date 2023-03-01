package ticket

import "github.com/labstack/echo/v4"

type handler struct {

}

func NewHandler(e *echo.Echo) *handler {
	h := handler{}

	e.POST("/purchase/:id", h.Purchase)

	return &h
}

func (ti *handler) Purchase(c echo.Context) error {

}
