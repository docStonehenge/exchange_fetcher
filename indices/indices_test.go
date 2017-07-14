package indices

import (
	"github.com/docStonehenge/exchange_fetcher/exchange"
	"strings"
	"testing"
)

func TestSplitCorrectlySeparatesIndicesOnJSONstring(t *testing.T) {
	results := []struct {
		body string
		exp  []string
	}{
		{body: "{\"indices\":[\"AAPL\", \"GOOGL\"]}", exp: []string{"AAPL", "GOOGL"}},
		{body: "{\"indices\":[\"MGLU3.SA\"]}", exp: []string{"MGLU3.SA"}},
		{body: "{\"indices\":[]}", exp: []string{}},
	}

	for _, r := range results {
		if actual := Split([]byte(r.body)); strings.Join(actual, ", ") != strings.Join(r.exp, ", ") {
			t.Fatalf("indices.Split should separate indices correctly, but result was %v", actual)
		}
	}
}

func TestSplitReturnsEmptyArrayWhenArgumentIsEmpty(t *testing.T) {
	results := []struct {
		body string
		exp  []string
	}{
		{body: "{}", exp: []string{}},
		{body: "", exp: []string{}},
	}

	for _, r := range results {
		if actual := Split([]byte(r.body)); strings.Join(actual, ", ") != strings.Join(r.exp, ", ") {
			t.Fatalf("indices.Split should return an empty array, but result was %v", actual)
		}
	}
}

func TestJoinReturnsParsedJSONExchanges(t *testing.T) {
	exchanges := make(map[string]exchange.Exchange)

	exp := "{\"Bar\":{\"Name\":\"Bar\",\"Symbol\":\"B\",\"PercentChange\":\"2%\",\"ChangeInPoints\":\"2.0\",\"LastTradeDate\":\"12/01/2017\",\"LastTradeTime\":\"12:31pm\"},\"Foo\":{\"Name\":\"Foo\",\"Symbol\":\"F\",\"PercentChange\":\"2%\",\"ChangeInPoints\":\"2.0\",\"LastTradeDate\":\"12/01/2017\",\"LastTradeTime\":\"12:31pm\"}}"

	exchanges["Foo"] = exchange.Exchange{
		Name:           "Foo",
		Symbol:         "F",
		PercentChange:  "2%",
		ChangeInPoints: "2.0",
		LastTradeDate:  "12/01/2017",
		LastTradeTime:  "12:31pm",
	}

	exchanges["Bar"] = exchange.Exchange{
		Name:           "Bar",
		Symbol:         "B",
		PercentChange:  "2%",
		ChangeInPoints: "2.0",
		LastTradeDate:  "12/01/2017",
		LastTradeTime:  "12:31pm",
	}

	jsonBody, _ := Join(exchanges)

	if string(jsonBody) != exp {
		t.Fatalf("Built JSON response should be equal to %v, but is %v", exp, string(jsonBody))
	}
}
