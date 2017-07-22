package indices

import (
	"encoding/json"
	"github.com/docStonehenge/exchange_fetcher/exchange"
)

func Split(body []byte) (indices []string) {
	var idxJSON map[string]interface{}

	json.Unmarshal(body, &idxJSON)

	if indicesNode, ok := idxJSON["indices"].([]interface{}); ok {
		for _, idx := range indicesNode {
			if idxString, ok := idx.(string); ok {
				indices = append(indices, idxString)
			}
		}
	}

	return
}

func Join(exchanges map[string]exchange.Exchange) ([]byte, error) {
	body, err := json.Marshal(exchanges)

	if err == nil {
		return body, nil
	}

	return nil, err
}
