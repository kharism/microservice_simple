package repository

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
func TestCreateUser(t *testing.T) {
	userRepo := NewUser()
	userModel := model.User{}

	Convey("Clean up", t, func() {
		cli1, err := db.NewClient()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err = cli1.Database(viper.GetString("db")).Collection(userModel.TableName()).DeleteMany(ctx, bson.M{})
		So(err, ShouldBeNil)
		Convey("Try Save", func() {
			userModel.Name = "adasd"
			userModel.Username = "KKK"
			userModel.Password = "adasdasd"
			userModel.PasswordHash = "asdadsas"
			_, err = userRepo.Save(userModel)
			So(err, ShouldBeNil)
			Convey("Get By Username", func() {
				user2, err := userRepo.FindByUsername(userModel.Username)
				So(err, ShouldBeNil)
				So(user2.Username, ShouldEqual, userModel.Username)
				So(user2.Password, ShouldBeEmpty)
			})
		})

	})
}
