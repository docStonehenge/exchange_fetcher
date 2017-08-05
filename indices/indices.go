package indices

import (
	"encoding/json"
	"github.com/docStonehenge/exchange_fetcher/exchange"
	"strings"
	"unicode"
)

func SplitJSONBody(body []byte) (indices []string) {
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

func SplitListBody(body string) (indices []string) {
	removeSpacesAndCommas := func(character rune) bool {
		return unicode.IsSpace(character) ||
			character == ',' || character == ';'
	}

	indices = strings.FieldsFunc(body, removeSpacesAndCommas)
	return
}

func Join(exchanges map[string]exchange.Exchange) ([]byte, error) {
	body, err := json.Marshal(exchanges)

	if err == nil {
		return body, nil
	}

	return nil, err
}
