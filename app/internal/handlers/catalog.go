package handlers

import "github.com/labstack/echo/v4"

// CatalogPage handles the "Catalog" page request.
func CatalogPage(c echo.Context) error {
	// Render the "Catalog" page
	return c.File("./internal/templates/catalog.html")
}
