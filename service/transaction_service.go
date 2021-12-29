package service

import (
	"log"
	"os"
	"sync"

	"github.com/kharism/microservice_simple/model"
	"github.com/kharism/microservice_simple/repository"
)

func HandleTransaction(transChan <-chan model.Transaction) {
	itemService := NewItem()
	logger := log.New(os.Stdout, "[TRN ROUTINE]", log.Ldate|log.Ltime|log.Lshortfile)
	transService := NewTransaction()
	logger.Println("Start Logging")
TRANS_LOOP:
	for trans := range transChan {
		logger.Println("Receiving Transaction", trans.CartId)
		backupItem := []model.Item{} //for rollback purpose
		itemIds := []string{}
		detailMap := map[string]model.TransactionDetail{}
		for _, item := range trans.Details {
			itemIds = append(itemIds, item.ItemId)
			detailMap[item.ItemId] = item
		}
		query := map[string]interface{}{}
		query["_id"] = map[string]interface{}{"$in": itemIds}
		logger.Println("Start Fetching item")
		backupItem, _, err := itemService.Find(query, 0, int64(len(itemIds)))

		if err != nil {
			logger.Println("Error Fetching item")
			trans.Status = model.STATUS_FAILED
			transService.Save(trans)
			continue
		}
		logger.Println("Done Getting Item", len(backupItem))
		if len(backupItem) != len(trans.Details) {
			logger.Println("Some item is not found")
			trans.Status = model.STATUS_FAILED
			transService.Save(trans)
			continue
		}
		committedItem := []model.Item{}
		logger.Println("Checking stock")
		for _, item := range backupItem {
			if item.Stock < detailMap[item.ID].Amount {
				logger.Println("Not enough item in stock", item.ProductName, "in", trans.CartId)
				logger.Println("Notify user")
				trans.Status = model.STATUS_FAILED
				transService.Save(trans)
				continue TRANS_LOOP
			}
			item.Stock -= detailMap[item.ID].Amount
			committedItem = append(committedItem, item)
		}
		//committing to db
		for _, i := range committedItem {
			itemService.Save(i)
		}
		trans.Status = model.STATUS_DONE
		transService.Save(trans)
		logger.Println("Done Transaction", trans.CartId)
	}
}

// IAuth auth service interface
type ITransaction interface {
	FindById(id string) (model.Transaction, error)
	//HideById(id string) error
	Save(item model.Transaction) (model.Transaction, error)
	Find(filter interface{}, skip, take int64) ([]model.Transaction, int64, error)
}
type transaction struct {
	transaction func() repository.ITransaction
}

// NewAuth create new service instance
func NewTransaction() ITransaction {
	return transaction{
		transaction: repository.NewTransaction,
	}
}
func (r transaction) FindById(id string) (model.Transaction, error) {
	return r.transaction().FindByID(id)
}

func (r transaction) Save(item model.Transaction) (model.Transaction, error) {
	//now := time.Now()
	//item.LastUpdatedDate = &now
	return r.transaction().Save(item)
}
func (r transaction) Find(filter interface{}, skip, take int64) ([]model.Transaction, int64, error) {
	var count int64
	if take == 0 {
		take = 10
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		var err error
		count, err = r.transaction().Count(filter)
		if err != nil {
			count = -1
		}
	}(&wg)
	result, err := r.transaction().FindAll(filter, skip, take)
	if err != nil {
		return result, count, err
	}
	wg.Wait()
	return result, count, nil
}
