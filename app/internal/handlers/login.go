package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

// AuthPage handlers the "Auth" page request.
func LoginPage(c echo.Context) error {
	// Render the "Auth" page
	return c.File("./internal/templates/login.html")
}

// LoginPost handlers the "Auth" page post request.
func LoginPost(c echo.Context) error {
	var LoginRequest LoginRequest
	if err := c.Bind(&LoginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request",
		})
	}

	fmt.Println("LoginRequest: ", LoginRequest)
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Form data received successfully",
	})
}
