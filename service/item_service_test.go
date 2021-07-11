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

func TestCreateItem(t *testing.T) {
	itemService := NewItem()
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
			newItem, err := itemService.Save(itemModel)
			So(err, ShouldBeNil)
			So(newItem.ID, ShouldNotBeEmpty)
			item1Id := newItem.ID
			//mass insert
			for i := 0; i < 10; i++ {
				itemModel.ID = ""
				itemModel.Price += 0.9
				newItem, err = itemService.Save(itemModel)
				So(err, ShouldBeNil)
				So(newItem.ID, ShouldNotBeEmpty)
			}
			items, count, err := itemService.Find(map[string]interface{}{}, 0, 5)
			So(err, ShouldBeNil)
			So(count, ShouldEqual, 11)
			So(len(items), ShouldEqual, 5)
			item1, err := itemService.FindById(item1Id)
			So(err, ShouldBeNil)
			So(item1.ID, ShouldNotBeEmpty)
		})
	})
}
