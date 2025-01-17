package serviceproduct

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"

	"ijahinventory/data/database"
	"ijahinventory/data/model/log"
	"ijahinventory/data/model/product"
	"ijahinventory/http/response"
)

var (
	db  *gorm.DB
	err error
)

func openDb() {
	db, err = database.OpenDb()
}

func closeDb() {
	db.Close()
}

//if sku not found, store new product
//else update existing price and adding the stock count and log it
//if sku not found and other product's fields are empty, return error
func StoreProduct(item product.Product,
	countOrder int,
	receiptNumber string,
	note string) (*product.Product, error) {
	openDb()
	defer closeDb()

	stockCount := item.StockCount
	newPrice := item.Price
	checkRow := db.First(&item, product.Product{Sku: item.Sku})
	if !checkRow.RecordNotFound() {
		//add the stock count and update the price
		item.StockCount += stockCount
		item.Price = newPrice
	} else {
		if item.Name == "" || item.Price == 0 || item.StockCount == 0 {
			return nil,
				errors.New("name, price, or stock_count cannot be empty")
		}
		if err := db.Create(&item).Error; err != nil {
			return nil, err
		}
	}

	logIngoing := log.LogIngoing{
		Timestamp:     time.Now(),
		CountOrder:    countOrder,
		CountReceived: item.StockCount,
		Product:       item,
		BuyPrice:      item.Price,
		Note:          note,
		ReceiptNumber: receiptNumber,
		TotalPrice:    item.Price * item.StockCount,
	}
	if err := db.Create(&logIngoing).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func LogOutgoing(
	sku string,
	countOutgoing int,
	note string,
) (*response.ResponseLogOutgoing, error) {
	openDb()
	defer closeDb()

	var item product.Product
	if db.First(
		&item,
		product.Product{Sku: sku},
	).RecordNotFound() {
		return nil, errors.New("Item with sku " + sku + " not found")
	}
	if item.StockCount < countOutgoing {
		return nil, errors.New("Item out of stock")
	}

	item.StockCount -= countOutgoing
	logOutgoing := log.LogOutgoing{
		Timestamp:     time.Now(),
		Product:       item,
		SalePrice:     item.Price,
		CountOutgoing: countOutgoing,
		TotalPrice:    item.Price * countOutgoing,
		Note:          note,
	}
	//update the stock count and log it
	if err := db.Create(&logOutgoing).Error; err != nil {
		return nil, err
	}

	//custom response body to remove unnecessary fields shown
	response := response.ResponseLogOutgoing{
		Timestamp:     logOutgoing.Timestamp,
		Sku:           logOutgoing.Product.Sku,
		ProductName:   logOutgoing.Product.Name,
		CountOutgoing: logOutgoing.CountOutgoing,
		Price:         logOutgoing.Product.Price,
		TotalPrice:    logOutgoing.TotalPrice,
		Note:          logOutgoing.Note,
	}

	return &response, nil
}
