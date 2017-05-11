package basket

import (
	"gopkg.in/resty.v0"
	"strconv"
	"github.com/arvydasvapsva/product_catalog/app/models"
	"github.com/revel/revel"
	"github.com/arvydasvapsva/product_catalog/app/repositories"
)

func AddProduct(basketId string, Product models.Product) bool {
	resp, _ := resty.R().Post("http://demoshop.oxid-esales.com/professional-edition/" + strconv.Itoa(Product.ProductId))

	revel.INFO.Printf("Response from the backend:\n", resp.RawResponse)
	repositories.StoreBasketItem(basketId, Product)

	// nicolo: usually the idiomatic way to fail in go is to return an error
	return resp.Status() == "200 OK"
}
