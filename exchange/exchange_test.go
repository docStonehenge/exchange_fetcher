package exchange

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetch(t *testing.T) {
	jsonResult := "{\"query\":{\"results\":{\"quote\":{\"symbol\":\"^n225\",\"Change_PercentChange\":\"-172.98 - -0.91%\",\"Change\":\"-172.98\"}}}}"

	testServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, jsonResult)
			},
		),
	)

	defer testServer.Close()

	result, _ := Fetch(testServer.URL)

	if result.rawResult != jsonResult {
		t.Fatalf("Raw result fetched should be %v, is %v", jsonResult, result.rawResult)
	}
}
func TestFetchWithRequestError(t *testing.T) {
	mockURL := "foo.bar"

	_, err := Fetch(mockURL)

	if err == nil {
		t.Fatal("Fetching invalid URL should return an error, but nothing happened.")
	}
}

func TestBuildURLWithOneIndex(t *testing.T) {
	index := []string{"^BVSP"}
	expected := "https://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20yahoo.finance.quotes%20where%20symbol%20in%20(%22^BVSP%22)&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys"

	actual := BuildURL(index)

	if actual != expected {
		t.Fatalf("Expected URL to be %s, is %s", expected, actual)
	}
}

func TestBuildURLWithMoreThanOneIndex(t *testing.T) {
	indexes := []string{"^BVSP", "GOOGL"}
	expected := "https://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20yahoo.finance.quotes%20where%20symbol%20in%20(%22^BVSP,GOOGL%22)&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys"

	actual := BuildURL(indexes)

	if actual != expected {
		t.Fatalf("Expected URL to be %s, is %s", expected, actual)
	}
}

func TestParseForOneIndex(t *testing.T) {
	jsonResult := "{\"query\":{\"results\":{\"quote\":{\"Name\":\"Nikkei 225\",\"Symbol\":\"^n225\",\"ChangeInPercent\":\"-0.91%\",\"Change\":\"-172.98\",\"LastTradeDate\":\"4/14/2017\",\"LastTradeTime\":\"3:15pm\"}}}}"

	exchangeResult := ExchangesResult{rawResult: jsonResult}

	exchangeResult.Parse()

	for key, exchange := range exchangeResult.Exchanges {
		if key != "Nikkei 225" {
			t.Fatalf("Exchanges list should have key %s, but it is %s", "Nikkei 225", key)
		}

		if exchange.Name != "Nikkei 225" {
			t.Fatalf("Exchange name should be %s, is %s", "Nikkei 225", exchange.Name)
		}

		if exchange.Symbol != "^n225" {
			t.Fatalf("Exchange symbol should be %s, is %s", "^n225", exchange.Symbol)
		}

		if exchange.PercentChange != "-0.91%" {
			t.Fatalf("Parsed percent change should be %s, is %s", "-0.91%", exchange.PercentChange)
		}

		if exchange.ChangeInPoints != "-172.98" {
			t.Fatalf("Parsed change in points should be %s, is %s", "-172.98", exchange.ChangeInPoints)
		}

		if exchange.LastTradeDate != "4/14/2017" {
			t.Fatalf("Parsed 'last trade date' should be %s, is %s", "4/14/2017", exchange.LastTradeDate)
		}

		if exchange.LastTradeTime != "3:15pm" {
			t.Fatalf("Parsed 'last trade time' should be %s, is %s", "3:15pm", exchange.LastTradeTime)
		}
	}
}

func TestParseForMoreThanOneIndex(t *testing.T) {
	exchangeResult := ExchangesResult{
		rawResult: "{\"query\":{\"results\":{\"quote\":[{\"Name\":\"Nikkei 225\",\"Symbol\":\"^n225\",\"ChangeInPercent\":\"-0.91%\",\"Change\":\"-172.98\",\"LastTradeDate\":\"4/14/2017\",\"LastTradeTime\":\"3:15pm\"},{\"Name\":\"Alphabet Inc.\",\"Symbol\":\"GOOGL\",\"ChangeInPercent\":\"-0.09%\",\"Change\":\"-0.76\",\"LastTradeDate\":\"4/13/2017\",\"LastTradeTime\":\"4:00pm\"}]}}}",
	}

	exchangeResult.Parse()

	expectedList := map[string]Exchange{
		"Nikkei 225":    Exchange{Name: "Nikkei 225", Symbol: "^n225", PercentChange: "-0.91%", ChangeInPoints: "-172.98", LastTradeDate: "4/14/2017", LastTradeTime: "3:15pm"},
		"Alphabet Inc.": Exchange{Name: "Alphabet Inc.", Symbol: "GOOGL", PercentChange: "-0.09%", ChangeInPoints: "-0.76", LastTradeDate: "4/13/2017", LastTradeTime: "4:00pm"},
	}

	nikkei := exchangeResult.Exchanges["Nikkei 225"]

	if nikkei != expectedList["Nikkei 225"] {
		t.Fatalf("Parsed exchanges list should have a %s exchange and it should be equal to %v, but it is %v", "Nikkei 225", expectedList["Nikkei 225"], nikkei)
	}

	google := exchangeResult.Exchanges["Alphabet Inc."]

	if google != expectedList["Alphabet Inc."] {
		t.Fatalf("Parsed exchanges list should have a %s exchange and it should be equal to %v, but it is %v", "Alphabet Inc.", expectedList["Alphabet Inc."], google)
	}
}
