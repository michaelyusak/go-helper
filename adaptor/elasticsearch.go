package adaptor

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v9"
)

func ConnectElasticSearch(config elasticsearch.Config) (*elasticsearch.Client, error) {
	es, err := elasticsearch.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("[adaptor][ConnectElasticSearch][Open] error: %w", err)
	}

	return es, nil
}
