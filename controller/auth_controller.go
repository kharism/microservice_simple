package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/eaciit/toolkit"
	"github.com/kharism/microservice_simple/util"

	"github.com/kharism/microservice_simple/service"

	"github.com/kharism/microservice_simple/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

const jwtKeyID = "_id"

// IAuthRestAPI user controller interface
// RegisterUser -> the register user request handler
// VerifyLogin  -> this is the one handling login request
type IAuthRestAPI interface {
	VerifyLogin(w http.ResponseWriter, r *http.Request)
	RegisterUser(w http.ResponseWriter, r *http.Request)
	Register() http.Handler
}

// authentication controller struct
// simple implementation with mongodb backend as prototyping backend
type authController struct {
	tokenAuth *jwtauth.JWTAuth
	auth      func() service.IAuth
	//rkas      func() service.IRKAS
}

// create new auth controller
func NewAuth(tokenAuth *jwtauth.JWTAuth) IAuthRestAPI {
	return &authController{
		auth: service.NewAuth,
		//rkas:      service.NewRKAS,
		tokenAuth: tokenAuth,
	}
}

// handler for register new user
func (c *authController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	data := model.User{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		util.WriteJSONError(w, err)
		return
	}
	err := c.auth().RegisterUser(data)
	if err != nil {
		util.WriteJSONError(w, err)
		return
	}
	response := toolkit.M{}
	response["IsError"] = false
	util.WriteJSONData(w, response)
}

// handler for login. On success login return the userdata with JWT Token.
// Probably need to remove unneccessary field due to security reason
func (c *authController) VerifyLogin(w http.ResponseWriter, r *http.Request) {
	var (
		user model.User
		err  error
	)

	// parse basic auth request
	data := model.User{}

	if err = json.NewDecoder(r.Body).Decode(&data); err != nil {
		util.WriteJSONError(w, err)
		return
	}

	user, err = c.auth().VerifyPassword(data.Username, data.Password)
	if err != nil {
		util.WriteJSONError(w, err)
		return
	}

	// use user id and username info on generating jwt token
	expirationInSecond := viper.GetInt64("jwt_expiration_duration")
	_, tokenString, _ := c.tokenAuth.Encode(jwt.MapClaims{
		jwtKeyID:   user.ID,
		"username": user.Username,
		"exp":      int64(time.Now().Add(time.Second * time.Duration(expirationInSecond)).Unix()),
	})

	user.Token = tokenString

	util.WriteJSONData(w, user)
}

func (c *authController) Register() http.Handler {
	r := chi.NewRouter()

	r.Post("/", c.VerifyLogin)
	r.Post("/registeruser", c.RegisterUser)

	return r
}
