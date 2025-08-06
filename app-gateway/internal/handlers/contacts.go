package handlers

import "github.com/labstack/echo/v4"

// ContactsPage handles the request for the "Contacts" page.
func ContactsPage(c echo.Context) error {
	// Отправка HTML файла страницы "Контакты"
	return c.File("./internal/templates/contacts.html")
}
