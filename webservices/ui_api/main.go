package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/kharism/microservice_simple/controller"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	debugging bool
	token     *jwtauth.JWTAuth
)
var (
	itemAPI controller.IItemRestAPI
)

func init() {
	token = jwtauth.New("HS256", []byte("somethingSecret"), nil)

	viper.SetConfigName("api")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath("../../config/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	debugging = viper.GetBool(`debug`)
	mongourl := os.Getenv("MONGO_URI")
	if mongourl != "" {
		viper.Set("uri", mongourl)
	}

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   debugging,
		FullTimestamp: true,
	})
	itemAPI = controller.NewItem(token)
}

func main() {
	r := chi.NewRouter()

	logger := logrus.New()
	logger.SetFormatter(&log.TextFormatter{
		ForceColors:   debugging,
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
		r.Mount("/item", itemAPI.Register())
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
	err := http.ListenAndServe(serverAddress, r)
	if err != nil {
		log.Fatal("unable to start web server", err.Error())
	}
}
