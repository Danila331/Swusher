package handlers

import "github.com/labstack/echo/v4"

// HowItWorksPage handles the "How It Works" page request.
func HowItWorksPage(c echo.Context) error {
	// Render the "How It Works" page
	return c.File("./internal/templates/how-it-works.html")
}
