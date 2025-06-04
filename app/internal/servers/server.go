package servers

import (
	"github.com/Danila331/ShareHub/internal/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	"go.uber.org/zap"
)

// StartServer initializes and starts the Echo server on port 8080.
func StartServer(logger *zap.Logger, pool *pgxpool.Pool) {
	// Start echo web app
	app := echo.New()

	// Include static files
	app.Static("/", "./internal/static")

	// WE need delete in future
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

	// Handlers roots
	app.GET("/", handlers.MainPage)
	app.GET("/about", handlers.AboutPage)
	app.GET("/contacts", handlers.ContactsPage)
	app.GET("/how-it-works", handlers.HowItWorksPage)
	app.GET("/catalog", handlers.CatalogPage)

	// Register group
	register := app.Group("/register")
	register.GET("/", handlers.RegisterPage)
	register.POST("/", handlers.RegisterPost)
	register.POST("/yandex", handlers.RegisterByYandexPost)

	// Login group
	login := app.Group("/login")
	login.GET("/", handlers.LoginPage)
	login.POST("/", handlers.LoginPost)
	login.POST("/yandex", handlers.LoginByYandexPost)
	login.POST("/google", handlers.LoginByGooglePost)

	// Start server on host
	app.Logger.Fatal(app.Start(":8080"))
}
