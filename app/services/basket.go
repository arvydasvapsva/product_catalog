package services

import (
	"github.com/arvydasvapsva/product_catalog/app/models"
	"errors"
	"fmt"
)

type BasketStorageInterface interface {
	AddProducts(map[int]BasketItem) (int, error)
	UpdateProducts(map[int]BasketItem) (int, error)
	GetParentStorage() BasketStorageInterface
	FindBasketItems(basketId string) map[int] models.BasketItem
}

type Basket struct {
	basketStorage BasketStorageInterface
}

func NewBasketService(basketStorage BasketStorageInterface) Basket {
	return Basket{basketStorage}
}

type BasketItem struct {
	ProductId	int
	Name string
	Price float32
	BasketId string
	Amount float32
}

func (basket *Basket) getBasketStorage() BasketStorageInterface {
	return basket.basketStorage
}

func (basket *Basket) AddProduct(basketId string, Product models.Product, amount float32) (int, error) {
	basketItems := map[int]BasketItem{}
	basketItems[Product.ProductId] = BasketItem{Product.ProductId, Product.Name, Product.Price, basketId, amount}

	_, err := basket.getBasketStorage().AddProducts(basketItems)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Product \"%s\" was not added to the basket", Product.Name))
	}

	return Product.ProductId, nil
}

func (basket *Basket) UpdateProduct(basketId string, updateProducts map[int]float32) (string, error) {
	basketItems := map[int]BasketItem{}
	for ProductId, amount := range updateProducts {
		basketItems[ProductId] = BasketItem{ProductId, "",0, basketId, amount}
	}

	_, err := basket.getBasketStorage().UpdateProducts(basketItems)
	if err != nil {
		return "", errors.New("Basker was not updated.")
	}

	return "Basket was successfully updated.", nil
}

func (basket *Basket) FindBasketItems(basketId string) map[int] models.BasketItem  {
	return basket.getBasketStorage().FindBasketItems(basketId)
}