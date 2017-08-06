package indices

import (
	"github.com/docStonehenge/exchange_fetcher/exchange"
	"strings"
	"testing"
)

func TestSplitJSONBodyCorrectlySeparatesIndicesOnJSONstring(t *testing.T) {
	results := []struct {
		body string
		exp  []string
	}{
		{body: "{\"indices\":[\"AAPL\", \"GOOGL\"]}", exp: []string{"AAPL", "GOOGL"}},
		{body: "{\"indices\":[\"MGLU3.SA\"]}", exp: []string{"MGLU3.SA"}},
		{body: "{\"indices\":[]}", exp: []string{}},
	}

	for _, r := range results {
		if actual := SplitJSONBody([]byte(r.body)); strings.Join(actual, ", ") != strings.Join(r.exp, ", ") {
			t.Fatalf("indices.SplitJSONBody should separate indices correctly, but result was %v", actual)
		}
	}
}

func TestSplitJSONBodyReturnsEmptyArrayWhenArgumentIsEmpty(t *testing.T) {
	results := []struct {
		body string
		exp  []string
	}{
		{body: "{}", exp: []string{}},
		{body: "", exp: []string{}},
	}

	for _, r := range results {
		if actual := SplitJSONBody([]byte(r.body)); strings.Join(actual, ", ") != strings.Join(r.exp, ", ") {
			t.Fatalf("indices.SplitAsJSON should return an empty array, but result was %v", actual)
		}
	}
}

func TestSplitListBodyCorrectlySeparatesIndicesOnString(t *testing.T) {
	results := []struct {
		body string
		exp  []string
	}{
		{body: "AAPL, GOOGL", exp: []string{"AAPL", "GOOGL"}},
		{body: "MGLU3.SA", exp: []string{"MGLU3.SA"}},
		{body: "MGLU3.SA,BBSE3.SA", exp: []string{"MGLU3.SA", "BBSE3.SA"}},
		{body: "MGLU3.SA,BBSE3.SA,  AAPL, GOOGL", exp: []string{"MGLU3.SA", "BBSE3.SA", "AAPL", "GOOGL"}},
		{body: "MGLU3.SA,BBSE3.SA,  AAPL; GOOGL;TPIS3.SA", exp: []string{"MGLU3.SA", "BBSE3.SA", "AAPL", "GOOGL", "TPIS3.SA"}},
		{body: "", exp: []string{}},
	}

	for _, r := range results {
		actual := SplitListBody(r.body)

		if len(actual) != len(r.exp) {
			t.Fatalf("indices.SplitListBody should return an array with %d elements, but its size is %d", len(r.exp), len(actual))
		}

		for count, index := range actual {
			if index != r.exp[count] {
				t.Fatalf("indices.SplitListBody should return an array with element as %v, but element is %v", r.exp[count], index)
			}
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
