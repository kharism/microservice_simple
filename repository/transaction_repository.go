package repository

import (
	"context"
	"time"

	db "github.com/kharism/microservice_simple/connection"
	"github.com/kharism/microservice_simple/model"
	"github.com/kharism/microservice_simple/util"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ITransaction interface {
	Count(filters interface{}) (int64, error)
	FindAll(filters interface{}, skip, take int64) ([]model.Transaction, error)
	FindByID(id string) (model.Transaction, error)
	Save(data model.Transaction) (model.Transaction, error)
}

type transactionRepo struct {
	client func() (*mongo.Client, error)
}

func NewTransaction() ITransaction {
	return &transactionRepo{
		client: db.NewClient,
	}
}
func (r *transactionRepo) Count(filters interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	item := model.Transaction{}
	mongoCli, err := r.client()
	defer cancel()
	if err != nil {
		return -1, err
	}
	defer mongoCli.Disconnect(ctx)
	db := mongoCli.Database(viper.GetString("db"))
	query := bson.M{}
	query2 := filters.(map[string]interface{})
	for key, val := range query2 {
		query[key] = val
	}
	//query["Visible"] = true

	count, err := db.Collection(item.TableName()).CountDocuments(ctx, query)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (r *transactionRepo) Save(data model.Transaction) (model.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	mongoCli, err := r.client()
	defer cancel()
	if err != nil {
		return model.Transaction{}, err
	}
	defer mongoCli.Disconnect(ctx)
	db := mongoCli.Database(viper.GetString("db"))
	if data.ID == "" {
		data.ID = util.RandString(23)
		insertOneRes, err := db.Collection(data.TableName()).InsertOne(ctx, data)
		if err != nil {
			return model.Transaction{}, err
		}
		data.SetID([]interface{}{insertOneRes.InsertedID})
	} else {
		_, err = db.Collection(data.TableName()).ReplaceOne(ctx, bson.M{"_id": data.ID}, data)
		if err != nil {
			return model.Transaction{}, err
		}
	}

	return data, nil
}
func (r *transactionRepo) FindByID(id string) (model.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	item := model.Transaction{}
	mongoCli, err := r.client()
	defer cancel()
	if err != nil {
		return model.Transaction{}, err
	}
	defer mongoCli.Disconnect(ctx)
	db := mongoCli.Database(viper.GetString("db"))
	query := bson.M{}

	err = db.Collection(item.TableName()).FindOne(ctx, query).Decode(&item)
	if err != nil {
		return model.Transaction{}, err
	}
	return item, nil
}

// fetch all record match in the filters. The filters for mongo is map[string]interface{}
func (r *transactionRepo) FindAll(filters interface{}, skip, take int64) ([]model.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	item := model.Transaction{}
	result := []model.Transaction{}
	mongoCli, err := r.client()
	defer cancel()
	if err != nil {
		return result, err
	}
	defer mongoCli.Disconnect(ctx)
	db := mongoCli.Database(viper.GetString("db"))
	query := bson.M{}
	query2 := filters.(map[string]interface{})
	for key, val := range query2 {
		query[key] = val
	}
	//query["Visible"] = true
	options := options.FindOptions{}
	options.Skip = &skip
	options.Limit = &take
	cursor, err := db.Collection(item.TableName()).Find(ctx, query, &options)
	if err != nil {
		return result, err
	}
	err = cursor.All(ctx, &result)
	return result, err
}
