package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"my-backend/internal/db"
	"my-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type categoryRequest struct {
	StockID      string `json:"StockID"`
	CategoryName string `json:"CategoryName"`
	Discription  string `json:"Discription"`
}

// ListCategories godoc
// @Summary      List categories by stock
// @Description  Returns categories filtered by StockID.
// @Tags         categories
// @Produce      json
// @Param        stockId  query  string  true  "Stock ID (UUID)"
// @Success      200  {array}   models.Categories
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/categories [get]
func ListCategories(c *fiber.Ctx) error {
	stockIDParam := strings.TrimSpace(c.Query("stockId"))
	if stockIDParam == "" {
		return fiber.NewError(fiber.StatusBadRequest, "stockId is required")
	}

	stockUUID, err := uuid.Parse(stockIDParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "stockId must be a valid UUID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection, err := db.CategoriesCollection(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "database unavailable")
	}

	cursor, err := collection.Find(ctx, bson.M{"StockID": stockUUID})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch categories")
	}
	defer cursor.Close(ctx)

	var categories []models.Categories
	if err := cursor.All(ctx, &categories); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to decode categories")
	}

	return c.JSON(categories)
}

// CreateCategories godoc
// @Summary      Bulk create categories
// @Description  Creates multiple categories in a single request.
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        payload  body      []categoryRequest  true  "List of categories"
// @Success      201  {array}   models.Categories
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/categories [post]
func CreateCategories(c *fiber.Ctx) error {
	var payload []categoryRequest
	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON payload")
	}

	if len(payload) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "at least one category is required")
	}

	categories := make([]models.Categories, len(payload))
	docs := make([]interface{}, len(payload))

	for i, cat := range payload {
		cat.StockID = strings.TrimSpace(cat.StockID)
		cat.CategoryName = strings.TrimSpace(cat.CategoryName)
		cat.Discription = strings.TrimSpace(cat.Discription)

		if cat.StockID == "" || cat.CategoryName == "" {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("StockID and CategoryName are required at index %d", i))
		}

		stockUUID, err := uuid.Parse(cat.StockID)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("StockID at index %d must be a valid UUID", i))
		}

		category := models.Categories{
			CategoryID:   uuid.New(),
			StockID:      stockUUID,
			CategoryName: cat.CategoryName,
			Discription:  cat.Discription,
		}
		categories[i] = category
		docs[i] = category
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection, err := db.CategoriesCollection(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "database unavailable")
	}

	if _, err := collection.InsertMany(ctx, docs); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create categories")
	}

	return c.Status(fiber.StatusCreated).JSON(categories)
}

// DeleteCategory godoc
// @Summary      Delete a category
// @Description  Sets Category to null on related products (matching StockID and category name) then deletes the category document.
// @Tags         categories
// @Produce      json
// @Param        categoryId  path  string  true  "Category ID (UUID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/categories/{categoryId} [delete]
func DeleteCategory(c *fiber.Ctx) error {
	categoryIDParam := strings.TrimSpace(c.Params("categoryId"))
	if categoryIDParam == "" {
		return fiber.NewError(fiber.StatusBadRequest, "categoryId is required")
	}

	categoryUUID, err := uuid.Parse(categoryIDParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "categoryId must be a valid UUID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	categoriesCol, err := db.CategoriesCollection(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "database unavailable")
	}
	productsCol, err := db.ProductsCollection(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "database unavailable")
	}

	var category models.Categories
	if err := categoriesCol.FindOne(ctx, bson.M{"CategoryID": categoryUUID}).Decode(&category); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fiber.NewError(fiber.StatusNotFound, "category not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch category")
	}

	updateRes, err := productsCol.UpdateMany(
		ctx,
		bson.M{"StockID": category.StockID, "Category": category.CategoryName},
		bson.M{"$set": bson.M{"Category": nil}},
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to update related products")
	}

	deleteRes, err := categoriesCol.DeleteOne(ctx, bson.M{"CategoryID": categoryUUID})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to delete category")
	}

	return c.JSON(fiber.Map{
		"updated_products": updateRes.ModifiedCount,
		"deleted_category": deleteRes.DeletedCount,
	})
}
