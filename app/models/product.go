package models

import "strconv"

type Product struct {
	ProductId		int
	Name, Description	string
	Price			float32
	DetailsUrl, ImageUrl	string
}

func NewProduct(ProductId int, Name string, Description string, Price float32) Product {
	var DetailsUrl = "details/" +  strconv.Itoa(ProductId);
	var ImageUrl = "http://demoshop.oxid-esales.com/professional-edition/out/pictures/generated/product/1/390_245_75/lf_shane_1.jpg"
	return Product{ProductId,Name,Description,Price,DetailsUrl,ImageUrl}
}
