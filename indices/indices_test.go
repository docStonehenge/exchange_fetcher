package indices

import (
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
