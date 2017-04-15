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
