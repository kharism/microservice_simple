package controller

import (
	"net/http"

	"github.com/kharism/microservice_simple/model"
	"github.com/kharism/microservice_simple/service"
	"github.com/kharism/microservice_simple/util"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

type ICartRestAPI interface {
	GetCart(w http.ResponseWriter, r *http.Request)
	GetCarts(w http.ResponseWriter, r *http.Request)
	SaveCart(w http.ResponseWriter, r *http.Request)
	AddItemToCart(w http.ResponseWriter, r *http.Request)
	RemoveItemFromCart(w http.ResponseWriter, r *http.Request)
	//HideItem(w http.ResponseWriter, r *http.Request)
	Register() http.Handler
}

type cartController struct {
	tokenAuth *jwtauth.JWTAuth
	cart      func() service.ICart
	item      func() service.IItem
	//rkas      func() service.IRKAS
}

func NewCart(token *jwtauth.JWTAuth) ICartRestAPI {
	return &cartController{
		tokenAuth: token,
		cart:      service.NewCart,
		item:      service.NewItem,
		//rkas:      service.NewRKAS,
	}
}

// get cart by id
// http method: GET
// router /{id}
// payload: None
func (c *cartController) GetCart(w http.ResponseWriter, r *http.Request) {
	id := util.URLParam(r, "id")
	item, err := c.cart().FindById(id)
	if err != nil {
		util.WriteJSONError(w, err)
	}
	util.WriteJSONData(w, item)
}

// get carts by id
// http method: POST
// router /list
// payload: {
//	Skip int64
//	Take int64
// }
func (c *cartController) GetCarts(w http.ResponseWriter, r *http.Request) {
	items := []model.Cart{}
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
	items, count, err := c.cart().Find(map[string]interface{}{}, param.Skip, param.Take)
	if err != nil {
		util.WriteJSONError(w, err)
		return
	}
	util.WriteJSONDataWithTotal(w, items, count)
}

// get carts by id
// http method: POST or PUT
// router / for POST /{id} for PUT
// payload: follow th cart model
func (c *cartController) SaveCart(w http.ResponseWriter, r *http.Request) {
	userID, err := util.GetClaimStringFromJWT(r, jwtKeyID)
	if err != nil {
		util.WriteJSONError(w, err)
		return
	}
	item := model.Cart{}
	err = util.ParsePayload(r, &item)
	if err != nil {
		util.WriteJSONError(w, err)
		return
	}
	item.CreatedBy = userID
	item, err = c.cart().Save(item)
	if err != nil {
		util.WriteJSONError(w, err)
		return
	}
	util.WriteJSONData(w, item)
}

type AddItemPayload struct {
	ItemId string
	Amount int64
}

// add item to cart
// http method: PUT
// router : /push/{id} for PUT
// payload: {
//    ItemId string
//    Amount int64
//}
func (c *cartController) AddItemToCart(w http.ResponseWriter, r *http.Request) {
	id := util.URLParam(r, "id")
	item, err := c.cart().FindById(id)
	if err != nil {
		util.WriteJSONError(w, err)
		return
	}
	//item.CreatedBy = userID
	userID, _ := util.GetClaimStringFromJWT(r, jwtKeyID)
	if userID != "" {
		item.CreatedBy = userID
		item.UserId = userID
	}
	itemPayload := AddItemPayload{}
	err = util.ParsePayload(r, &itemPayload)
	if err != nil {
		util.WriteJSONError(w, err)
		return
	}
	itemCart, err := c.item().FindById(itemPayload.ItemId)
	if err != nil {
		util.WriteJSONError(w, err)
		return
	}
	item.Items = append(item.Items, model.ItemAmmount{Item: itemCart, Ammount: itemPayload.Amount})
	item, err = c.cart().Save(item)
	if err != nil {
		util.WriteJSONError(w, err)
		return
	}
	util.WriteJSONData(w, item, "Success Add Item")
}

// add item to cart
// http method: PUT
// router : /pop/{id} for PUT
// payload: {
//    ItemId string
//    Amount int64
//}
func (c *cartController) RemoveItemFromCart(w http.ResponseWriter, r *http.Request) {
	id := util.URLParam(r, "id")
	item, err := c.cart().FindById(id)
	if err != nil {
		util.WriteJSONError(w, err)
		return
	}
	userID, _ := util.GetClaimStringFromJWT(r, jwtKeyID)
	if userID != "" {
		item.CreatedBy = userID
		item.UserId = userID
	}
	itemPayload := AddItemPayload{}
	err = util.ParsePayload(r, &itemPayload)
	if err != nil {
		util.WriteJSONError(w, err)
		return
	}
	item.RemoveItem(itemPayload.ItemId)
	item, err = c.cart().Save(item)
	if err != nil {
		util.WriteJSONError(w, err)
		return
	}
	util.WriteJSONData(w, item, "Success Remove Item")
}
func (c *cartController) Register() http.Handler {
	r := chi.NewRouter()

	r.Post("/list", c.GetCarts)
	//r.Post("/", c.GetItems)
	r.Get("/{id}", c.GetCart)
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(c.tokenAuth))
		r.Post("/", c.SaveCart)
		r.Put("/{id}", c.SaveCart)
		r.Put("/push/{id}", c.AddItemToCart)
		r.Put("/pop/{id}", c.RemoveItemFromCart)
		//r.Delete("/{id}", c.HideItem)
	})

	//r.Post("/registeruser", c.GetItems)

	return r
}
