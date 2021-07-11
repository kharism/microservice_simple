package service

import (
	"sync"
	"time"

	"github.com/kharism/microservice_simple/model"
	"github.com/kharism/microservice_simple/repository"
	//"kano/simwas/pkg/module/toolkit"
)

// IAuth auth service interface
type IItem interface {
	FindById(id string) (model.Item, error)
	HideById(id string) error
	Save(item model.Item) (model.Item, error)
	Find(filter interface{}, skip, take int64) ([]model.Item, int64, error)
}
type item struct {
	item func() repository.IItem
}

// NewAuth create new service instance
func NewItem() IItem {
	return item{
		item: repository.NewItem,
	}
}
func (r item) FindById(id string) (model.Item, error) {
	return r.item().FindByID(id)
}
func (r item) HideById(id string) error {
	return r.item().HideById(id)
}
func (r item) Save(item model.Item) (model.Item, error) {
	now := time.Now()
	item.LastUpdatedDate = &now
	return r.item().Save(item)
}
func (r item) Find(filter interface{}, skip, take int64) ([]model.Item, int64, error) {
	var count int64
	if take == 0 {
		take = 10
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		var err error
		count, err = r.item().Count(filter)
		if err != nil {
			count = -1
		}
	}(&wg)
	result, err := r.item().FindAll(filter, skip, take)
	if err != nil {
		return result, count, err
	}
	wg.Wait()
	return result, count, nil
}
