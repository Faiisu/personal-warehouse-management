package handlers

import (
	"context"
	"strings"
	"time"

	"my-backend/internal/db"
	"my-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type createEventRequest struct {
	EventOwner string `json:"EventOwner"`
	Title      string `json:"Title"`
	StartAt    string `json:"StartAt"` // RFC3339
	EndAt      string `json:"EndAt"`   // RFC3339
	Location   string `json:"Location"`
	Status     string `json:"Status"`
}

// ListEvents godoc
// @Summary      List events
// @Description  Returns all events.
// @Tags         events
// @Produce      json
// @Success      200  {array}   models.Events
// @Failure      500  {object}  map[string]string
// @Router       /api/events [get]
func ListEvents(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection, err := db.EventsCollection(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "database unavailable")
	}

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch events")
	}
	defer cursor.Close(ctx)

	var events []models.Events
	if err := cursor.All(ctx, &events); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to decode events")
	}

	return c.JSON(events)
}

// CreateEvent godoc
// @Summary      Create a new event
// @Description  Creates an event with owner, title, time range, and optional location/status.
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        payload  body      createEventRequest  true  "Event data"
// @Success      201  {object}  models.Events
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/events [post]
func CreateEvent(c *fiber.Ctx) error {
	var req createEventRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON payload")
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Location = strings.TrimSpace(req.Location)
	req.Status = strings.TrimSpace(req.Status)
	req.EventOwner = strings.TrimSpace(req.EventOwner)

	if req.EventOwner == "" || req.Title == "" || req.StartAt == "" || req.EndAt == "" {
		return fiber.NewError(fiber.StatusBadRequest, "EventOwner, Title, StartAt, and EndAt are required")
	}

	eventOwnerUUID, err := uuid.Parse(req.EventOwner)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "EventOwner must be a valid UUID")
	}

	startAt, err := time.Parse(time.RFC3339, req.StartAt)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "StartAt must be RFC3339 timestamp")
	}

	endAt, err := time.Parse(time.RFC3339, req.EndAt)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "EndAt must be RFC3339 timestamp")
	}

	if endAt.Before(startAt) {
		return fiber.NewError(fiber.StatusBadRequest, "EndAt must be after StartAt")
	}

	status := req.Status
	if status == "" {
		status = "OPEN"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	collection, err := db.EventsCollection(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "database unavailable")
	}

	event := models.Events{
		EventID:    uuid.New(),
		EventOwner: eventOwnerUUID,
		Title:      req.Title,
		StartAt:    startAt,
		EndAt:      endAt,
		Location:   req.Location,
		Status:     status,
		CreatedAt:  time.Now().UTC(),
	}

	if _, err := collection.InsertOne(ctx, event); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fiber.NewError(fiber.StatusConflict, "event already exists")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create event")
	}

	return c.Status(fiber.StatusCreated).JSON(event)
}
