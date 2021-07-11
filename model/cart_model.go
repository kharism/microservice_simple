package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Status int

const (
	WaitingForCheckedOut Status = iota
	WaitingForPayment
	Processing
	Delivery
	Done
)

func (d Status) String() string {
	return [...]string{"Waiting For Checked Out", "Waiting for Payment", "Processing", "Delivery", "Done"}[d]
}

type ItemAmmount struct {
	Item
	Ammount int64 `bson:"Amount" json:"Amount"`
}
type Cart struct {
	ID              string        `bson:"_id" json:"ID"`
	Items           []ItemAmmount `bson:"Items" json:"Items"`
	CheckedOut      bool          `bson:"CheckedOut" json:"CheckedOut"`
	UserId          string        `bson:"UserId" json:"UserId"`
	Total           float64       `bson:"Total" json:"Total"`
	DelivaryAddress string        `bson:"DeliveryAddress" json:"DeliveryAddress"`

	Status          Status     `bson:"Status" json:"Status"`
	CreatedBy       string     `bson:"CreatedBy" json:"CreatedBy"`
	CreatedDate     *time.Time `bson:"CreatedDate" json:"CreatedDate"`
	UpdatedBy       string     `bson:"UpdatedBy" json:"UpdatedBy"`
	LastUpdatedDate *time.Time `bson:"LastUpdatedDate" json:"LastUpdatedDate"`
}

func (m *Cart) TableName() string {
	return "Cart"
}
func (m *Cart) SetID(values []interface{}) {
	id := values[0]
	if v, ok := id.(string); ok {
		m.ID = v
	} else if v, ok := id.(primitive.ObjectID); ok {
		m.ID = v.Hex()
	}
}

// GetID get model id
func (m *Cart) GetID() ([]string, []interface{}) {
	return []string{"ID"}, []interface{}{m.ID}
}
func (m *Cart) GetTotal() float64 {
	sum := float64(0.0)
	for _, val := range m.Items {
		sum += val.Price * float64(val.Ammount)
	}
	return sum
}

func (m *Cart) RemoveItem(itemId string) {
	newItem := []ItemAmmount{}
	for _, v := range m.Items {
		if v.ID != itemId {
			newItem = append(newItem, v)
		}
	}
	m.Items = newItem
}
