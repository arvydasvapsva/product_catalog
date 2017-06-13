package basket

import (
	"github.com/arvydasvapsva/product_catalog/app/services"
	"github.com/arvydasvapsva/product_catalog/app/models"
	"github.com/arvydasvapsva/product_catalog/app/repositories/client"
)

type RedisBasket struct {
	ParentBasketStorage services.BasketStorageInterface
}

func (basket *RedisBasket) GetParentStorage() services.BasketStorageInterface  {
	return basket.ParentBasketStorage
}

func (basket *RedisBasket) AddProducts(basketItems map[int]services.BasketItem) (int, error)  {
	parent := basket.GetParentStorage()
	if parent != nil {
		_, err := parent.AddProducts(basketItems)
		if err != nil {
			return 0, err
		}
	}


	client := client.RedisClient{}
	client.GetClient()


	return 1, nil
}

func (basket *RedisBasket) UpdateProducts(basketItems map[int]services.BasketItem) (int, error) {
	parent := basket.GetParentStorage()
	if parent != nil {
		_, err := parent.UpdateProducts(basketItems)
		if err != nil {
			return 0, err
		}
	}

	return 1, nil
}

func (basket *RedisBasket) FindBasketItems(basketId string) map[int] models.BasketItem  {
	return nil
}