package repositories

import (
	"github.com/arvydasvapsva/product_catalog/app/models"
	elastic "gopkg.in/olivere/elastic.v5"
	"context"
	"reflect"
	"encoding/json"
	"github.com/revel/revel"
	"github.com/OwnLocal/goes"
	"os"
	"net/url"
)



type Tweet struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

func FindProducts() map[int] models.Product {



	h := os.Getenv("TEST_ELASTICSEARCH_HOST")
	if h == "" {
		h = "localhost"
	}

	p := os.Getenv("TEST_ELASTICSEARCH_PORT")
	if p == "" {
		p = "9200"
	}

	var conn = goes.NewClient(h, p)

	var query = map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"match_all": map[string]interface{}{},
				},
			},
		},
		"from":   0,
		"size":   100,
		"fields": []string{"onefield"},
		"filter": map[string]interface{}{
			"range": map[string]interface{}{
				"somefield": map[string]interface{}{
					"from":          "some date",
					"to":            "some date",
					"include_lower": false,
					"include_upper": false,
				},
			},
		},
	}

	extraArgs := make(url.Values, 1)

	searchResults, err := conn.Search(query, []string{"someindex"}, []string{""}, extraArgs)

	if err != nil {
		panic(err)
	}



	revel.INFO.Print(searchResults);


	ctx := context.Background()

	// Create a client
	//client, err := elastic.NewClient()
	client, err := elastic.NewClient(elastic.SetErrorLog(revel.INFO))//, elastic.SetBasicAuth("elastic", "changeme"), elastic.SetErrorLog(revel.INFO))
	if err != nil {
		// Handle error
		panic(err)
	}

	// Create an index
	_, err = client.CreateIndex("twitter").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	// Add a document to the index
	tweet := Tweet{User: "olivere", Message: "Take Five"}
	_, err = client.Index().
		Index("twitter").
		Type("tweet").
		Id("1").
		BodyJson(tweet).
		Refresh("true").
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	// Search with a term query
	termQuery := elastic.NewTermQuery("user", "olivere")
	searchResult, err := client.Search().
		Index("twitter").   // search in index "twitter"
		Query(termQuery).   // specify the query
		Sort("user", true). // sort by "user" field, ascending
		From(0).Size(10).   // take documents 0-9
		Pretty(true).       // pretty print request and response JSON
		Do(ctx)             // execute
	if err != nil {
		// Handle error
		panic(err)
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	revel.INFO.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over iterating the hits, see below.
	var ttyp Tweet
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(Tweet); ok {
			revel.INFO.Printf("Tweet by %s: %s\n", t.User, t.Message)
		}
	}
	// TotalHits is another convenience function that works even when something goes wrong.
	revel.INFO.Printf("Found a total of %d tweets\n", searchResult.TotalHits())

	// Here's how you iterate through results with full control over each step.
	if searchResult.Hits.TotalHits > 0 {
		revel.INFO.Printf("Found a total of %d tweets\n", searchResult.Hits.TotalHits)

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t Tweet
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}

			// Work with tweet
			revel.INFO.Printf("Tweet by %s: %s\n", t.User, t.Message)
		}
	} else {
		// No hits
		revel.INFO.Print("Found no tweets\n")
	}

	// Delete the index again
	_, err = client.DeleteIndex("twitter").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}








	return map[int] models.Product{
		1: models.NewProduct(1, "Wakeboard Shane", "The professional model by Shane Bonifay", 389.00),
		2: models.NewProduct(2, "Wakeboard GROOVE", "A stylish wakeboard with a fabtastic performance", 329.00),
		3: models.NewProduct(3, "Wakeboard S4", "The professional model by Phillip Soven", 389.00),
	}
}

func FindProduct(ProductId int) models.Product {
	var all = FindProducts()
	return all[ProductId];
}