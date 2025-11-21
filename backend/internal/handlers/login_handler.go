package handlers

import (
	"context"
	"errors"
	"strings"
	"time"

	"my-backend/internal/db"
	"my-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

// LoginUser godoc
// @Summary      Login user
// @Description  Authenticate user by email and password.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        payload  body      loginRequest  true  "Login credentials"
// @Success      200  {object}  models.Users
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/login [post]
func LoginUser(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON payload")
	}

	req.Email = strings.TrimSpace(req.Email)

	if req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "email and password are required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection, err := db.UsersCollection(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "database unavailable")
	}

	var user models.Users
	if err := collection.FindOne(ctx, bson.M{"Email": req.Email}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}

	return c.JSON(user)
}
