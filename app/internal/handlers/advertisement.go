package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

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

// AdvertisementDelete handles the deletion of an advertisement.
func AdvertisementDelete(c echo.Context) error { // надо переделать на Delete
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

// ItemData представляет данные для страницы объявления
type ItemData struct {
	Title         string
	User          *UserData
	Advertisement *AdvertisementData
}

// AdvertisementData представляет данные для страницы объявления
type AdvertisementData struct {
	ID           int
	Title        string
	MainImage    string
	Images       []string
	Rating       float64
	ReviewsCount int
	Location     string
	Price        string
	PricePerDay  string
	Deposit      int
	Description  string
	Specs        []string
	RentalTerms  []string
	Owner        OwnerData
	Reviews      []ReviewData
}

// OwnerData представляет данные владельца
type OwnerData struct {
	ID           string
	Name         string
	Avatar       string
	Rating       float64
	ReviewsCount int
	IsVerified   bool
}

// ReviewData представляет данные отзыва
type ReviewData struct {
	ID     string
	User   UserData
	Rating float64
	Text   string
	Date   string
}

// UserData представляет данные пользователя
type UserData struct {
	Name   string
	Avatar string
}

// AdvertisementPage handles the request to view a specific advertisement.
func AdvertisementPage(c echo.Context) error {
	// id := c.Param("id")

	// advertisement := &advertisements.Advertisement{ID: id}
	// advertisement, err := advertisement.ReadByID(c.Request().Context(), c.Get("pool").(*pgxpool.Pool))
	// if err != nil {
	// 	c.Get("logger").(*zap.Logger).Error("failed to get advertisement",
	// 		zap.String("id", id),
	// 		zap.Error(err),
	// 	)
	// 	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка получения объявления"})
	// }

	// return c.JSON(http.StatusOK, map[string]interface{}{"advertisement": advertisement})

	userId := c.Get("user_id").(string)
	user := users.User{ID: userId}
	err := user.ReadByID(c.Request().Context(), c.Get("pool").(*pgxpool.Pool))
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to get user",
			zap.String("user_id", userId),
			zap.Error(err),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка получения пользователя"})
	}

	data := AdvertisementData{
		Title:     "Камера Sony A7III с объективом 24-70mm f/2.8",
		MainImage: "/images/items/camera-main.jpg",
		Images: []string{
			"/images/items/camera-main.jpg",
			"/images/items/camera-side.jpg",
			"/images/items/camera-back.jpg",
			"/images/items/camera-lens.jpg",
			"/images/items/camera-case.jpg",
		},
		Rating:       4.8,
		ReviewsCount: 24,
		Location:     "Москва, Центральный округ",
		Price:        "1,500 ₽",
		Deposit:      50000,
		Description:  "Полнокадровая беззеркальная камера Sony A7III с матрицей 24.2 Мп, системой автофокусировки с 693 точками и 5-осевой стабилизацией изображения. В комплекте идет профессиональный объектив Sony FE 24-70mm f/2.8 GM. Идеально подходит для профессиональной фото- и видеосъемки. Камера в отличном состоянии, используется аккуратно.",
		Specs: []string{
			"Матрица: 24.2 Мп, полный кадр",
			"Стабилизация: 5-осевая",
			"Автофокус: 693 точки",
			"Видео: 4K 30p",
			"Батарея: NP-FZ100",
			"Объектив: Sony FE 24-70mm f/2.8 GM",
			"Вес: 650г (с объективом)",
			"Размеры: 126.9 x 95.6 x 73.7 мм",
		},
		RentalTerms: []string{
			"Минимальный срок аренды: 1 день",
			"Максимальный срок аренды: 30 дней",
			"Требуется залог: 50,000 ₽",
			"Страховка: включена в стоимость",
			"Доставка: возможна по Москве (+500 ₽)",
			"Возврат: в день окончания аренды",
		},
		Owner: OwnerData{
			ID:           "owner123",
			Name:         "Александр Петров",
			Avatar:       "/images/users/owner-avatar.jpg",
			Rating:       4.9,
			ReviewsCount: 156,
			IsVerified:   true,
		},
		Reviews: []ReviewData{
			{
				ID: "review1",
				User: UserData{
					Name:   "Мария Козлова",
					Avatar: "/images/users/user1.jpg",
				},
				Rating: 5.0,
				Text:   "Отличная камера, всё работает идеально. Владелец очень ответственный и приятный в общении. Камера была в идеальном состоянии, объектив чистый. Рекомендую!",
				Date:   "12 марта 2024",
			},
			{
				ID: "review2",
				User: UserData{
					Name:   "Дмитрий Соколов",
					Avatar: "/images/users/user2.jpg",
				},
				Rating: 4.5,
				Text:   "Хорошая камера, качество фото отличное. Владелец пунктуальный, встретил вовремя. Единственный минус - немного тяжелая для длительной съемки.",
				Date:   "8 марта 2024",
			},
			{
				ID: "review3",
				User: UserData{
					Name:   "Анна Иванова",
					Avatar: "/images/users/user3.jpg",
				},
				Rating: 5.0,
				Text:   "Профессиональная техника в отличном состоянии. Владелец подробно объяснил все настройки. Снимала свадьбу, результат превзошел ожидания!",
				Date:   "5 марта 2024",
			},
		},
	}

	err = c.Render(http.StatusOK, "advertisement", data)
	if err != nil {
		c.Get("logger").(*zap.Logger).Error("failed to render advertisement page",
			zap.Error(err),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка отображения объявления"})
	}

	return nil
}
