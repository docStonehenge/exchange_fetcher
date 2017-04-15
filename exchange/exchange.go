package exchange

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type ExchangeResult struct {
	rawResult, PercentChange string
	ChangeInPoints           float64
}

func Fetch(url string) (*ExchangeResult, error) {
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	return &ExchangeResult{rawResult: string(bytes.TrimSpace(body))}, nil
}
