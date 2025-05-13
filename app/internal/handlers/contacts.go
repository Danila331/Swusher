package handlers

import "github.com/labstack/echo/v4"

// ContactsPage handles the "Contacts" page request.
func ContactsPage(c echo.Context) error {
	// Render the "Contacts" page
	return c.File("./internal/templates/contacts.html")
}
