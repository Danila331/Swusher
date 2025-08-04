package servers

import (
	"html/template"
	"io"

	"github.com/Danila331/ShareHub/internal/handlers"
	"github.com/Danila331/ShareHub/internal/midlewary"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	"go.uber.org/zap"
)

// TemplateRenderer структура для рендера шаблонов
type TemplateRenderer struct {
	Templates *template.Template
}

// Render метод для рендера шаблонов
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

// StartServer создает и запускает Echo сервер с заданным логгером и пулом соединений с базой данных.
func StartServer(logger *zap.Logger, pool *pgxpool.Pool) {
	// Создаем новый Echo сервер
	app := echo.New()

	// Включаем статические файлы
	app.Static("/", "./internal/static")

	// Настройка рендера
	renderer := &TemplateRenderer{
		Templates: template.Must(template.ParseGlob("./internal/templates/*.html")),
	}

	app.Renderer = renderer
	// Устанавливаем обработчик для шаблонов
	app.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
			c.Response().Header().Set("Pragma", "no-cache")
			c.Response().Header().Set("Expires", "0")
			c.Response().Header().Set("Surrogate-Control", "no-store")
			c.Set("pool", pool)
			c.Set("logger", logger)
			return next(c)
		}
	})

	// Устанавливаем маршруты для страниц
	app.GET("/", handlers.MainPage)
	app.GET("/about", handlers.AboutPage)
	app.GET("/contacts", handlers.ContactsPage)
	app.GET("/how-it-works", handlers.HowItWorksPage)
	app.GET("/catalog", handlers.CatalogPage)

	// Register группа
	register := app.Group("/register")
	register.GET("/", handlers.RegisterPage, midlewary.IsAuthMiddleware)
	register.POST("/", handlers.RegisterPost)
	register.POST("/yandex", handlers.RegisterByYandexPost)

	// Login группа
	login := app.Group("/login")
	login.GET("/", handlers.LoginPage, midlewary.IsAuthMiddleware)
	login.POST("/", handlers.LoginPost)
	login.POST("/yandex", handlers.LoginByYandexPost)
	login.POST("/google", handlers.LoginByGooglePost)

	// Profile-saler группа
	profileSaler := app.Group("/profile-saler")
	profileSaler.GET("/", handlers.ProfileSalerPage, midlewary.JWTMiddleware)
	profileSaler.GET("/items", handlers.ProfileSalerItemsPage, midlewary.JWTMiddleware)
	profileSaler.GET("/rental-history", handlers.ProfileSalerRentalHistoryPage, midlewary.JWTMiddleware)
	profileSaler.GET("/reviews", handlers.ProfileSalerReviewsPage, midlewary.JWTMiddleware)
	profileSaler.GET("/earnings", handlers.ProfileSalerEarningsPage, midlewary.JWTMiddleware)
	profileSaler.GET("/settings", handlers.ProfileSalerSettingsPage, midlewary.JWTMiddleware)

	// Advertisement группа
	advertisement := app.Group("/advertisement")
	advertisement.GET("/:id", handlers.AdvertisementPage, midlewary.JWTMiddleware)
	advertisement.DELETE("/:id", handlers.AdvertisementDelete, midlewary.JWTMiddleware)
	advertisement.GET("/add", handlers.AdvertisementAddPage, midlewary.JWTMiddleware)
	advertisement.POST("/add", handlers.AdvertisementAddPost, midlewary.JWTMiddleware)

	// Запуск сервера на хосте
	app.Logger.Fatal(app.Start(":8081"))
}
