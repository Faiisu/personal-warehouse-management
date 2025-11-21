package routes

import (
	"my-backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	app.Get("/api", handlers.APIIndex)
	app.Get("/api/health", handlers.HealthCheck)
	app.Post("/api/register", handlers.RegisterUser)
	app.Post("/api/login", handlers.LoginUser)
}
