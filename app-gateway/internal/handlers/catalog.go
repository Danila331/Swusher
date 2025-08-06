package handlers

import "github.com/labstack/echo/v4"

// CatalogPage handles the request for the "Catalog" page.
func CatalogPage(c echo.Context) error {
	// Отправка HTML файла каталога
	return c.File("./internal/templates/catalog.html")
}
