package connection

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/spf13/viper"
)

var (
	connectionString string
)

// NewClient return mongo client stuff. You need to close the connection yourself once done.
func NewClient() (*mongo.Client, error) {
	connectionString := viper.GetString("uri")

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client, nil
}
