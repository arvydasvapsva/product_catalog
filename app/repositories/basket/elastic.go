package basket

import (
	elastic "gopkg.in/olivere/elastic.v5"
	"context"
	"strconv"
	"github.com/arvydasvapsva/product_catalog/app/repositories/client"
	"github.com/arvydasvapsva/product_catalog/app/services"
	"github.com/arvydasvapsva/product_catalog/app/models"
	"encoding/json"
)

type ElasticBasket struct {
	ParentBasketStorage services.BasketStorageInterface
}

func (basket *ElasticBasket) GetParentStorage() services.BasketStorageInterface  {
	return basket.ParentBasketStorage
}

func getClient(logging bool) *elastic.Client {
	return client.ElasticClient{}.GetClient(logging)
}

func (basket *ElasticBasket) AddProducts(basketItems map[int]services.BasketItem) (int, error)  {
	parent := basket.GetParentStorage()
	if parent != nil {
		_, err := parent.AddProducts(basketItems)
		if err != nil {
			return 0, err
		}
	}

	ctx := context.Background()
	elasticClient := getClient(true)

	for _, basketItem := range basketItems {
		elasticClient.Update().
			Index("basket").
			Id(basketItem.BasketId + strconv.Itoa(basketItem.ProductId)).
			Type("basketItem").
			Upsert(map[string]interface{}{"ProductId": basketItem.ProductId, "Name": basketItem.Name, "Price": basketItem.Price, "Amount": basketItem.Amount}).
			Script(elastic.NewScript("ctx._source.Amount += params.Amount").Lang("painless").Param("Amount", basketItem.Amount)).
			Refresh("true").
			Do(ctx)
	}

	return 1, nil
}

func (basket *ElasticBasket) UpdateProducts(basketItems map[int]services.BasketItem) (int, error) {
	parent := basket.GetParentStorage()
	if parent != nil {
		_, err := parent.UpdateProducts(basketItems)
		if err != nil {
			return 0, err
		}
	}

	ctx := context.Background()
	elasticClient := getClient(true)
	bulkRequest := elasticClient.Bulk()

	for _, basketItem := range basketItems {
		if basketItem.Amount == 0 {
			request := elastic.NewBulkDeleteRequest().
				Index("basket").
				Type("basketItem").
				Id(basketItem.BasketId + strconv.Itoa(basketItem.ProductId))
			bulkRequest = bulkRequest.Add(request)
		} else {
			request := elastic.NewBulkUpdateRequest().
				Index("basket").
				Type("basketItem").
				Id(basketItem.BasketId + strconv.Itoa(basketItem.ProductId)).
				Doc(map[string]interface{}{"Amount": basketItem.Amount})
			bulkRequest = bulkRequest.Add(request)
		}
	}

	bulkRequest.Refresh("true").Do(ctx)

	return 1, nil
}

func (basket *ElasticBasket) FindBasketItems(basketId string) map[int] models.BasketItem  {
	ctx := context.Background()
	elasticClient := getClient(false)

	// Search with a term query
	searchResult, err := elasticClient.Search().
		Index("basket").
		Type("basketItem").
		Pretty(true).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	result := map[int] models.BasketItem{}

	// Here's how you iterate through results with full control over each step.
	if searchResult.Hits.TotalHits > 0 {

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			var basketItem models.BasketItem
			err := json.Unmarshal(*hit.Source, &basketItem)
			if err != nil {
				// Deserialization failed
			}

			result[basketItem.ProductId] = basketItem
		}
	} else {
		// No hits
	}

	return result
}