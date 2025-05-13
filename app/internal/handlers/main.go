package handlers

import (
	"github.com/labstack/echo/v4"
)

// MainPage handles the main page request.
func MainPage(c echo.Context) error {
	// Render the main page
	return c.File("./internal/templates/index.html")
}
