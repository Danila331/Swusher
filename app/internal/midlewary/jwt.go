package midlewary

import (
	"net/http"

	"github.com/Danila331/ShareHub/pkg/jwt"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware проверяет JWT-токен из cookie "token" и кладёт user_id в context
func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil || cookie.Value == "" {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		userID, err := jwt.ParseToken(cookie.Value)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		c.Set("user_id", userID)
		return next(c)
	}
}
