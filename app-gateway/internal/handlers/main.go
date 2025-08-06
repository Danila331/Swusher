package handlers

import (
	"github.com/labstack/echo/v4"
)

// MainPage handles the request for the main page.
func MainPage(c echo.Context) error {
	// Отправка HTML файла главной страницы
	return c.File("./internal/templates/index.html")
}
