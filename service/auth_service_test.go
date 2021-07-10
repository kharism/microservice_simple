package service

import (
	"context"
	"testing"
	"time"

	db "github.com/kharism/microservice_simple/connection"
	"github.com/kharism/microservice_simple/model"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	viper.SetConfigName("api_test")
	viper.SetConfigType("json")
	viper.AddConfigPath("../config/")
	viper.AddConfigPath("./config/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func TestRegisterAndLogin(t *testing.T) {
	auth := NewAuth()
	userModel := model.User{}

	Convey("Clean up", t, func() {
		cli1, err := db.NewClient()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err = cli1.Database(viper.GetString("db")).Collection(userModel.TableName()).DeleteMany(ctx, bson.M{})
		So(err, ShouldBeNil)
		Convey("Try Register", func() {
			userModel.Username = "KK"
			userModel.Password = "SDSS"
			err := auth.RegisterUser(userModel)
			So(err, ShouldBeNil)
			user2, err := auth.VerifyPassword("KK", "SDSS")
			So(err, ShouldBeNil)
			So(user2.Password, ShouldBeEmpty)
		})
	})
}
