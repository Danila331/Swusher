package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// RegisterRequest represents the structure of the registration form.
type RegisterRequest struct {
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

// RegisterPage handles the "Register" page request.
func RegisterPage(c echo.Context) error {
	// Render the "Register" page
	return c.File("./internal/templates/register.html")
}

// RegisterPost handles the form submission for the "Register" post reauest.
func RegisterPost(c echo.Context) error {
	var RegisterRequest RegisterRequest
	if err := c.Bind(&RegisterRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request",
		})
	}

	fmt.Println("RegisterRequest: ", RegisterRequest)
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Form data received successfully",
	})
}
