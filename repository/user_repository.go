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
)

type IUser interface {
	//Count(filters ...*dbflex.Filter) (int64, error)
	//Delete(filters ...*dbflex.Filter) error
	//DeleteByID(id string) error
	//FindAll(param *BaseParam) ([]model.User, error)
	//FindAtasan(id string) ([]model.User, error)
	//FindByID(id string) (model.User, error)
	FindByUsername(username string) (model.User, error)
	Save(data model.User) (model.User, error)
}

type userRepo struct {
	client func() (*mongo.Client, error)
}

func NewUser() IUser {
	return &userRepo{
		client: db.NewClient,
	}
}
func (r *userRepo) Save(data model.User) (model.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	mongoCli, err := r.client()
	defer cancel()
	if err != nil {
		return model.User{}, err
	}
	defer mongoCli.Disconnect(ctx)
	db := mongoCli.Database(viper.GetString("db"))
	if data.ID == "" {
		data.ID = util.RandString(23)
		insertOneRes, err := db.Collection(data.TableName()).InsertOne(ctx, data)
		if err != nil {
			return model.User{}, err
		}
		data.SetID([]interface{}{insertOneRes.InsertedID})
	} else {
		_, err = db.Collection(data.TableName()).ReplaceOne(ctx, bson.M{"_id": data.ID}, data)
		if err != nil {
			return model.User{}, err
		}
	}

	return data, nil
}
func (r *userRepo) FindByUsername(username string) (model.User, error) {
	data := model.User{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoCli, err := r.client()
	if err != nil {
		return model.User{}, err
	}

	defer mongoCli.Disconnect(ctx)
	db := mongoCli.Database(viper.GetString("db"))
	query := bson.M{"Username": username}

	err = db.Collection(data.TableName()).FindOne(ctx, query).Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}
