package main

import (
	//"kano/simwas/pkg/middleware"

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
	authAPI controller.IAuthRestAPI
)

func init() {
	token = jwtauth.New("HS256", []byte("somethingSecret"), nil)

	viper.SetConfigName("api")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath("../../config/")

	debugging = viper.GetBool(`debug`)

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   debugging,
		FullTimestamp: true,
	})
	authAPI = controller.NewAuth(token)
}

func main() {
	r := chi.NewRouter()

	logger := logrus.New()
	logger.SetFormatter(&log.TextFormatter{
		ForceColors:   debugging,
		FullTimestamp: true,
	})

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
		r.Mount("/auth", authAPI.Register())
	})
}
