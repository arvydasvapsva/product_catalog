package basket

import (
	"gopkg.in/resty.v0"
	"errors"
	"github.com/arvydasvapsva/product_catalog/app/services"
	"github.com/arvydasvapsva/product_catalog/app/models"
)

type Remote struct {
	ParentBasketStorage services.BasketStorageInterface
}

func (basket *Remote) GetParentStorage() services.BasketStorageInterface  {
	return basket.ParentBasketStorage
}

func (basket *Remote) AddProducts(basketItems map[int]services.BasketItem) (int, error)  {
	parent := basket.GetParentStorage()
	if parent != nil {
		_, err := parent.AddProducts(basketItems)
		if err != nil {
			return 0, err
		}
	}

	for _, basketItem := range basketItems {
		resp, _ := resty.R().
			SetBody(basketItem).
			Post("http://localhost:8888/add-to-basket")

		if resp.Status() != "200 OK" {
			return 0, errors.New("Product was not added to the basket")
		}
	}

	return 1, nil
}

func (basket *Remote) UpdateProducts(basketItems map[int]services.BasketItem) (int, error) {
	parent := basket.GetParentStorage()
	if parent != nil {
		_, err := parent.UpdateProducts(basketItems)
		if err != nil {
			return 0, err
		}
	}

	resp, _ := resty.R().
		SetBody(basketItems).
		Post("http://localhost:8888/update-basket")

	if resp.Status() != "200 OK" {
		return 0, errors.New("Basket was not updated.")
	}

	return 1, nil
}

func (basket *Remote) FindBasketItems(basketId string) map[int] models.BasketItem  {
	return nil
}