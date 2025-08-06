package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Danila331/Swusher/internal/models/users"
	"github.com/Danila331/Swusher/pkg/hash"
	"github.com/Danila331/Swusher/pkg/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// LoginRequest структура запроса для логина.
type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

// YandexLoginRequest структура запроса для логина через Яндекс.
type YandexLoginRequest struct {
	AccessToken string `json:"access_token"`
}

// GoogleLoginRequest структура запроса для логина через Google.
type GoogleLoginRequest struct {
	Credential string `json:"credential"`
}

// LoginPage handles the request for the "Auth" page.
func LoginPage(c echo.Context) error {
	// Отправка HTML файла страницы "Авторизация"
	return c.File("./internal/templates/login.html")
}

// LoginPost handles the "Auth" page post request.
func LoginPost(c echo.Context) error {
	var loginRequest LoginRequest
	if err := c.Bind(&loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request",
		})
	}

	user := users.User{
		Email: loginRequest.Email,
	}

	// Проверяем, существует ли пользователь с таким email
	// Если нет, возвращаем ошибку
	err := user.ReadByEmail(c.Request().Context(), c.Get("pool").(*pgxpool.Pool))
	fmt.Println("Login attempt for user:", user.Email, user.Password)
	if err != nil || !hash.CheckPasswordHash(loginRequest.Password, user.Password) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Неверный email или пароль"})
	}

	// Генерируем JWT токен
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	// Кладём токен в httpOnly cookie
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = false // true для https
	cookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Login successful",
	})
}

// LoginByYandexPost handlers the "Auth" page post request for Yandex.
func LoginByYandexPost(c echo.Context) error {
	var req YandexLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Некорректные данные"})
	}

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

	// Получаем email и имя пользователя
	email, ok := yandexUser["default_email"].(string)
	if !ok || email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Не удалось получить email пользователя от Яндекса"})
	}

	// Проверяем, есть ли пользователь с таким email
	user := users.User{Email: email}
	err = user.ReadByEmail(c.Request().Context(), c.Get("pool").(*pgxpool.Pool))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Пользователь не найден"})
	}

	// Генерируем JWT токен
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Ошибка генерации токена"})
	}

	// Кладём токен в httpOnly cookie
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = false // true для https
	cookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Yandex OAuth успешен",
	})
}

// LoginByGooglePost handlers the "Auth" page post request for Google.
func LoginByGooglePost(c echo.Context) error {
	var req GoogleLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Некорректные данные"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message":    "Google OAuth успешен",
		"credential": req.Credential,
	})
}
