package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/eaciit/toolkit"
	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	db "github.com/kharism/microservice_simple/connection"
	"github.com/kharism/microservice_simple/model"
	"github.com/kharism/microservice_simple/repository"
	"github.com/kharism/microservice_simple/service"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	token := jwtauth.New("HS256", []byte("secretTest"), nil)
	viper.SetConfigName("api_test")
	viper.SetConfigType("json")
	viper.AddConfigPath("../config/")
	viper.AddConfigPath("./config/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	transChan := make(chan model.Transaction, 100)
	go service.HandleTransaction(transChan)
	//start webservice
	authAPI := NewAuth(token)
	itemAPI := NewItem(token)
	cartAPI := NewCart(token, transChan)
	r := chi.NewRouter()

	logger := logrus.New()
	logger.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.Info("server api run on DEBUG mode")
	log.SetLevel(log.DebugLevel)

	r.Use(chiMiddleware.RequestID)

	r.Use(chiMiddleware.Recoverer)

	// disable cache control
	r.Use(chiMiddleware.NoCache)

	// apply gzip compression
	r.Use(chiMiddleware.Compress(5, "gzip"))

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)
	r.Group(func(r chi.Router) {
		r.Get("/ping", ping)
		//r.Use(jwtauth.Verifier(token))
		r.Mount("/auth", authAPI.Register())
		r.Mount("/item", itemAPI.Register())
		r.Mount("/cart", cartAPI.Register())
	})

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Debugf("[%s] %s", method, route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		log.Errorf("walk function error : %s\n", err.Error())
	}

	serverAddress := viper.GetString("address")
	log.Infof("server api run at %s", serverAddress)

	go func() {
		err = http.ListenAndServe(serverAddress, r)
		if err != nil {
			log.Fatal("unable to start web server", err.Error())
		}
	}()
}
func ToStringReader(payload toolkit.M) io.Reader {
	return strings.NewReader(toolkit.JsonString(payload))
}
func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("pong"))
}
func ProcessResponse(resp *http.Response) (toolkit.M, error) {
	responseJson := toolkit.M{}
	contentByte, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return responseJson, err
	}
	err = json.Unmarshal(contentByte, &responseJson)
	if err != nil {
		return responseJson, err
	}
	return responseJson, nil
}
func TestCreateUser(t *testing.T) {
	//userService := service.NewAuth()
	userModel := model.User{}

	Convey("Clean up", t, func() {
		cli1, err := db.NewClient()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err = cli1.Database(viper.GetString("db")).Collection(userModel.TableName()).DeleteMany(ctx, bson.M{})
		So(err, ShouldBeNil)
		payload := toolkit.M{}
		client := &http.Client{}
		Convey("Try Register", func() {
			payload["Username"] = "admin"
			payload["Password"] = "PasswordXX"
			payloadReader := strings.NewReader(toolkit.JsonString(payload))

			resp, err := client.Post("http://localhost:8098/auth/registeruser", "application/json", payloadReader)
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp.StatusCode, ShouldEqual, 200)
			user, err := repository.NewUser().FindByUsername("admin")
			So(err, ShouldBeNil)
			So(user.Username, ShouldEqual, "admin")
			So(user.PasswordHash, ShouldNotBeEmpty)
			//try to register again, should be error
			payloadReader = strings.NewReader(toolkit.JsonString(payload))

			resp, err = client.Post("http://localhost:8098/auth/registeruser", "application/json", payloadReader)
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp.StatusCode, ShouldEqual, 500)
			respJson, err := ProcessResponse(resp)
			So(err, ShouldBeNil)
			So(respJson["Message"], ShouldEqual, "User Sudah ada")
			Convey("Try Login", func() {
				//t.Log("Try Login Fail")
				payload["Username"] = "admin"
				payload["Password"] = "Password"
				payloadReader := strings.NewReader(toolkit.JsonString(payload))
				resp, err := client.Post("http://localhost:8098/auth", "application/json", payloadReader)
				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
				So(resp.StatusCode, ShouldEqual, 500)
				//respJson, err = ProcessResponse(resp)
				//So(err, ShouldBeNil)
				//t.Log(respJson)
				//So(respJson["Message"], ShouldEqual, "User Sudah ada")
				payload["Username"] = "admin"
				payload["Password"] = "PasswordXX"
				payloadReader = strings.NewReader(toolkit.JsonString(payload))
				resp, err = client.Post("http://localhost:8098/auth", "application/json", payloadReader)
				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
				So(resp.StatusCode, ShouldEqual, 200)
			})
		})

	})
}
func AddRequestBearer(req *http.Request, token string) {
	req.Header.Add("Authorization", "BEARER "+token)
}
func TestItems(t *testing.T) {
	userModel := model.User{}
	itemModel := model.Item{}
	Convey("Clean up", t, func() {
		cli1, err := db.NewClient()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err = cli1.Database(viper.GetString("db")).Collection(userModel.TableName()).DeleteMany(ctx, bson.M{})
		So(err, ShouldBeNil)
		_, err = cli1.Database(viper.GetString("db")).Collection(itemModel.TableName()).DeleteMany(ctx, bson.M{})
		So(err, ShouldBeNil)
		payload := toolkit.M{}
		client := &http.Client{}
		//create basic user
		payload["Username"] = "admin"
		payload["Password"] = "PasswordXX"
		payloadReader := strings.NewReader(toolkit.JsonString(payload))

		resp, err := client.Post("http://localhost:8098/auth/registeruser", "application/json", payloadReader)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		So(resp.StatusCode, ShouldEqual, 200)
		//login
		payload["Username"] = "admin"
		payload["Password"] = "PasswordXX"
		payloadReader = strings.NewReader(toolkit.JsonString(payload))
		resp, err = client.Post("http://localhost:8098/auth", "application/json", payloadReader)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		So(resp.StatusCode, ShouldEqual, 200)

		responseJson, err := ProcessResponse(resp)
		So(err, ShouldBeNil)
		So(responseJson, ShouldNotBeNil)
		//t.Log(responseJson)
		token := responseJson["Data"].(map[string]interface{})["Token"].(string)
		Convey("Test Item API", func() {
			//create item without login should fail
			itemModel.ID = ""
			itemModel.ProductName = "LLL1212"
			itemModel.Price = 10.0
			itemModel.Visible = true
			payloadReader = strings.NewReader(toolkit.JsonString(itemModel))
			resp, err = client.Post("http://localhost:8098/item", "application/json", payloadReader)
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp.StatusCode, ShouldNotEqual, 200)
			//responseJson, err := ProcessResponse(resp)
			So(err, ShouldBeNil)
			//t.Log(responseJson)

			//create item with login token
			payloadReader = strings.NewReader(toolkit.JsonString(itemModel))
			req, err := http.NewRequest("POST", "http://localhost:8098/item", payloadReader)
			So(err, ShouldBeNil)
			AddRequestBearer(req, token)
			resp, err = client.Do(req)
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			responseJson, err = ProcessResponse(resp)
			So(err, ShouldBeNil)
			//t.Log(responseJson)
			So(resp.StatusCode, ShouldEqual, 200)
			itemId := responseJson["Data"].(map[string]interface{})["ID"]
			So(itemId, ShouldNotBeEmpty)

			//get item without login
			resp, err = client.Get("http://localhost:8098/item/" + itemId.(string))
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp.StatusCode, ShouldEqual, 200)
			payload = map[string]interface{}{
				"Skip": 0,
				"Take": 10,
			}
			payloadReader = strings.NewReader(toolkit.JsonString(payload))
			resp, err = client.Post("http://localhost:8098/item/list", "application/json", payloadReader)
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			responseJson, err = ProcessResponse(resp)
			So(err, ShouldBeNil)
			//t.Log(responseJson)
			So(resp.StatusCode, ShouldEqual, 200)

			//update item
			t.Log("update Item")
			itemModel.ID = itemId.(string)
			itemModel.ProductName = "LLL1212"
			itemModel.Price = 13.0
			itemModel.Visible = true
			payloadReader = strings.NewReader(toolkit.JsonString(itemModel))
			request, err := http.NewRequest("PUT", "http://localhost:8098/item/"+itemId.(string), payloadReader)
			So(err, ShouldBeNil)
			resp, err = client.Do(request)
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			//resp, err = client.("http://localhost:8098/item", "application/json", payloadReader)

			//"delete" item without login
			request, err = http.NewRequest("DELETE", "http://localhost:8098/item/"+itemId.(string), nil)
			So(err, ShouldBeNil)
			resp, err = client.Do(request)
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			responseJson, err = ProcessResponse(resp)
			So(err, ShouldBeNil)
			//t.Log(responseJson)
			So(resp.StatusCode, ShouldEqual, 200)

		})
	})
}
func TestCarts(t *testing.T) {
	userModel := model.User{}
	itemModel := model.Item{}
	cartModel := model.Cart{}
	transModel := model.Transaction{}
	cartRepo := repository.NewCart()
	//transactionRepo := repository.NewCart()
	//itemRepo := repository.NewItem()
	Convey("Clean up", t, func() {
		//time.Sleep(10 * time.Second)
		cli1, err := db.NewClient()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err = cli1.Database(viper.GetString("db")).Collection(userModel.TableName()).DeleteMany(ctx, bson.M{})
		So(err, ShouldBeNil)
		_, err = cli1.Database(viper.GetString("db")).Collection(itemModel.TableName()).DeleteMany(ctx, bson.M{})
		So(err, ShouldBeNil)
		_, err = cli1.Database(viper.GetString("db")).Collection(cartModel.TableName()).DeleteMany(ctx, bson.M{})
		So(err, ShouldBeNil)
		_, err = cli1.Database(viper.GetString("db")).Collection(transModel.TableName()).DeleteMany(ctx, bson.M{})
		So(err, ShouldBeNil)
		payload := toolkit.M{}
		client := &http.Client{}
		//create basic user
		payload["Username"] = "admin"
		payload["Password"] = "PasswordXX"
		payloadReader := strings.NewReader(toolkit.JsonString(payload))

		resp, err := client.Post("http://localhost:8098/auth/registeruser", "application/json", payloadReader)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		So(resp.StatusCode, ShouldEqual, 200)
		payload["Username"] = "admin2"
		payloadReader = strings.NewReader(toolkit.JsonString(payload))
		resp, err = client.Post("http://localhost:8098/auth/registeruser", "application/json", payloadReader)
		payload["Username"] = "admin3"
		payloadReader = strings.NewReader(toolkit.JsonString(payload))
		resp, err = client.Post("http://localhost:8098/auth/registeruser", "application/json", payloadReader)
		//login
		payload["Username"] = "admin"
		payload["Password"] = "PasswordXX"
		payloadReader = strings.NewReader(toolkit.JsonString(payload))
		resp, err = client.Post("http://localhost:8098/auth", "application/json", payloadReader)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		So(resp.StatusCode, ShouldEqual, 200)

		responseJson, err := ProcessResponse(resp)
		So(err, ShouldBeNil)
		So(responseJson, ShouldNotBeNil)
		//t.Log(responseJson)
		token := responseJson["Data"].(map[string]interface{})["Token"].(string)

		payload["Username"] = "admin2"
		payloadReader = strings.NewReader(toolkit.JsonString(payload))
		resp, err = client.Post("http://localhost:8098/auth", "application/json", payloadReader)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		So(resp.StatusCode, ShouldEqual, 200)

		responseJson, err = ProcessResponse(resp)
		So(err, ShouldBeNil)
		So(responseJson, ShouldNotBeNil)
		//t.Log(responseJson)
		token2 := responseJson["Data"].(map[string]interface{})["Token"].(string)

		payload["Username"] = "admin3"
		payloadReader = strings.NewReader(toolkit.JsonString(payload))
		resp, err = client.Post("http://localhost:8098/auth", "application/json", payloadReader)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		So(resp.StatusCode, ShouldEqual, 200)

		responseJson, err = ProcessResponse(resp)
		So(err, ShouldBeNil)
		So(responseJson, ShouldNotBeNil)
		//t.Log(responseJson)
		token3 := responseJson["Data"].(map[string]interface{})["Token"].(string)

		itemModel.ID = ""
		itemModel.ProductName = "LLL1212"
		itemModel.Price = 10.0
		itemModel.Stock = 200
		itemModel.Visible = true
		payloadReader = strings.NewReader(toolkit.JsonString(itemModel))
		req, err := http.NewRequest("POST", "http://localhost:8098/item", payloadReader)
		So(err, ShouldBeNil)
		AddRequestBearer(req, token)
		resp, err = client.Do(req)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		responseJson, err = ProcessResponse(resp)
		So(err, ShouldBeNil)
		itemId := responseJson["Data"].(map[string]interface{})["ID"].(string)
		Convey("Test Cart", func() {
			cartModel = model.Cart{}
			cartModel.ID = ""
			cartModel.DelivaryAddress = "adadoiqueiqowe"
			payloadReader = strings.NewReader(toolkit.JsonString(cartModel))
			req, err := http.NewRequest("POST", "http://localhost:8098/cart", payloadReader)
			So(err, ShouldBeNil)
			AddRequestBearer(req, token)
			resp, err = client.Do(req)
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			responseJson, err = ProcessResponse(resp)
			So(err, ShouldBeNil)
			// t.Log(responseJson)
			cartId := responseJson["Data"].(map[string]interface{})["ID"].(string)

			resp, err = client.Get("http://localhost:8098/cart/" + cartId)
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp.StatusCode, ShouldEqual, 200)
			payload = map[string]interface{}{
				"Skip": 0,
				"Take": 10,
			}
			payloadReader = strings.NewReader(toolkit.JsonString(payload))
			resp, err = client.Post("http://localhost:8098/cart/list", "application/json", payloadReader)
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			responseJson, err = ProcessResponse(resp)
			So(err, ShouldBeNil)
			So(responseJson["Data"], ShouldNotBeNil)

			payload = map[string]interface{}{
				"ItemId": itemId,
				"Amount": 100,
			}
			payloadReader = strings.NewReader(toolkit.JsonString(payload))
			req, err = http.NewRequest("PUT", "http://localhost:8098/cart/push/"+cartId, payloadReader)
			So(err, ShouldBeNil)
			AddRequestBearer(req, token)
			resp, err = client.Do(req)
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp.StatusCode, ShouldEqual, 200)
			tempCart, err := cartRepo.FindByID(cartId)
			So(err, ShouldBeNil)
			So(len(tempCart.Items), ShouldEqual, 1)
			Convey("Test Pop", func() {
				payloadReader = strings.NewReader(toolkit.JsonString(payload))
				req, err = http.NewRequest("PUT", "http://localhost:8098/cart/pop/"+cartId, payloadReader)
				So(err, ShouldBeNil)
				AddRequestBearer(req, token)
				resp, err = client.Do(req)
				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
				tempCart, err = cartRepo.FindByID(cartId)
				So(err, ShouldBeNil)
				So(len(tempCart.Items), ShouldEqual, 0)
				// So(err, ShouldBeNil)
			})
			//
			// Convey("Test Checkout", func() {
			// 	// test checkout
			// 	req, err = http.NewRequest("POST", "http://localhost:8098/cart/checkout/"+cartId, nil)
			// 	So(err, ShouldBeNil)
			// 	AddRequestBearer(req, token)
			// 	resp, err = client.Do(req)
			// 	So(err, ShouldBeNil)
			// 	So(resp, ShouldNotBeNil)
			// 	//transactionRepo.Count()

			// 	responseJson, err = ProcessResponse(resp)
			// 	So(err, ShouldBeNil)
			// 	// t.Log(responseJson)
			// 	transId := responseJson["Data"].(map[string]interface{})["ID"].(string)
			// 	So(transId, ShouldNotBeEmpty)
			// 	cc, _ := transactionRepo.Count(map[string]interface{}{})
			// 	fmt.Println("TRANSID", transId, cc)
			// 	//time.Sleep(2 * time.Second)

			// 	trans, err := transactionRepo.FindByID(transId)
			// 	So(err, ShouldBeNil)
			// 	So(trans.Status, ShouldEqual, model.STATUS_DONE)
			// 	itemNew, err := itemRepo.FindByID(itemId)
			// 	So(err, ShouldBeNil)
			// 	So(itemNew.Stock, ShouldEqual, 100)
			// })
			Convey("Test Parallel Checkout", func() {
				//time.Sleep(10)
				cartModel := model.Cart{}
				cartModel.ID = ""
				cartModel.DelivaryAddress = "Address2"
				payloadReader = strings.NewReader(toolkit.JsonString(cartModel))
				req, err = http.NewRequest("POST", "http://localhost:8098/cart", payloadReader)
				So(err, ShouldBeNil)
				AddRequestBearer(req, token2)
				resp, err = client.Do(req)
				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
				responseJson, err = ProcessResponse(resp)
				So(err, ShouldBeNil)
				// t.Log(responseJson)
				cartId2 := responseJson["Data"].(map[string]interface{})["ID"].(string)

				cartModel = model.Cart{}
				cartModel.ID = ""
				cartModel.DelivaryAddress = "Address3"
				payloadReader = strings.NewReader(toolkit.JsonString(cartModel))
				req, err = http.NewRequest("POST", "http://localhost:8098/cart", payloadReader)
				So(err, ShouldBeNil)
				AddRequestBearer(req, token3)
				resp, err = client.Do(req)
				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
				responseJson, err = ProcessResponse(resp)
				So(err, ShouldBeNil)
				// t.Log(responseJson)
				cartId3 := responseJson["Data"].(map[string]interface{})["ID"].(string)

				t.Log("Cart IDs", cartId, cartId2, cartId3)

				payload = map[string]interface{}{
					"ItemId": itemId,
					"Amount": 100,
				}
				payloadReader = strings.NewReader(toolkit.JsonString(payload))
				req, err = http.NewRequest("PUT", "http://localhost:8098/cart/push/"+cartId2, payloadReader)
				So(err, ShouldBeNil)
				AddRequestBearer(req, token2)
				resp, err = client.Do(req)
				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
				So(resp.StatusCode, ShouldEqual, 200)

				tempCart, err := cartRepo.FindByID(cartId2)
				So(err, ShouldBeNil)
				So(len(tempCart.Items), ShouldEqual, 1)

				payloadReader = strings.NewReader(toolkit.JsonString(payload))
				req, err = http.NewRequest("PUT", "http://localhost:8098/cart/push/"+cartId3, payloadReader)
				So(err, ShouldBeNil)
				AddRequestBearer(req, token3)
				resp, err = client.Do(req)
				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
				So(resp.StatusCode, ShouldEqual, 200)

				checkoutCarts := []string{cartId, cartId2, cartId3}
				tokens := []string{token, token2, token3}
				// failedChan := make(chan string)
				//wg := sync.WaitGroup{}
				//wg.Add(1)
				for idx, i := range checkoutCarts {
					go func(cartid, token string) {
						req, _ := http.NewRequest("POST", "http://localhost:8098/cart/checkout/"+cartId, nil)
						AddRequestBearer(req, token)
						client.Do(req)

					}(i, tokens[idx])
				}
				time.Sleep(10 * time.Second)
				query := map[string]interface{}{"CartId": map[string]interface{}{"$in": checkoutCarts}}
				transactions, err := cartRepo.FindAll(query, 0, int64(len(checkoutCarts)))
				So(err, ShouldBeNil)
				for _, trans := range transactions {
					if trans.Status == model.Status(model.STATUS_FAILED) {
						fmt.Println("transaction id", trans.ID, "Failed")
					}
				}
			})
		})

	})
}
