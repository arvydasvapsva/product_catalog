package client

import (
	elastic "gopkg.in/olivere/elastic.v5"
	"github.com/revel/revel"
)

type ElasticClient struct {
}

func (*ElasticClient) GetClient(traceLog bool) *elastic.Client  {
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