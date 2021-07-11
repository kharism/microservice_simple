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

func TestCreateCart(t *testing.T) {
	itemService := NewItem()
	itemModel := model.Item{}
	cartService := NewCart()
	cartModel := model.Cart{}

	Convey("Clean up", t, func() {
		cli1, err := db.NewClient()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err = cli1.Database(viper.GetString("db")).Collection(itemModel.TableName()).DeleteMany(ctx, bson.M{})
		So(err, ShouldBeNil)
		_, err = cli1.Database(viper.GetString("db")).Collection(cartModel.TableName()).DeleteMany(ctx, bson.M{})
		So(err, ShouldBeNil)
		itemModel.ID = ""
		itemModel.ProductName = "lajsdlajskzcxc"
		itemModel.Price = 4.0
		itemModel.Visible = true
		newItem, err := itemService.Save(itemModel)
		//item1Id := newItem.ID
		itemModel.ID = ""
		itemModel.ProductName = "fasdas"
		itemModel.Price = 10.0
		itemModel.Visible = true
		newItem2, err := itemService.Save(itemModel)
		Convey("Try Save", func() {
			cartModel.ID = ""
			itemAmmount := model.ItemAmmount{Item: newItem, Ammount: 4}
			itemAmmount2 := model.ItemAmmount{Item: newItem2, Ammount: 10}
			cartModel.Items = []model.ItemAmmount{itemAmmount, itemAmmount2}
			cartModel.DelivaryAddress = "aaopqweqnlqenl"
			cartModel2, err := cartService.Save(cartModel)
			So(err, ShouldBeNil)
			So(cartModel2.ID, ShouldNotBeEmpty)

			cartModel3, err := cartService.FindById(cartModel2.ID)
			So(err, ShouldBeNil)
			So(cartModel3.ID, ShouldNotBeEmpty)

			carts, size, err := cartService.Find(map[string]interface{}{}, 0, 10)
			So(err, ShouldBeNil)
			So(size, ShouldEqual, 1)
			So(len(carts), ShouldEqual, 1)

		})
	})
}
