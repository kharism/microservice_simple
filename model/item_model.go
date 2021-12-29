package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID              string     `bson:"_id" json:"ID"`
	ProductName     string     `bson:"Name" json:"Name"`
	Price           float64    `bson:"Price" json:"Price"`
	Visible         bool       `bson:"Visible" json:"Visible"`
	CreatedBy       string     `bson:"CreatedBy" json:"CreatedBy"`
	CreatedDate     *time.Time `bson:"CreatedDate" json:"CreatedDate"`
	UpdatedBy       string     `bson:"UpdatedBy" json:"UpdatedBy"`
	Stock           int        `bson:"Stock" json:"Stock"`
	LastUpdatedDate *time.Time `bson:"LastUpdatedDate" json:"LastUpdatedDate"`
}

func (m *Item) TableName() string {
	return "Item"
}
func (m *Item) SetID(values []interface{}) {
	id := values[0]
	if v, ok := id.(string); ok {
		m.ID = v
	} else if v, ok := id.(primitive.ObjectID); ok {
		m.ID = v.Hex()
	}
}

// GetID get model id
func (m *Item) GetID() ([]string, []interface{}) {
	return []string{"ID"}, []interface{}{m.ID}
}
