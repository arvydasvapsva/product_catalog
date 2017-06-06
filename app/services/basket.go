package basket

import (
	"gopkg.in/resty.v0"
	"github.com/arvydasvapsva/product_catalog/app/models"
	"github.com/revel/revel"
	"github.com/arvydasvapsva/product_catalog/app/repositories"
	"errors"
	"fmt"
)

type remoteBasketItem struct {
	ProductId	int
	BasketId	string
	Amount		float32
}

func AddProduct(basketId string, Product models.Product, amount float32) (ProductId int, err error) {
	resp, _ := resty.R().
		SetBody(remoteBasketItem{Product.ProductId, basketId, amount}).
		Post("http://localhost:8888/add-to-basket")

	revel.INFO.Printf("Response from the backend:\n", resp.RawResponse)

	if resp.Status() != "200 OK" {
		return 0, errors.New(fmt.Sprintf("Product \"%s\" was not added to the basket", Product.Name))
	} else {
		repositories.StoreBasketItem(basketId, Product, amount)
	}

	return Product.ProductId, nil
}

func UpdateProduct(basketId string, updateProducts map[int]float32) (message string, err error) {

	var remoteBasketItems = map[int]remoteBasketItem{}

	for ProductId, amount := range updateProducts {
		remoteBasketItems[ProductId] = remoteBasketItem{ProductId, basketId, amount}
	}

	resp, _ := resty.R().
		SetBody(remoteBasketItems).
		Post("http://localhost:8888/update-basket")

	revel.INFO.Printf("Response from the backend:\n", resp.RawResponse)

	if resp.Status() != "200 OK" {
		return "", errors.New(fmt.Sprintf("Product \"%s\" was not updated.", "aaa"))
	} else {
		repositories.UpdateBasketItem(basketId, updateProducts)
	}

	return "Basket was successfully updated.", nil
}
