package controller

import (
	"context"
	"encoding/json"
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
	//start webservice
	authAPI := NewAuth(token)
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
		r.Mount("/auth", authAPI.Register())
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
