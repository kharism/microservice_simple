package controller

import (
	"net/http"

	"github.com/kharism/microservice_simple/model"
	"github.com/kharism/microservice_simple/service"
	"github.com/kharism/microservice_simple/util"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

type IItemRestAPI interface {
	GetItem(w http.ResponseWriter, r *http.Request)
	GetItems(w http.ResponseWriter, r *http.Request)
	SaveItem(w http.ResponseWriter, r *http.Request)
	HideItem(w http.ResponseWriter, r *http.Request)
	Register() http.Handler
}

type itemController struct {
	tokenAuth *jwtauth.JWTAuth
	item      func() service.IItem
	//rkas      func() service.IRKAS
}

func NewItem(token *jwtauth.JWTAuth) IItemRestAPI {
	return &itemController{
		tokenAuth: token,
		item:      service.NewItem,
		//rkas:      service.NewRKAS,
	}
}

func (c *itemController) GetItem(w http.ResponseWriter, r *http.Request) {
	id := util.URLParam(r, "id")
	item, err := c.item().FindById(id)
	if err != nil {
		util.WriteJSONError(w, err)
	}
	util.WriteJSONData(w, item)
}
func (c *itemController) HideItem(w http.ResponseWriter, r *http.Request) {
	id := util.URLParam(r, "id")
	err := c.item().HideById(id)
	if err != nil {
		util.WriteJSONError(w, err)
	}
	util.WriteJSONData(w, model.Item{})
}

type requestPayload struct {
	Skip int64
	Take int64
}

func (c *itemController) GetItems(w http.ResponseWriter, r *http.Request) {
	items := []model.Item{}
	param := requestPayload{}
	err := util.ParsePayload(r, &param)
	if err != nil {
		if err.Error() != "EOF" {
			util.WriteJSONError(w, err)
			return
		}
		param.Skip = 0
		param.Take = 10
	}
	items, count, err := c.item().Find(map[string]interface{}{}, param.Skip, param.Take)
	if err != nil {
		util.WriteJSONError(w, err)
		return
	}
	util.WriteJSONDataWithTotal(w, items, count)
}
func (c *itemController) SaveItem(w http.ResponseWriter, r *http.Request) {
	userID, err := util.GetClaimStringFromJWT(r, jwtKeyID)
	if err != nil {
		util.WriteJSONError(w, err)
		return
	}
	item := model.Item{}
	err = util.ParsePayload(r, &item)
	if err != nil {
		util.WriteJSONError(w, err)
		return
	}
	item.CreatedBy = userID
	item, err = c.item().Save(item)
	if err != nil {
		util.WriteJSONError(w, err)
		return
	}
	util.WriteJSONData(w, item)
}
func (c *itemController) Register() http.Handler {
	r := chi.NewRouter()

	r.Post("/list", c.GetItems)
	//r.Post("/", c.GetItems)
	r.Get("/{id}", c.GetItem)
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(c.tokenAuth))
		r.Post("/", c.SaveItem)
		r.Put("/{id}", c.SaveItem)
		r.Delete("/{id}", c.HideItem)
	})

	//r.Post("/registeruser", c.GetItems)

	return r
}
