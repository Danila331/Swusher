package servers

import (
	"github.com/Danila331/ShareHub/internal/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	"go.uber.org/zap"
)

// StartServer создает и запускает Echo сервер с заданным логгером и пулом соединений с базой данных.
func StartServer(logger *zap.Logger, pool *pgxpool.Pool) {
	// Создаем новый Echo сервер
	app := echo.New()

	// Включаем статические файлы
	app.Static("/", "./internal/static")

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
	register.GET("/", handlers.RegisterPage)
	register.POST("/", handlers.RegisterPost)
	register.POST("/yandex", handlers.RegisterByYandexPost)

	// Login группа
	login := app.Group("/login")
	login.GET("/", handlers.LoginPage)
	login.POST("/", handlers.LoginPost)
	login.POST("/yandex", handlers.LoginByYandexPost)
	login.POST("/google", handlers.LoginByGooglePost)

	// Запуск сервера на хосте
	app.Logger.Fatal(app.Start(":8080"))
}
