package servers

import (
	"github.com/labstack/echo/v4"
)

// StartServer initializes and starts the Echo server on port 8080.
func StartServer() {
	app := echo.New()
	app.Logger.Fatal(app.Start(":8080"))
}
