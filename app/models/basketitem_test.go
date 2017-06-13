package models_test

import (
	"testing"
	"github.com/arvydasvapsva/product_catalog/app/models"
	"reflect"
)

func TestBasketItem(t *testing.T)  {
	basketItem := new(models.BasketItem)
	if reflect.TypeOf(basketItem).String() != "*models.BasketItem" {
		t.Errorf("Unexpected BasketItem type %s", reflect.TypeOf(basketItem).String())
	}
}
