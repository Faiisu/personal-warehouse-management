package handlers

import "github.com/gofiber/fiber/v2"

// APIIndex lists available HTTP endpoints for quick discovery.
// @Summary      API index
// @Description  Lists available API endpoints and brief descriptions.
// @Tags         meta
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /api [get]
func APIIndex(c *fiber.Ctx) error {
	endpoints := []fiber.Map{
		{
			"method":      "GET",
			"path":        "/api/health",
			"description": "Health check endpoint",
		},
		{
			"method":      "POST",
			"path":        "/api/register",
			"description": "Register a new user",
		},
		{
			"method":      "POST",
			"path":        "/api/login",
			"description": "Login user with email and password",
		},
		{
			"method":      "POST",
			"path":        "/api/events",
			"description": "Create a new event",
		},
		{
			"method":      "GET",
			"path":        "/api/events",
			"description": "List events",
		},
	}

	return c.JSON(fiber.Map{
		"name":      "Event Blog API",
		"endpoints": endpoints,
	})
}
