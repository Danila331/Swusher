package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Danila331/Swusher/internal/models/users"
	"github.com/Danila331/Swusher/pkg/hash"
	"github.com/Danila331/Swusher/pkg/jwt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// RegisterRequest структура запроса для регистрации.
type RegisterRequest struct {
	Name     string `json:"name" form:"name"`
	LastName string `json:"lastname" form:"lastname"`
	Email    string `json:"email" form:"email"`
	Phone    string `json:"phone" form:"phone"`
	Password string `json:"password" form:"password"`
}

// RegisterPage handles the "Register" page request.
func RegisterPage(c echo.Context) error {
	// Отправка HTML файла страницы "Регистрация"
	return c.File("./internal/templates/register.html")
}

// RegisterPost handles the form submission for the "Register" post request.
func RegisterPost(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// TODO: добавить валидацию req.Email, req.Phone, req.Password

	hashPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to hash password", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}

	// Создаем нового пользователя
	// Если пользователь с таким email уже существует, возвращаем ошибку
	var user users.User
	user.Name = req.Name
	user.LastName = req.LastName
	user.Email = req.Email
	user.Phone = req.Phone
	user.Password = hashPassword
	user.Role = "user"
	user.IsVerified = false
	user.PhotoPath = "default.png"

	if user.Nickname == "" {
		user.Nickname = fmt.Sprintf("%s_%s_%s", user.Name, user.LastName, uuid.New().String()[:8])
	}

	err = user.Create(c.Request().Context(), c.Get("pool").(*pgxpool.Pool))
	if err != nil {
		if strings.Contains(err.Error(), "уже существует") {
			return c.JSON(http.StatusConflict, map[string]string{"error": "User already exists"})
		}
		c.Get("logger").(*zap.Logger).Error("failed to create user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
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

	c.Get("logger").(*zap.Logger).Info("user registered", zap.String("email", user.Email))
	return c.JSON(http.StatusCreated, map[string]string{"message": "User registered successfully"})
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

	var user users.User
	user.Email, _ = yandexUser["default_email"].(string)
	if user.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Не удалось получить email пользователя от Яндекса"})
	}

	// Добавляем пользователя в базу данных
	user.Name, _ = yandexUser["first_name"].(string)
	user.LastName, _ = yandexUser["last_name"].(string)
	user.Phone = ""    // Яндекс не предоставляет номер телефона
	user.Password = "" // Пароль не нужен, т.к. вход через OAuth
	user.Role = "user"
	user.IsVerified = false
	user.PhotoPath = "default.png"

	if user.Nickname == "" {
		user.Nickname = fmt.Sprintf("%s_%s_%s", user.Name, user.LastName, uuid.New().String()[:8])
	}

	err = user.Create(c.Request().Context(), c.Get("pool").(*pgxpool.Pool))
	if err != nil {
		if strings.Contains(err.Error(), "уже существует") {
			return c.JSON(http.StatusConflict, map[string]string{"message": "Пользователь с таким email уже существует"})
		}
		c.Get("logger").(*zap.Logger).Error("failed to create user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Ошибка создания пользователя"})
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

	c.Get("logger").(*zap.Logger).Info("user registered via Yandex OAuth", zap.String("email", user.Email))
	return c.JSON(http.StatusOK, yandexUser)
}
