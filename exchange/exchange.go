package exchange

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const baseURL = "https://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20yahoo.finance.quotes%20where%20symbol%20in%20(%22INDEXES%22)&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys&callback="

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

func BuildURL(indexes []string) string {
	match := regexp.MustCompile("INDEXES")
	indexesList := strings.Join(indexes, ",")

	url := match.ReplaceAllString(baseURL, indexesList)

	return url
}
