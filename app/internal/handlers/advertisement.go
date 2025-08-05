package handlers

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Danila331/ShareHub/internal/models/advertisements"
	"github.com/Danila331/ShareHub/internal/models/users"
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

// AdvertisementEditPage handles the "Edit Advertisement" page request.
func AdvertisementEditPage(c echo.Context) error {
	id := c.Param("id")

	// Получаем информацию об объявлении из базы данных
	var advertisement advertisements.Advertisement
	advertisement.ID = id
	err := advertisement.ReadByID(c.Request().Context(), c.Get("pool").(*pgxpool.Pool))
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to retrieve advertisement data",
			zap.String("id", id),
			zap.Error(err),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка получения объявления"})
	}

	// Проверяем, что пользователь является владельцем объявления
	if advertisement.UserID != c.Get("user_id").(string) {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "У вас нет прав для редактирования этого объявления"})
	}

	// Отправляем данные в шаблон для отображения
	err = c.Render(http.StatusOK, "edit-advertisement", advertisement)

	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to render edit advertisement page", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка отображения страницы редактирования объявления"})
	}

	return nil
}

// AdvertisementEditPost handles the form submission for editing an existing advertisement.
func AdvertisementEditPost(c echo.Context) error {
	advertisementID := c.Param("id")
	logger := c.Get("logger").(*zap.Logger)
	pool := c.Get("pool").(*pgxpool.Pool)

	// Получаем текстовые данные из формы
	var advertisement advertisements.Advertisement
	advertisement.ID = advertisementID
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
		logger.Error("missing required fields in advertisement form",
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

	// Получаем текущие фотографии объявления
	err := advertisement.ReadByID(c.Request().Context(), pool)
	if err != nil {
		logger.Error("failed to read advertisement", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка чтения объявления"})
	}

	// Получаем список удаленных серверных фотографий
	deletedPhotosJSON := c.FormValue("deleted_photos")
	var deletedPhotos []string
	if deletedPhotosJSON != "" {
		// Здесь нужно добавить импорт "encoding/json"
		// import "encoding/json"
		// err = json.Unmarshal([]byte(deletedPhotosJSON), &deletedPhotos)
		// if err != nil {
		//     logger.Error("failed to parse deleted photos", zap.Error(err))
		//     return c.JSON(http.StatusBadRequest, map[string]string{"error": "Ошибка обработки удаленных фотографий"})
		// }
		// Пока используем простой подход - разбиваем по запятой
		if deletedPhotosJSON != "[]" {
			// Убираем квадратные скобки и разбиваем по запятой
			cleaned := deletedPhotosJSON[1 : len(deletedPhotosJSON)-1] // убираем []
			if cleaned != "" {
				// Простое разбиение по запятой (для простоты)
				// В реальном проекте лучше использовать json.Unmarshal
				parts := strings.Split(cleaned, ",")
				for _, part := range parts {
					part = strings.Trim(part, `" `) // убираем кавычки и пробелы
					if part != "" {
						deletedPhotos = append(deletedPhotos, part)
					}
				}
			}
		}
	}

	// Удаляем фотографии из файловой системы
	for _, deletedPhoto := range deletedPhotos {
		photoPath := filepath.Join("./internal/static/images/uploads", deletedPhoto)
		err := os.Remove(photoPath)
		if err != nil && !os.IsNotExist(err) {
			logger.Warn("failed to delete photo file",
				zap.String("path", photoPath),
				zap.Error(err))
		} else {
			logger.Info("deleted photo file", zap.String("path", photoPath))
		}
	}

	// Обработка новых загруженных фотографий
	form, err := c.MultipartForm()
	if err != nil {
		logger.Error("failed to get multipart form", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка обработки формы"})
	}

	files := form.File
	var newPhotoPaths []string
	newPhotoPathsSet := make(map[string]bool) // Для предотвращения дублирования

	// Обрабатываем новые фотографии (photo_*)
	for key, fileHeaders := range files {
		if len(fileHeaders) > 0 && strings.HasPrefix(key, "photo_") {
			for _, fileHeader := range fileHeaders {
				// Генерируем уникальное имя файла
				ext := filepath.Ext(fileHeader.Filename)
				newFilename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), generateRandomString(8), ext)

				// Сохраняем файл
				dst := filepath.Join("./internal/static/images/uploads", newFilename)
				filePath := filepath.Join("uploads", newFilename)

				// Проверяем, что файл еще не добавлен
				if !newPhotoPathsSet[filePath] {
					newPhotoPaths = append(newPhotoPaths, filePath)
					newPhotoPathsSet[filePath] = true
				}

				src, err := fileHeader.Open()
				if err != nil {
					logger.Error("failed to open uploaded file",
						zap.String("filename", fileHeader.Filename),
						zap.Error(err),
					)
					return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка открытия файла"})
				}
				defer src.Close()

				out, err := os.Create(dst)
				if err != nil {
					logger.Error("failed to create destination file",
						zap.String("destination", dst),
						zap.Error(err),
					)
					return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка сохранения файла"})
				}
				defer out.Close()

				_, err = io.Copy(out, src)
				if err != nil {
					logger.Error("failed to copy file",
						zap.String("filename", fileHeader.Filename),
						zap.Error(err),
					)
					return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка копирования файла"})
				}

				logger.Info("saved new photo",
					zap.String("original_name", fileHeader.Filename),
					zap.String("saved_as", newFilename))
			}
		}
	}

	// Обрабатываем существующие фотографии (existing_photo_*)
	var existingPhotos []string
	existingPhotosSet := make(map[string]bool) // Для предотвращения дублирования
	for key, values := range form.Value {
		if strings.HasPrefix(key, "existing_photo_") {
			for _, value := range values {
				// Проверяем, что фотография не была удалена
				shouldKeep := true
				for _, deletedPhoto := range deletedPhotos {
					if value == deletedPhoto {
						shouldKeep = false
						break
					}
				}
				if shouldKeep && !existingPhotosSet[value] {
					existingPhotos = append(existingPhotos, value)
					existingPhotosSet[value] = true
				}
			}
		}
	}

	// Объединяем существующие и новые фотографии
	advertisement.PhotoPaths = append(existingPhotos, newPhotoPaths...)

	logger.Info("updating advertisement photos",
		zap.String("advertisement_id", advertisementID),
		zap.Int("existing_photos", len(existingPhotos)),
		zap.Int("new_photos", len(newPhotoPaths)),
		zap.Int("deleted_photos", len(deletedPhotos)),
		zap.Int("total_photos", len(advertisement.PhotoPaths)),
		zap.Strings("existing_photo_paths", existingPhotos),
		zap.Strings("new_photo_paths", newPhotoPaths),
		zap.Strings("deleted_photo_paths", deletedPhotos))

	// Обновляем объявление в базе данных
	err = advertisement.Update(c.Request().Context(), pool)
	if err != nil {
		logger.Error("failed to update advertisement", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка обновления объявления"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"success": true})
}

// AdvertisementDelete handles the deletion of an advertisement.
func AdvertisementDelete(c echo.Context) error {
	id := c.Param("id")

	advertisement := &advertisements.Advertisement{ID: id}
	err := advertisement.Delete(c.Request().Context(), c.Get("pool").(*pgxpool.Pool))
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to delete advertisement",
			zap.String("id", id),
			zap.Error(err),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка удаления объявления"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"success": true})
}

type AdvertisementPageData struct {
	User          users.User
	Advertisement advertisements.Advertisement
	IsOwner       bool
}

// AdvertisementPage handles the request to view a specific advertisement.
func AdvertisementPage(c echo.Context) error {
	advertisementId := c.Param("id")
	userId := c.Get("user_id").(string)

	// Получаем информацию об объявлении from the database
	var advertisement advertisements.Advertisement
	advertisement.ID = advertisementId
	err := advertisement.ReadByID(c.Request().Context(), c.Get("pool").(*pgxpool.Pool))
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to get advertisement",
			zap.String("id", advertisementId),
			zap.Error(err),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка получения объявления"})
	}

	var user users.User
	user.ID = userId
	err = user.ReadByID(c.Request().Context(), c.Get("pool").(*pgxpool.Pool))
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to get user",
			zap.String("id", userId),
			zap.Error(err),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка получения пользователя"})
	}

	data := AdvertisementPageData{
		User:          user,
		Advertisement: advertisement,
		IsOwner:       advertisement.UserID == user.ID,
	}

	err = c.Render(http.StatusOK, "advertisement", data)
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to render advertisement page",
			zap.String("id", advertisementId),
			zap.Error(err),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка отображения страницы объявления"})
	}

	return nil
}

// generateRandomString генерирует случайную строку заданной длины
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// В случае ошибки используем простую генерацию
			result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		} else {
			result[i] = charset[num.Int64()]
		}
	}
	return string(result)
}
