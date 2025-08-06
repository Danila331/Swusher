package midlewary

import (
	"net/http"

	"github.com/Danila331/Swusher/pkg/jwt"
	"github.com/labstack/echo/v4"
)

func IsAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil || cookie.Value == "" {
			return next(c)
		}

		userID, err := jwt.ParseToken(cookie.Value)
		if err != nil {
			return next(c)
		}

		c.Set("user_id", userID)
		c.Redirect(http.StatusSeeOther, "/profile-saler/")
		return next(c)
	}
}
