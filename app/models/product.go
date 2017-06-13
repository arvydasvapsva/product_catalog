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
	DetailsUrl := "/details/" +  strconv.Itoa(ProductId)
	BuyUrl := "/buy/" +  strconv.Itoa(ProductId)
	return Product{ProductId,Name,Description,Price,DetailsUrl,"", "", BuyUrl}
}

func (product *Product) SetImageUrl(ImageUrl string) {
	product.ImageUrl = ImageUrl
}

func (product *Product) SetDetailsImageUrl(DetailsImageUrl string)  {
	product.DetailsImageUrl = DetailsImageUrl
}