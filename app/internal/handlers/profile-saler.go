package handlers

import (
	"net/http"
	"strconv"

	"github.com/Danila331/ShareHub/internal/models/advertisements"
	"github.com/Danila331/ShareHub/internal/models/users"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// ProfileSaler struct of data for ProfileSalerPage
type ProfileSaler struct {
	Name         string
	Surname      string
	Email        string
	Photo        string
	Rating       float64
	ReviewsCount int
	IsVerified   bool
}

// ProfileSalerItemsPage struct of data for ProfileSalerItemsPage
type ProfileSalerItems struct {
	ProfileSaler
	Advertisements []advertisements.Advertisement
}

// ProfileSalerPage handles the "Profile Saler" page request.
func ProfileSalerPage(c echo.Context) error {
	// Получаем данные пользователя из базы данных
	var user users.User
	user.ID = c.Get("user_id").(string)
	err := user.ReadByID(c.Request().Context(), c.Get("pool").(*pgxpool.Pool))
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to retrieve user data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user data"})
	}

	// Подготовка данных для отображения на странице
	profileSaler := ProfileSaler{
		Name:         user.Name,
		Surname:      user.LastName,
		Email:        user.Email,
		Photo:        user.PhotoPath,
		Rating:       4.5, // Здесь должна быть логика для получения рейтинга
		ReviewsCount: 10,  // Здесь должна быть логика для получения количества отзыв
		IsVerified:   user.IsVerified,
	}

	// Отправляем данные в шаблон для отображения
	err = c.Render(http.StatusOK, "profile-saler", profileSaler)
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to render profile saler page", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to render profile saler page"})
	}
	return nil
}

// ProfileSalerItemsPage handles the "Profile Saler Items" page request.
func ProfileSalerItemsPage(c echo.Context) error {
	// Получаем данные пользователя из базы данных
	var user users.User
	user.ID = c.Get("user_id").(string)
	err := user.ReadByID(c.Request().Context(), c.Get("pool").(*pgxpool.Pool))
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to retrieve user data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user data"})
	}

	// Получаем параметры пагинации из запроса
	limitParam := c.QueryParam("limit")
	offsetParam := c.QueryParam("offset")

	// Устанавливаем значения по умолчанию для limit и offset
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10 // Значение по умолчанию
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0 // Значение по умолчанию
	}

	// Получаем объявления пользователя из базы данных
	var advertisement advertisements.Advertisement
	advertisement.UserID = c.Get("user_id").(string)
	ads, err := advertisement.ReadAllByUserID(c.Request().Context(), c.Get("pool").(*pgxpool.Pool), limit, offset)
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to retrieve advertisements", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve advertisements"})
	}

	// Подготовка данных для отображения на странице
	profileSalerItems := ProfileSalerItems{
		ProfileSaler: ProfileSaler{
			Name:         user.Name,
			Surname:      user.LastName,
			Email:        user.Email,
			Photo:        user.PhotoPath,
			Rating:       4.5, // Здесь должна быть логика для получения рейтинга
			ReviewsCount: 10,  // Здесь должна быть логика для получения количества отзывов
			IsVerified:   user.IsVerified,
		},
		Advertisements: ads,
	}

	// Отправляем данные в шаблон для отображения
	err = c.Render(http.StatusOK, "my-items", profileSalerItems)
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to render profile saler items page", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to render profile saler items page"})
	}
	return nil
}

// ProfileSalerRentalHistoryPage handles the "Profile Saler Rental History" page request.
func ProfileSalerRentalHistoryPage(c echo.Context) error {
	// Получаем данные пользователя из базы данных
	var user users.User
	user.ID = c.Get("user_id").(string)
	err := user.ReadByID(c.Request().Context(), c.Get("pool").(*pgxpool.Pool))
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to retrieve user data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user data"})
	}

	// Подготовка данных для отображения на странице
	profileSaler := ProfileSaler{
		Name:         user.Name,
		Surname:      user.LastName,
		Email:        user.Email,
		Photo:        user.PhotoPath,
		Rating:       4.5, // Здесь должна быть логика для получения рейтинга
		ReviewsCount: 10,  // Здесь должна быть логика для получения количества отзыв
		IsVerified:   user.IsVerified,
	}

	// Отправляем данные в шаблон для отображения
	err = c.Render(http.StatusOK, "rental-history", profileSaler)
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to render profile saler page", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to render profile saler page"})
	}
	return nil
}

// ProfileSalerReviewsPage handles the "Profile Saler Reviews" page request.
func ProfileSalerReviewsPage(c echo.Context) error {
	// Получаем данные пользователя из базы данных
	var user users.User
	user.ID = c.Get("user_id").(string)
	err := user.ReadByID(c.Request().Context(), c.Get("pool").(*pgxpool.Pool))
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to retrieve user data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user data"})
	}

	// Подготовка данных для отображения на странице
	profileSaler := ProfileSaler{
		Name:         user.Name,
		Surname:      user.LastName,
		Email:        user.Email,
		Photo:        user.PhotoPath,
		Rating:       4.5, // Здесь должна быть логика для получения рейтинга
		ReviewsCount: 10,  // Здесь должна быть логика для получения количества отзыв
		IsVerified:   user.IsVerified,
	}

	// Отправляем данные в шаблон для отображения
	err = c.Render(http.StatusOK, "reviews", profileSaler)
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to render profile saler page", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to render profile saler page"})
	}
	return nil
}

// ProfileSalerEarningsPage handles the "Profile Saler Earnings" page request.
func ProfileSalerEarningsPage(c echo.Context) error {
	// Получаем данные пользователя из базы данных
	var user users.User
	user.ID = c.Get("user_id").(string)
	err := user.ReadByID(c.Request().Context(), c.Get("pool").(*pgxpool.Pool))
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to retrieve user data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user data"})
	}

	// Подготовка данных для отображения на странице
	profileSaler := ProfileSaler{
		Name:         user.Name,
		Surname:      user.LastName,
		Email:        user.Email,
		Photo:        user.PhotoPath,
		Rating:       4.5, // Здесь должна быть логика для получения рейтинга
		ReviewsCount: 10,  // Здесь должна быть логика для получения количества отзыв
		IsVerified:   user.IsVerified,
	}

	// Отправляем данные в шаблон для отображения
	err = c.Render(http.StatusOK, "earnings", profileSaler)
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to render profile saler page", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to render profile saler page"})
	}
	return nil
}

// ProfileSalerSettingsPage handles the "Profile Saler Settings" page request.
func ProfileSalerSettingsPage(c echo.Context) error {
	// Получаем данные пользователя из базы данных
	var user users.User
	user.ID = c.Get("user_id").(string)
	err := user.ReadByID(c.Request().Context(), c.Get("pool").(*pgxpool.Pool))
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to retrieve user data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user data"})
	}

	// Подготовка данных для отображения на странице
	profileSaler := ProfileSaler{
		Name:         user.Name,
		Surname:      user.LastName,
		Email:        user.Email,
		Photo:        user.PhotoPath,
		Rating:       4.5, // Здесь должна быть логика для получения рейтинга
		ReviewsCount: 10,  // Здесь должна быть логика для получения количества отзыв
		IsVerified:   user.IsVerified,
	}

	// Отправляем данные в шаблон для отображения
	err = c.Render(http.StatusOK, "settings", profileSaler)
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to render profile saler page", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to render profile saler page"})
	}
	return nil
}
