package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// LoginRequest represents the structure of the login form.
type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

// YandexLoginRequest represents the structure of the Yandex login request.
type YandexLoginRequest struct {
	AccessToken string `json:"access_token"`
}

// GoogleLoginRequest represents the structure of the Google login request.
type GoogleLoginRequest struct {
	Credential string `json:"credential"`
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

// LoginByYandexPost handlers the "Auth" page post request for Yandex.
func LoginByYandexPost(c echo.Context) error {
	fmt.Println("Yandex OAuth login")
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

// LoginByGooglePost handlers the "Auth" page post request for Google.
func LoginByGooglePost(c echo.Context) error {
	var req GoogleLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Некорректные данные"})
	}
	fmt.Printf("Google OAuth: %+v\n", req)
	// Здесь можно добавить валидацию credential через Google API, если нужно
	return c.JSON(http.StatusOK, map[string]string{
		"message":    "Google OAuth успешен",
		"credential": req.Credential,
	})
}
