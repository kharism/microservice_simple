package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionStatus int

const (
	STATUS_DONE     TransactionStatus = iota // use this if transaction is done without any error
	STATUS_CHECKING                          // transaction is being verified
	STATUS_FAILED                            // transaction is failed to be fulfilled
)

type Transaction struct {
	ID               string              `bson:"_id" json:"ID"`
	UserId           string              `bson:"UserId" json:"UserId"`
	CartId           string              `bson:"CartId" json:"CartId"`
	TreansactionDate time.Time           `bson:"TreansactionDate" json:"TreansactionDate"`
	Details          []TransactionDetail `bson:"Details" json:"Details"`
	Status           TransactionStatus   `bson:"Status" json:"Status"`
	Total            int                 `bson:"Total" json:"Total"`
}

func (t *Transaction) TableName() string {
	return "Transaction"
}

func (m *Transaction) SetID(values []interface{}) {
	id := values[0]
	if v, ok := id.(string); ok {
		m.ID = v
	} else if v, ok := id.(primitive.ObjectID); ok {
		m.ID = v.Hex()
	}
}

type TransactionDetail struct {
	ItemId   string `bson:"ItemId" json:"ItemId"`
	Subtotal int    `bson:"Subtotal" json:"Subtotal"`
	Amount   int    `bson:"Amount" json:"Amount"`
}
