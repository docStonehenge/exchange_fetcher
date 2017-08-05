package slice

import (
	"fmt"
	"github.com/docStonehenge/exchange_fetcher/indices"
)

type StringSlice []string

type StringSliceError struct {
	message string
}

func (slice *StringSlice) Set(value string) error {
	*slice = StringSlice(indices.SplitListBody(value))

	if len(*slice) == 0 {
		return &StringSliceError{"List of values should not be empty."}
	}

	return nil
}

func (slice *StringSlice) String() string {
	return fmt.Sprintf("%s", *slice)
}

func (e *StringSliceError) Error() string {
	return e.message
}
