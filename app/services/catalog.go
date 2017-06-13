package services

import "github.com/arvydasvapsva/product_catalog/app/models"

type CatalogSearch struct {
	Key string ""
	From	int
	Size	int
}

type CatalogInterface interface {
	FindProducts(search CatalogSearch) (map[int] models.Product)
	FindProduct(ProductId int) (models.Product)
}