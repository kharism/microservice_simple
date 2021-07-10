package model

import (
	"time"

	"github.com/kharism/microservice_simple/util"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              string     `bson:"_id" json:"ID"`
	Name            string     `bson:"Name" json:"Name"`
	NIP             string     `bson:"NIP" json:"NIP"`
	Username        string     `bson:"Username" json:"Username"`
	Password        string     `bson:"-" json:"Password"`
	PasswordHash    string     `bson:"PasswordHash" json:"-"`
	Email           string     `bson:"Email" json:"Email"`
	Token           string     `bson:"-" json:"Token"`
	LastLogin       *time.Time `bson:"LastLogin" json:"LastLogin"`
	CreatedBy       string     `bson:"CreatedBy" json:"CreatedBy"`
	CreatedDate     *time.Time `bson:"CreatedDate" json:"CreatedDate"`
	UpdatedBy       string     `bson:"UpdatedBy" json:"UpdatedBy"`
	LastUpdatedDate *time.Time `bson:"LastUpdatedDate" json:"LastUpdatedDate"`
}

func (m *User) TableName() string {
	return "User"
}
func (m *User) SetID(values []interface{}) {
	id := values[0]
	if v, ok := id.(string); ok {
		m.ID = v
	} else if v, ok := id.(primitive.ObjectID); ok {
		m.ID = v.Hex()
	}
}

// GetID get model id
func (m *User) GetID() ([]string, []interface{}) {
	return []string{"ID"}, []interface{}{m.ID}
}

// VerifyPassword verify plain password input with hashed password
func (m *User) VerifyPassword(password string) bool {
	return util.HashCompare(password, m.PasswordHash)
}
