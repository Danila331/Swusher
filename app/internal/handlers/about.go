package handlers

import "github.com/labstack/echo/v4"

// AboutPage handles the "About" page request.
func AboutPage(c echo.Context) error {
	// Render the "About" page
	return c.File("./internal/templates/about.html")
}
