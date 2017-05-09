package models

import "strconv"

type Product struct {
	ProductId		int
	Name, Description	string
	Price			float32
	DetailsUrl, ImageUrl, DetailsImageUrl	string
	BuyUrl string
}

func NewProduct(ProductId int, Name string, Description string, Price float32) Product {
	var DetailsUrl = "/details/" +  strconv.Itoa(ProductId)
	var BuyUrl = "/buy/" +  strconv.Itoa(ProductId)
	return Product{ProductId,Name,Description,Price,DetailsUrl,"", "", BuyUrl}
}

func (product Product) SetImageUrl(ImageUrl string) Product {
	product.ImageUrl = ImageUrl
	return product
}

func (product Product) SetDetailsImageUrl(DetailsImageUrl string) Product {
	product.DetailsImageUrl = DetailsImageUrl
	return product
}