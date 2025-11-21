package models

import "github.com/google/uuid"

// Products represents a product in stock.
type Products struct {
	ProductID   uuid.UUID `bson:"ProductID" json:"ProductID"`
	StockID     uuid.UUID `bson:"StockID" json:"StockID"`
	ProductName string    `bson:"ProductName" json:"ProductName"`
	Category    string    `bson:"Category,omitempty" json:"Category,omitempty"`
	Unit        string    `bson:"Unit,omitempty" json:"Unit,omitempty"`
	ProductQty  int       `bson:"ProductQty" json:"ProductQty"`
}
