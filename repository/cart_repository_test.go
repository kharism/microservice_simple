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

func TestCreateCart(t *testing.T) {
	itemRepo := NewItem()
	cartRepo := NewCart()
	itemModel := model.Item{}
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
		newItem, err := itemRepo.Save(itemModel)
		//item1Id := newItem.ID
		itemModel.ID = ""
		itemModel.ProductName = "fasdas"
		itemModel.Price = 10.0
		itemModel.Visible = true
		newItem2, err := itemRepo.Save(itemModel)
		//item2Id := newItem2.ID
		Convey("Try Save", func() {
			ammount := model.ItemAmmount{newItem, 4}
			//ammount.Ammount = 4
			cartModel.Items = []model.ItemAmmount{ammount}
			cartModel.DelivaryAddress = "asdasdadsd"
			cartModel.CheckedOut = false
			So(cartModel.GetTotal(), ShouldEqual, 16.0)
			cart2, err := cartRepo.Save(cartModel)
			So(err, ShouldBeNil)
			cartId := cart2.ID
			cartModel.ID = cartId
			ammount2 := model.ItemAmmount{newItem2, 10}
			//ammount2.Ammount = 10
			cartModel.Items = append(cartModel.Items, ammount2)
			So(cartModel.GetTotal(), ShouldEqual, 116.0)
			cart3, err := cartRepo.Save(cartModel)
			So(err, ShouldBeNil)
			So(cart3.ID, ShouldEqual, cart2.ID)

			aa, err := cartRepo.Count(map[string]interface{}{})
			So(err, ShouldBeNil)
			So(aa, ShouldEqual, 1)

			_, err = cartRepo.FindByID(cart3.ID)
			So(err, ShouldBeNil)

			carts4, err := cartRepo.FindAll(map[string]interface{}{}, 0, 10)
			So(err, ShouldBeNil)
			So(len(carts4), ShouldEqual, 1)
		})
	})
}
