package slice

import (
	"testing"
)

func TestStringSliceSetSetsSplittedStringOnPointerSlice(t *testing.T) {
	var s StringSlice

	exp := []string{"AAPL", "GOOGL"}

	s.Set("AAPL, GOOGL")

	if length := len(s); length != 2 {
		t.Fatalf("StringSlice pointer should have length on 2, but has length of %v", length)
	}

	for i, v := range s {
		if v != exp[i] {
			t.Fatalf("StringSlice pointer should have values of %v, but have %v", exp, s)
		}
	}
}

func TestStringSliceReturnsStringSliceErrorWhenValueIsEmpty(t *testing.T) {
	var s StringSlice

	err := s.Set("")

	if err == nil {
		t.Fatal("StringSlice.Set should return error on empty slice value")
	}
}

func TestStringSliceStringMethodReturnsArrayRepresentationOfPointer(t *testing.T) {
	var s StringSlice
	s.Set("AAPL, GOOGL")

	if actual := s.String(); actual != "[AAPL GOOGL]" {
		t.Fatalf("StringSlice.String should return array representation of slice values, but returned %v", actual)
	}
}

func TestStringSliceErrorReturnsCorrectMessage(t *testing.T) {
	err := StringSliceError{"Error Test"}

	if m := err.Error(); m != err.message {
		t.Fatalf("StringSliceError.Error should return message from struct point, but returned %v", m)
	}
}
