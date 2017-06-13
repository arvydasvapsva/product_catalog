package models_test

import (
	"testing"
	"github.com/arvydasvapsva/product_catalog/app/models"
)

func TestNewProduct(t *testing.T) {
	productId := 1
	name := "Name"
	description := "Description"
	price := float32(100.00)

	product := models.NewProduct(productId, name, description, price)

	if product.ProductId != productId || product.Name != name || product.Description != description || product.Price != price {
		t.Errorf("Somethign wrong with the product creation, field values do not match")
	}
}

func TestProduct_SetImageUrl(t *testing.T) {
	testUrl := "ImageUrl"

	product := models.Product{}
	product.SetImageUrl(testUrl)

	if testUrl != product.ImageUrl {
		t.Errorf("Image URL does not match %s != %s", testUrl, product.ImageUrl)
	}
}

func TestProduct_SetDetailsImageUrl(t *testing.T) {
	testDetailsUrl := "DetailsImageUrl"

	product := models.Product{}
	product.SetDetailsImageUrl(testDetailsUrl)

	if testDetailsUrl != product.DetailsImageUrl {
		t.Errorf("Details image URL does not match %s != %s", testDetailsUrl, product.DetailsImageUrl)
	}
}