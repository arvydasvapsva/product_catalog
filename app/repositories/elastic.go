package repositories

import (
	"github.com/arvydasvapsva/product_catalog/app/models"
	elastic "gopkg.in/olivere/elastic.v5"
	"context"
	"encoding/json"
	"github.com/revel/revel"
	"strconv"
)

type Search struct {
	Key string ""
	From	int
	Size	int
}

func getClient(traceLog bool) *elastic.Client  {
	// Create a client

	if traceLog {
		client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"), elastic.SetBasicAuth("elastic", "changeme"), elastic.SetErrorLog(revel.INFO), elastic.SetSniff(false), elastic.SetTraceLog(revel.INFO), elastic.SetGzip(false))
		if err != nil {
			// Handle error
			panic(err)
		}

		return client
	} else {
		client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"), elastic.SetBasicAuth("elastic", "changeme"), elastic.SetErrorLog(revel.INFO), elastic.SetSniff(false), elastic.SetGzip(false))
		if err != nil {
			// Handle error
			panic(err)
		}
		return client
	}
}

func buildQuery(search Search) elastic.Query {
	if len(search.Key) == 0 {
		return elastic.NewMatchAllQuery()
	} else {
		return elastic.NewQueryStringQuery(search.Key)
	}
}

func FindProducts(search Search) map[int] models.Product {

	ctx := context.Background()
	var client = getClient(false)

	// Search with a term query
	searchResult, err := client.Search().
		Index("catalog").
		Type("product").
		Query(buildQuery(search)).   // specify the query
		//Sort("user", true).
		//From(search.From).
		//Size(search.Size).
		Pretty(true).       // pretty print request and response JSON
		Do(ctx)             // execute
	if err != nil {
		// Handle error
		panic(err)
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	revel.INFO.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	result := map[int] models.Product{}

	// Here's how you iterate through results with full control over each step.
	if searchResult.Hits.TotalHits > 0 {
		revel.INFO.Printf("Found a total of %d products\n", searchResult.Hits.TotalHits)

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			var product models.Product
			err := json.Unmarshal(*hit.Source, &product)
			if err != nil {
				// Deserialization failed
				revel.INFO.Print("Deserialization failed\n")
			}

			revel.INFO.Printf("Image Url: %s\n", product.ImageUrl)

			result[product.ProductId] = models.NewProduct(product.ProductId, product.Name, product.Description, product.Price).SetImageUrl(product.ImageUrl)
		}
	} else {
		// No hits
		revel.INFO.Print("Found no products\n")
	}

	return result
}

func FindProduct(ProductId int) models.Product {
	var all = FindProducts(Search{})
	return all[ProductId]
}

func FindBasketItems(basketId string) map[int] models.BasketItem  {

	ctx := context.Background()
	var client = getClient(false)

	// Search with a term query
	searchResult, err := client.Search().
		Index("basket").
		Type("basketItem").
		Pretty(true).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	//revel.INFO.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	result := map[int] models.BasketItem{}

	// Here's how you iterate through results with full control over each step.
	if searchResult.Hits.TotalHits > 0 {
		revel.INFO.Printf("Found a total of %d basket Items\n", searchResult.Hits.TotalHits)

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			var basketItem models.BasketItem
			err := json.Unmarshal(*hit.Source, &basketItem)
			if err != nil {
				// Deserialization failed
				revel.INFO.Print("Deserialization failed\n")
			}

			result[basketItem.ProductId] = basketItem
		}
	} else {
		// No hits
		revel.INFO.Print("No basket items found\n")
	}


	return result
}

func StoreBasketItem(basketId string, product models.Product, amount float32) models.BasketItem  {

	ctx := context.Background()
	var client = getClient(true)

	// Add a document to the index
	basketItem := models.BasketItem{product.ProductId, product.Name, product.Price, amount}

	client.Update().
		Index("basket").
		Id(basketId + strconv.Itoa(product.ProductId)).
		Type("basketItem").
		Upsert(map[string]interface{}{"ProductId": product.ProductId, "Name": product.Name, "Price": product.Price, "Amount": amount}).
		Script(elastic.NewScript("ctx._source.Amount += params.Amount").Lang("painless").Param("Amount", amount)).
		Refresh("true").
		Do(ctx)

	revel.INFO.Printf("Basket item \"%s\" was added", basketItem.Name)

	return basketItem
}

func UpdateBasketItem(basketId string, updateProducts map[int]float32) {

	ctx := context.Background()
	var client = getClient(true)

	bulkRequest := client.Bulk()
	for ProductId, amount := range updateProducts {
		if amount == 0 {
			var request = elastic.NewBulkDeleteRequest().
				Index("basket").
				Type("basketItem").
				Id(basketId + strconv.Itoa(ProductId))
			bulkRequest = bulkRequest.Add(request)
		} else {
			var request = elastic.NewBulkUpdateRequest().
				Index("basket").
				Type("basketItem").
				Id(basketId + strconv.Itoa(ProductId)).
				Doc(map[string]interface{}{"Amount": amount})
			bulkRequest = bulkRequest.Add(request)
		}
	}

	bulkRequest.Refresh("true").Do(ctx)
}
