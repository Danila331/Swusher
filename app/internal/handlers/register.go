package handlers

import (
	"encoding/json"
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

// RegisterByYandexPost handlers the "Auth" page post request for Yandex.
func RegisterByYandexPost(c echo.Context) error {
	fmt.Println("Yandex OAuth register")
	var req YandexLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Некорректные данные"})
	}
	fmt.Printf("Yandex OAuth: %+v\n", req)

	// Запрос к Яндекс API для получения информации о пользователе
	client := &http.Client{}
	apiReq, err := http.NewRequest("GET", "https://login.yandex.ru/info", nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Ошибка создания запроса к Яндекс API"})
	}
	apiReq.Header.Set("Authorization", "OAuth "+req.AccessToken)

	resp, err := client.Do(apiReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Ошибка запроса к Яндекс API"})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Не удалось получить данные пользователя от Яндекса"})
	}

	var yandexUser map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&yandexUser); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Ошибка декодирования ответа Яндекса"})
	}

	fmt.Printf("Yandex user info: %+v\n", yandexUser)
	return c.JSON(http.StatusOK, yandexUser)
}
