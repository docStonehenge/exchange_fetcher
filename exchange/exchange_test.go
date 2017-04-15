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
	expected := "https://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20yahoo.finance.quotes%20where%20symbol%20in%20(%22^BVSP%22)&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys&callback="

	actual := BuildURL(index)

	if actual != expected {
		t.Fatalf("Expected URL to be %s, is %s", expected, actual)
	}
}

func TestBuildURLWithMoreThanOneIndex(t *testing.T) {
	indexes := []string{"^BVSP", "GOOGL"}
	expected := "https://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20yahoo.finance.quotes%20where%20symbol%20in%20(%22^BVSP,GOOGL%22)&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys&callback="

	actual := BuildURL(indexes)

	if actual != expected {
		t.Fatalf("Expected URL to be %s, is %s", expected, actual)
	}
}
