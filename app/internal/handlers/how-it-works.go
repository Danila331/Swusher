package handlers

import "github.com/labstack/echo/v4"

// HowItWorksPage handles the request for the "How It Works" page.
func HowItWorksPage(c echo.Context) error {
	// Отправка HTML файла страницы "Как это работает"
	return c.File("./internal/templates/how-it-works.html")
}
