package repository

import (
	"context"
	"errors"
	"time"

	db "github.com/kharism/microservice_simple/connection"
	"github.com/kharism/microservice_simple/model"
	"github.com/kharism/microservice_simple/util"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IItem interface {
	Count(filters interface{}) (int64, error)
	HideById(id string) error
	FindAll(filters interface{}, skip, take int64) ([]model.Item, error)
	FindByID(id string) (model.Item, error)
	Save(data model.Item) (model.Item, error)
}

type itemRepo struct {
	client func() (*mongo.Client, error)
}

func NewItem() IItem {
	return &itemRepo{
		client: db.NewClient,
	}
}
func (r *itemRepo) Count(filters interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	item := model.Item{}
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
	query["Visible"] = true

	count, err := db.Collection(item.TableName()).CountDocuments(ctx, query)
	if err != nil {
		return -1, err
	}
	return count, nil
}
func (r *itemRepo) HideById(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	item := model.Item{}
	mongoCli, err := r.client()
	defer cancel()
	if err != nil {
		return err
	}
	defer mongoCli.Disconnect(ctx)
	db := mongoCli.Database(viper.GetString("db"))
	query := bson.M{"_id": id}

	_, err = db.Collection(item.TableName()).UpdateOne(ctx, query,
		bson.D{{
			"$set", bson.D{{"Visible", false}},
		}},
	)
	if err != nil {
		return err
	}
	return nil
}
func (r *itemRepo) Save(data model.Item) (model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	mongoCli, err := r.client()
	defer cancel()
	if err != nil {
		return model.Item{}, err
	}
	defer mongoCli.Disconnect(ctx)
	db := mongoCli.Database(viper.GetString("db"))
	if data.ID == "" {
		data.ID = util.RandString(23)
		insertOneRes, err := db.Collection(data.TableName()).InsertOne(ctx, data)
		if err != nil {
			return model.Item{}, err
		}
		data.SetID([]interface{}{insertOneRes.InsertedID})
	} else {
		_, err = db.Collection(data.TableName()).ReplaceOne(ctx, bson.M{"_id": data.ID}, data)
		if err != nil {
			return model.Item{}, err
		}
	}

	return data, nil
}
func (r *itemRepo) FindByID(id string) (model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	item := model.Item{}
	mongoCli, err := r.client()
	defer cancel()
	if err != nil {
		return model.Item{}, err
	}
	defer mongoCli.Disconnect(ctx)
	db := mongoCli.Database(viper.GetString("db"))
	query := bson.M{"_id": id}

	err = db.Collection(item.TableName()).FindOne(ctx, query).Decode(&item)
	if err != nil {
		return model.Item{}, err
	}
	if !item.Visible {
		return model.Item{}, errors.New("Not found")
	}
	return item, nil
}

// fetch all record match in the filters. The filters for mongo is map[string]interface{}
func (r *itemRepo) FindAll(filters interface{}, skip, take int64) ([]model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	item := model.Item{}
	result := []model.Item{}
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
	query["Visible"] = true
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
