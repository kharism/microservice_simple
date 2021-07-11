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

func TestCreateItem(t *testing.T) {
	itemRepo := NewItem()
	itemModel := model.Item{}

	Convey("Clean up", t, func() {
		cli1, err := db.NewClient()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err = cli1.Database(viper.GetString("db")).Collection(itemModel.TableName()).DeleteMany(ctx, bson.M{})
		So(err, ShouldBeNil)
		Convey("Try Save", func() {
			itemModel.ID = ""
			itemModel.ProductName = "lajsdlajskzcxc"
			itemModel.Price = 29.80
			itemModel.Visible = true
			newItem, err := itemRepo.Save(itemModel)
			So(err, ShouldBeNil)
			So(newItem.ID, ShouldNotBeEmpty)
			item1Id := newItem.ID
			//mass insert
			for i := 0; i < 10; i++ {
				itemModel.ID = ""
				itemModel.Price += 0.9
				newItem, err = itemRepo.Save(itemModel)
				So(err, ShouldBeNil)
				So(newItem.ID, ShouldNotBeEmpty)
			}
			size, err := itemRepo.Count(map[string]interface{}{})
			So(err, ShouldBeNil)
			So(size, ShouldEqual, 11)
			err = itemRepo.HideById(item1Id)
			So(err, ShouldBeNil)
			respList, err := itemRepo.FindAll(map[string]interface{}{}, 0, 4)
			So(err, ShouldBeNil)
			So(len(respList), ShouldEqual, 4)

			newItem, err = itemRepo.FindByID(item1Id)
			So(err, ShouldNotBeNil)
			size, err = itemRepo.Count(map[string]interface{}{})
			So(err, ShouldBeNil)
			So(size, ShouldEqual, 10)

			//test update
			itemModel.ID = item1Id
			itemModel.ProductName = "lajsdlajskzcxc"
			itemModel.Price = 29.80
			itemModel.Visible = true
			_, err = itemRepo.Save(itemModel)
			So(err, ShouldBeNil)

		})
	})
}
