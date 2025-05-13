package servers

import (
	"github.com/Danila331/ShareHub/internal/handlers"
	"github.com/labstack/echo/v4"

	"go.uber.org/zap"
)

// StartServer initializes and starts the Echo server on port 8080.
func StartServer(logger *zap.Logger) {
	// Start echo web app
	app := echo.New()

	// Include static files
	app.Static("/", "./internal/static")

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

	// Login group
	login := app.Group("/login")
	login.GET("/", handlers.LoginPage)
	login.POST("/", handlers.LoginPost)
	// Start server on host
	app.Logger.Fatal(app.Start(":8080"))
}
