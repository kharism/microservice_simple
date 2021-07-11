package service

import (
	"sync"
	"time"

	"github.com/kharism/microservice_simple/model"
	"github.com/kharism/microservice_simple/repository"
	//"kano/simwas/pkg/module/toolkit"
)

// IAuth auth service interface
type ICart interface {
	FindById(id string) (model.Cart, error)
	Save(item model.Cart) (model.Cart, error)
	Find(filter interface{}, skip, take int64) ([]model.Cart, int64, error)
}
type cart struct {
	cart func() repository.ICart
}

// NewAuth create new service instance
func NewCart() ICart {
	return cart{
		cart: repository.NewCart,
	}
}
func (r cart) FindById(id string) (model.Cart, error) {
	return r.cart().FindByID(id)
}
func (r cart) Save(item model.Cart) (model.Cart, error) {
	now := time.Now()
	item.LastUpdatedDate = &now
	return r.cart().Save(item)
}
func (r cart) Find(filter interface{}, skip, take int64) ([]model.Cart, int64, error) {
	var count int64
	if take == 0 {
		take = 10
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		var err error
		count, err = r.cart().Count(filter)
		if err != nil {
			count = -1
		}
	}(&wg)
	result, err := r.cart().FindAll(filter, skip, take)
	if err != nil {
		return result, count, err
	}
	wg.Wait()
	return result, count, nil
}
