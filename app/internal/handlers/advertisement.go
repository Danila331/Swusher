package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Danila331/ShareHub/internal/models/advertisements"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// AdvertisementAddPage handles the "Add Advertisement" page request.
func AdvertisementAddPage(c echo.Context) error {
	err := c.Render(http.StatusOK, "add-advertisement", nil)
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to render advertisement page", zap.Error(err))
		return err
	}
	return nil
}

// AdvertisementAddPost handles the form submission for adding a new advertisement.
func AdvertisementAddPost(c echo.Context) error {
	// Получение текстовых данных из формы
	var advertisement advertisements.Advertisement
	advertisement.UserID = c.Get("user_id").(string) // Получаем ID пользователя из контекста
	advertisement.Title = c.FormValue("title")
	advertisement.Category = c.FormValue("category")
	advertisement.Description = c.FormValue("description")
	advertisement.CostPerday, _ = strconv.ParseFloat(c.FormValue("cost_per_day"), 64)
	advertisement.CostPerWeek, _ = strconv.ParseFloat(c.FormValue("cost_per_week"), 64)
	advertisement.CostPerMonth, _ = strconv.ParseFloat(c.FormValue("cost_per_month"), 64)
	advertisement.Address = c.FormValue("location")
	advertisement.Geolocation_X, _ = strconv.ParseFloat(c.FormValue("geolocation_x"), 64)
	advertisement.Geolocation_Y, _ = strconv.ParseFloat(c.FormValue("geolocation_y"), 64)
	advertisement.RentalRules = c.FormValue("rental_rules")

	// Проверка обязательных полей
	if advertisement.Title == "" || advertisement.Category == "" || advertisement.Description == "" || advertisement.CostPerday <= 0 || advertisement.Address == "" {
		c.Get("logger").(*zap.Logger).Error("missing required fields in advertisement form",
			zap.String("title", advertisement.Title),
			zap.String("category", advertisement.Category),
			zap.String("description", advertisement.Description),
			zap.Float64("cost_per_day", advertisement.CostPerday),
			zap.String("location", advertisement.Address),
			zap.Float64("geolocation_x", advertisement.Geolocation_X),
			zap.Float64("geolocation_y", advertisement.Geolocation_Y),
			zap.String("rental_rules", advertisement.RentalRules),
		)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Все обязательные поля должны быть заполнены"})
	}

	// Обработка загруженных фотографий
	form, err := c.MultipartForm()
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to get multipart form",
			zap.Error(err),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка обработки формы"})
	}

	files := form.File
	uploadedFiles := []string{}
	filesPaths := []string{}

	for key, fileHeaders := range files {
		if len(fileHeaders) > 0 && key[:6] == "photo_" { // Проверяем, что это фото
			for _, fileHeader := range fileHeaders {
				// Сохраняем файл
				dst := filepath.Join("./internal/static/images/uploads", fileHeader.Filename)
				filePath := filepath.Join("uploads", fileHeader.Filename)
				filesPaths = append(filesPaths, filePath)

				src, err := fileHeader.Open()
				if err != nil {
					c.Get("logger").(*zap.Logger).Error("failed to open uploaded file",
						zap.String("filename", fileHeader.Filename),
						zap.Error(err),
					)
					return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка открытия файла"})
				}
				defer src.Close()

				out, err := os.Create(dst)
				if err != nil {
					c.Get("logger").(*zap.Logger).Error("failed to create destination file",
						zap.String("destination", dst),
					)
					return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка сохранения файла"})
				}
				defer out.Close()

				_, err = io.Copy(out, src)
				if err != nil {
					c.Get("logger").(*zap.Logger).Error("failed to copy file",
						zap.String("filename", fileHeader.Filename),
						zap.Error(err),
					)
					return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка копирования файла"})
				}

				uploadedFiles = append(uploadedFiles, dst)
			}
		}
	}

	advertisement.PhotoPaths = filesPaths
	err = advertisement.Create(c.Request().Context(), c.Get("pool").(*pgxpool.Pool))
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to add new advertisement",
			zap.Error(err),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка добавления объявления"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"success": true})
}
