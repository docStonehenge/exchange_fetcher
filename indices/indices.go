package indices

import "encoding/json"

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
