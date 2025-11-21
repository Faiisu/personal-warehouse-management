package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a registered user stored in MongoDB.
type User struct {
	ID           uuid.UUID `bson:"_id" json:"id"`
	Email        string    `bson:"email" json:"email"`
	DisplayName  string    `bson:"display_name,omitempty" json:"display_name,omitempty"`
	PasswordHash string    `bson:"password_hash,omitempty" json:"-"`
	AvatarURL    string    `bson:"avatar_url,omitempty" json:"avatar_url,omitempty"`
	Status       string    `bson:"status" json:"status"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updated_at"`
}
