package controllers

import (
	"github.com/revel/revel"
	"strconv"
	"github.com/arvydasvapsva/product_catalog/app/services"
	"fmt"
	"github.com/arvydasvapsva/product_catalog/app/routes"
	"github.com/arvydasvapsva/product_catalog/app/repositories/catalog"
	repositoriesBasket "github.com/arvydasvapsva/product_catalog/app/repositories/basket"
)

const REQUEST_PRODUCT_ID = "id"
const REQUEST_PRODUCT_ACTION_REMOVE = "remove"
const REQUEST_PRODUCT_ACTION_SEARCH = "search"
const REQUEST_BASKET_FIELD_AMOUNT = "basket[amount]"

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	var products = getCatalog().FindProducts(services.CatalogSearch{})
	return c.Render(products)
}

func (c App) Details() revel.Result {
	var Id = c.Params.Get(REQUEST_PRODUCT_ID)
	var ProductId, _ = strconv.Atoi(Id)
	var product = getCatalog().FindProduct(ProductId)
	return c.Render(product)
}

func getBasketId(c App) string  {
	return c.Session.ID()
}

// basket service builder
func getBasketService() services.Basket {
	remoteStorage := &repositoriesBasket.Remote{}
	localStorage := &repositoriesBasket.ElasticBasket{remoteStorage}

	return services.NewBasketService(localStorage)
}

// catalog repository builder
func getCatalog() services.CatalogInterface {
	return &catalog.Elastic{}

    // you can pass a different data source!
	//return &catalog.Memory{}
}

func (c App) Buy() revel.Result  {
	var Id = c.Params.Get(REQUEST_PRODUCT_ID)
	var ProductId, _ = strconv.Atoi(Id)
	var product = getCatalog().FindProduct(ProductId)

	basketService := getBasketService()

	var _, err = basketService.AddProduct(getBasketId(c), product, 1)

	if err != nil {
		c.Flash.Error(err.Error())
	} else {
		c.Flash.Success(fmt.Sprintf("Product \"%s\" was added to the basket.", product.Name))
	}

	c.FlashParams()

	return c.Redirect(routes.App.Index())
}

func (c App) Basket() revel.Result {
	basketService := getBasketService()

	var basketItems = basketService.FindBasketItems(getBasketId(c))
	var basketItemsCount float32
	for _, v := range basketItems {
		basketItemsCount += v.Amount
	}
	return c.Render(basketItems, basketItemsCount)
}

func (c App) BasketUpdate() revel.Result {

	var updateProducts = map[int]float32{}
	var RemovableBasketItem = c.Params.Get(REQUEST_PRODUCT_ACTION_REMOVE)

	basketService := getBasketService()

	if RemovableBasketItem != "" {
		var ProductId, _ = strconv.Atoi(RemovableBasketItem)
		updateProducts[ProductId] = 0
	} else {
		c.Params.Bind(&updateProducts, REQUEST_BASKET_FIELD_AMOUNT)
	}

	var message, err = basketService.UpdateProduct(getBasketId(c), updateProducts)

	if err != nil {
		c.Flash.Error(err.Error())
	} else {
		c.Flash.Success(message)
	}

	c.FlashParams()

	return c.Redirect(routes.App.BasketDetails())
}

func (c App) BasketDetails() revel.Result {
	basketService := getBasketService()

	var basketItems = basketService.FindBasketItems(getBasketId(c))
	return c.Render(basketItems)
}

func (c App) Search() revel.Result {
	searchKey := c.Params.Get(REQUEST_PRODUCT_ACTION_SEARCH)
	var products = getCatalog().FindProducts(services.CatalogSearch{"*" + searchKey + "*", 0, 10})

	flash := make(map[string]string)
	flash["success"] = fmt.Sprintf("%d items found.", len(products))

	return c.Render(products, flash)
}