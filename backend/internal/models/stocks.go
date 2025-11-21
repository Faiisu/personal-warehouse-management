package models

import "github.com/google/uuid"

// Stocks represents a user's stock list.
type Stocks struct {
	StockID   uuid.UUID `bson:"StockID" json:"StockID"`
	UserID    uuid.UUID `bson:"UserID" json:"UserID"`
	StockName string    `bson:"StockName" json:"StockName"`
}
