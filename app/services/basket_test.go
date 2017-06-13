package services_test

import (
	"testing"
	"errors"
	"github.com/arvydasvapsva/product_catalog/app/models"
	"github.com/arvydasvapsva/product_catalog/app/services"
)

//
type mock_RemoteBasketStorageSuccess struct {
	basketStorage services.BasketStorageInterface
}

func (basket *mock_RemoteBasketStorageSuccess) AddProducts(basketItems map[int]services.BasketItem) (int, error)  {
	parent := basket.GetParentStorage()
	if parent != nil {
		_, err := parent.AddProducts(basketItems)
		if err != nil {
			return 0, err
		}
	}

	return 1, nil
}

func (basket *mock_RemoteBasketStorageSuccess) UpdateProducts(basketItems map[int]services.BasketItem) (int, error) {
	parent := basket.GetParentStorage()
	if parent != nil {
		_, err := parent.UpdateProducts(basketItems)
		if err != nil {
			return 0, err
		}
	}

	return 1, nil
}

func (basket *mock_RemoteBasketStorageSuccess) GetParentStorage() services.BasketStorageInterface  {
	return basket.basketStorage
}

//
type mock_RemoteBasketStorageError struct {
	basketStorage services.BasketStorageInterface
}

func (basket *mock_RemoteBasketStorageError) AddProducts(basketItems map[int]services.BasketItem) (int, error)  {
	parent := basket.GetParentStorage()
	if parent != nil {
		_, err := parent.AddProducts(basketItems)
		if err != nil {
			return 0, err
		}
	}

	return 0, errors.New("Product was not added to the basket")
}

func (basket *mock_RemoteBasketStorageError) UpdateProducts(basketItems map[int]services.BasketItem) (int, error) {
	parent := basket.GetParentStorage()
	if parent != nil {
		_, err := parent.UpdateProducts(basketItems)
		if err != nil {
			return 0, err
		}
	}

	return 0, errors.New("Basket was not updated.")
}

func (basket *mock_RemoteBasketStorageError) GetParentStorage() services.BasketStorageInterface  {
	return basket.basketStorage
}

//
type mock_LocalBasketStorage struct {
	basketStorage services.BasketStorageInterface
}

func (basket *mock_LocalBasketStorage) AddProducts(basketItems map[int]services.BasketItem) (int, error)  {
	parent := basket.GetParentStorage()
	if parent != nil {
		_, err := parent.AddProducts(basketItems)
		if err != nil {
			return 0, err
		}
	}

	return 1, nil
}

func (basket *mock_LocalBasketStorage) UpdateProducts(basketItems map[int]services.BasketItem) (int, error)  {
	parent := basket.GetParentStorage()
	if parent != nil {
		_, err := parent.UpdateProducts(basketItems)
		if err != nil {
			return 0, err
		}
	}

	return 1, nil
}

func (basket *mock_LocalBasketStorage) GetParentStorage() services.BasketStorageInterface  {
	return basket.basketStorage
}

func TestBasket_AddProductSuccess(t *testing.T) {
	basketService := services.NewBasketService(&mock_LocalBasketStorage{&mock_RemoteBasketStorageSuccess{}})
	res, err := basketService.AddProduct("1", models.NewProduct(1, "Name", "Description", float32(100)), float32(100))

	if res == 0 {
		t.Errorf("Error %s", err)
	}
}

func TestBasket_AddProductError(t *testing.T) {
	basketService := services.NewBasketService(&mock_LocalBasketStorage{&mock_RemoteBasketStorageError{}})
	_, err := basketService.AddProduct("1", models.NewProduct(1, "Name", "Description", float32(100)), float32(100))

	if err == nil {
		t.Errorf("Unexpected, that product was added to basket")
	}
}

func TestBasket_UpdateProductSuccess(t *testing.T) {
	basketItems := make(map[int]float32)
	basketItems[1] = float32(100)

	basketService := services.NewBasketService(&mock_LocalBasketStorage{&mock_RemoteBasketStorageSuccess{}})
	res, err := basketService.UpdateProduct("1", basketItems)

	if res == "" {
		t.Errorf("Error %s", err)
	}
}

func TestBasket_UpdateProductError(t *testing.T) {
	basketItems := make(map[int]float32)
	basketItems[1] = float32(100)

	basketService := services.NewBasketService(&mock_LocalBasketStorage{&mock_RemoteBasketStorageError{}})
	_, err := basketService.UpdateProduct("1", basketItems)

	if err == nil {
		t.Errorf("Unexpected, that product was updated in basket")
	}
}