package models

import (
	"time"

	"github.com/google/uuid"
)

// Events represents an event stored in MongoDB.
type Events struct {
	EventID    uuid.UUID `bson:"EventID" json:"EventID"`
	EventOwner uuid.UUID `bson:"EventOwner" json:"EventOwner"`
	Title      string    `bson:"Title" json:"Title"`
	StartAt    time.Time `bson:"StartAt" json:"StartAt"`
	EndAt      time.Time `bson:"EndAt" json:"EndAt"`
	Location   string    `bson:"Location,omitempty" json:"Location,omitempty"`
	Status     string    `bson:"Status" json:"Status"`
	CreatedAt  time.Time `bson:"CreatedAt" json:"CreatedAt"`
}
