package v2

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFib(t *testing.T) {
	var cases = []struct {
		in       int // input
		expected int // expected result
	}{
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 3},
		{5, 5},
		{6, 8},
		{7, 13},
	}
	for _, c := range cases {
		actual := Fib(c.in)
		assert.Equal(t, c.expected, actual)
	}
}

func TestFib_Parallel(t *testing.T) {
	var cases = []struct {
		name     string
		in       int // input
		expected int // expected result
	}{
		{"1 的 Fib", 1, 1},
		{"2 的 Fib", 2, 1},
		{"3 的 Fib", 3, 2},
		{"4 的 Fib", 4, 3},
		{"5 的 Fib", 5, 5},
		{"6 的 Fib", 6, 8},
		{"7 的 Fib", 7, 13},
	}
	for _, c := range cases {
		cc := c
		t.Run(c.name, func(t *testing.T) {
			t.Log("time:", time.Now())
			t.Parallel()
			time.Sleep(1 * time.Second)
			actual := Fib(cc.in)
			assert.Equal(t, cc.expected, actual)
		})
	}
}