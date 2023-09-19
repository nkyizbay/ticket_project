package auth

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

var allowList = map[string]bool{
	"/register": true,
	"/login":    true,
}

func TokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	jwtSecretKey := viper.GetString("ONLINE_TICKET_GO_JWTKEY")

	return func(c echo.Context) error {
		if _, ok := allowList[c.Request().RequestURI]; ok {
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