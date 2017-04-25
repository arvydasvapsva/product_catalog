package controllers

import (
	"github.com/revel/revel"
	"github.com/arvydasvapsva/product_catalog/app/repositories"
	"strconv"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	var products = repositories.FindProducts();
	return c.Render(products)
}

func (c App) Details() revel.Result {
	var Id = c.Params.Get("id")
	var ProductId, _ = strconv.Atoi(Id)
	var product = repositories.FindProduct(ProductId);
	return c.Render(product)
}