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

func AddProduct(basketId string, Product models.Product) (ProductId int, err error) {
	resp, _ := resty.R().
		SetBody(remoteBasketItem{Product.ProductId, basketId, 100}).
		Post("http://sellercenter.local/test.php")

	revel.INFO.Printf("Response from the backend:\n", resp.RawResponse)
	repositories.StoreBasketItem(basketId, Product)

	if resp.Status() != "200 OK" {
		return 0, errors.New(fmt.Sprintf("Product \"%s\" was not added to the basket", Product.Name))
	}

	return Product.ProductId, nil
}
