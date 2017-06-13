package catalog

import (
	"github.com/arvydasvapsva/product_catalog/app/models"
	elastic "gopkg.in/olivere/elastic.v5"
	"context"
	"encoding/json"
	"github.com/revel/revel"
	"github.com/arvydasvapsva/product_catalog/app/repositories/client"
	"github.com/arvydasvapsva/product_catalog/app/services"
)

func buildQuery(search services.CatalogSearch) elastic.Query {
	if len(search.Key) == 0 {
		return elastic.NewMatchAllQuery()
	} else {
		return elastic.NewQueryStringQuery(search.Key)
	}
}

type Elastic struct {
}

func getClient(logging bool) *elastic.Client {
	return client.ElasticClient{}.GetClient(logging)
}

func (*Elastic) FindProducts(search services.CatalogSearch) map[int] models.Product {
	ctx := context.Background()
	elasticClient := getClient(false)

	// Search with a term query
	searchResult, err := elasticClient.Search().
		Index("catalog").
		Type("product").
		Query(buildQuery(search)).   // specify the query
		Pretty(true).       // pretty print request and response JSON
		Do(ctx)             // execute
	if err != nil {
		// Handle error
		panic(err)
	}

	result := map[int] models.Product{}

	// Here's how you iterate through results with full control over each step.
	if searchResult.Hits.TotalHits > 0 {
		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			var product models.Product
			err := json.Unmarshal(*hit.Source, &product)
			if err != nil {
				// Deserialization failed
			}

			newProduct := models.NewProduct(product.ProductId, product.Name, product.Description, product.Price)
			newProduct.SetImageUrl(product.ImageUrl)

			result[product.ProductId] = newProduct
		}
	} else {
		// No hits
		revel.INFO.Print("Found no products\n")
	}

	return result
}

func (elastic *Elastic) FindProduct(ProductId int) models.Product {
	all := elastic.FindProducts(services.CatalogSearch{})
	return all[ProductId]
}