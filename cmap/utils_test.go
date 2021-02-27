package cmap

import "testing"

func TestBKDRHash(t *testing.T) {
	input := "hello"
	var expected uint64 = 3092122322284475622
	actual := BKDRHash(input)
	if actual != expected {
		t.Fatalf("BKDRHash: expected: %#v, actual: %#v",
			expected, actual)
	}
}