package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeMD5(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"case1", "2d20f1cd235afc9538202b25f2e68428"},
		{"case2", "4abfb640cb9e63d4cb42b3f7f98a9796"},
		{"case3", "0027998a647175eb55a587edf515aef0"},
	}

	for i := range cases {
		res := EncodeMD5(cases[i].input)
		assert.Equal(t, cases[i].expected, res)
	}

}

func TestEncodeMD5WithSalt(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"case1", "d62c93b1d441dfdbd104f33270776be8"},
		{"case2", "b038fe94dcf230f3755d9429cd91be96"},
		{"case3", "b6289549d102057ba0ce7194cd9b1509"},
		{"Admin@111", "432fedceddbfcab540c809fbff4e616f"},
	}

	for i := range cases {
		res := EncodeMD5WithSalt(cases[i].input, "llsfhfhhf$jjfklsjn52522@@44ddddsdfsiwotpvbnusf")
		assert.Equal(t, cases[i].expected, res)
	}

}
