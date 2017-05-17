package controllers

import (
	"github.com/revel/revel"
	"github.com/arvydasvapsva/product_catalog/app/repositories"
	"strconv"
	"github.com/arvydasvapsva/product_catalog/app/services"
	"fmt"
	"github.com/arvydasvapsva/product_catalog/app/routes"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	var products = repositories.FindProducts()
	return c.Render(products)
}

func (c App) Details() revel.Result {
	var Id = c.Params.Get("id")
	var ProductId, _ = strconv.Atoi(Id)
	var product = repositories.FindProduct(ProductId)
	return c.Render(product)
}

func getBasketId(c App) string  {
	var sessionId = c.Session.ID()

	revel.INFO.Printf("SessionId %s", sessionId)

	return sessionId
}

func (c App) Buy() revel.Result  {
	var Id = c.Params.Get("id")
	var ProductId, _ = strconv.Atoi(Id)
	var product = repositories.FindProduct(ProductId)
	var _, err = basket.AddProduct(getBasketId(c), product)

	if err != nil {
		c.Flash.Error(err.Error())
	} else {
		c.Flash.Success(fmt.Sprintf("Product \"%s\" was added to the basket", product.Name))
	}

	c.FlashParams()

	return c.Redirect(routes.App.Index())
}

func (c App) Basket() revel.Result {
	var basketItems = repositories.FindBasketItems(getBasketId(c))
	var basketItemsCount = len(basketItems)
	return c.Render(basketItems, basketItemsCount)
}

func (c App) BasketDetails() revel.Result {
	var basketItems = repositories.FindBasketItems(getBasketId(c))

	return c.Render(basketItems)
}
