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

	// nicolo: 
	// usually the idiomatic way to fail in go is to return an error
	// moreover, at that level, i like to distinguish between system errors and business errors
	// so what I usually do, I return a tuple
	// func AddProduct(p Product) (sys, business error)
	// but feel free to return just one if you don't like ti
	return resp.Status() == "200 OK"
}
