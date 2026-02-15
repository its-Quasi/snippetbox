package assert

import (
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper() // dont show the error here, show the error where I was call
	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}
