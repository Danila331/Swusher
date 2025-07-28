package handlers

import "github.com/labstack/echo/v4"

// AboutPage handles the request for the "About" page
func AboutPage(c echo.Context) error {
	// Отправка HTML файла о странице "О нас"
	return c.File("./internal/templates/about.html")
}
