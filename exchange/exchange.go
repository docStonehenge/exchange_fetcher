package exchange

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const baseURL = "https://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20yahoo.finance.quotes%20where%20symbol%20in%20(%22INDEXES%22)&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys"

type ExchangesResult struct {
	rawResult string
	Exchanges map[string]Exchange
}

type Exchange struct {
	Name, Symbol, PercentChange, ChangeInPoints, LastTradeDate, LastTradeTime string
}

func Fetch(url string) (*ExchangesResult, error) {
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	return &ExchangesResult{rawResult: string(bytes.TrimSpace(body))}, nil
}

func BuildURL(indexes []string) string {
	match := regexp.MustCompile("INDEXES")
	indexesList := strings.Join(indexes, ",")

	url := match.ReplaceAllString(baseURL, indexesList)

	return url
}

func (ex *ExchangesResult) Parse() {
	ex.Exchanges = make(map[string]Exchange)
	exchanges := ex.parseRawResultToJSON()

	if exchange, ok := exchanges["quote"].(map[string]interface{}); ok {
		ex.setExchangeByNameKey(exchange)
	}

	if exchanges, ok := exchanges["quote"].([]interface{}); ok {
		for _, exchange := range exchanges {
			exchange := exchange.(map[string]interface{})
			ex.setExchangeByNameKey(exchange)
		}
	}
}

func (ex *ExchangesResult) parseRawResultToJSON() map[string]interface{} {
	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(ex.rawResult), &jsonMap)

	exchanges := jsonMap["query"].(map[string]interface{})
	return exchanges["results"].(map[string]interface{})
}

func (ex *ExchangesResult) setExchangeByNameKey(parsedExchange map[string]interface{}) {
	ex.Exchanges[parsedExchange["Name"].(string)] = Exchange{
		Name:           parsedExchange["Name"].(string),
		Symbol:         parsedExchange["Symbol"].(string),
		PercentChange:  parsedExchange["PercentChange"].(string),
		ChangeInPoints: parsedExchange["Change"].(string),
		LastTradeDate:  parsedExchange["LastTradeDate"].(string),
		LastTradeTime:  parsedExchange["LastTradeTime"].(string),
	}
}
