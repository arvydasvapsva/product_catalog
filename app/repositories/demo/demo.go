package demo

import (
	"github.com/arvydasvapsva/product_catalog/app/models"
)

func FindProducts() map[int] models.Product {
	return map[int] models.Product{
		1: models.NewProduct(1, "Wakeboard Shane", "The professional model by Shane Bonifay", 389.00),
		2: models.NewProduct(2, "Wakeboard GROOVE", "A stylish wakeboard with a fabtastic performance", 329.00),
		3: models.NewProduct(3, "Wakeboard S4", "The professional model by Phillip Soven", 389.00),
	}
}

func FindProduct(ProductId int) models.Product {
	var all = FindProducts()
	return all[ProductId];
}