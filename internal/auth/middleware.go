package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

func TokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	jwtSecretKey := os.Getenv("ONLINE_TICKET_GO_JWTKEY")

	return func(c echo.Context) error {
		if c.Request().RequestURI == "/register" || c.Request().RequestURI == "/login" {
			return next(c)
		}

		cookie, err := c.Cookie("token")
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		token := cookie.Value

		claim := Claims{}
		parsedTokenInfo, err := jwt.ParseWithClaims(token, &claim, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecretKey), nil
		})
		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				return c.String(http.StatusUnauthorized, "Please login again")
			}

			return c.String(http.StatusUnauthorized, "Please login again")
		}

		if !parsedTokenInfo.Valid {
			return c.String(http.StatusForbidden, "Invalid token")
		}

		c.Set("claim", claim)

		return next(c)
	}
}

func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claim, _ := c.Get("claim").(Claims)

		if claim.IsNotAdmin() {
			return c.String(http.StatusForbidden, "You have no authority")
		}

		return next(c)
	}
}
