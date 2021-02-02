package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXOR(t *testing.T) {
	key := "signature"
	cases := []struct {
		input    string
		expected string
	}{
		{"case1", "2d20f1cd235afc9538202b25f2e68428"},
		{"case2", "4abfb640cb9e63d4cb42b3f7f98a9796"},
		{"case3", "0027998a647175eb55a587edf515aef0"},
	}

	for i := range cases {
		en := EncodeXOR(cases[i].input, key)
		t.Logf("input: %s, encrypted: %s", cases[i].input, en)
		de, err := DecodeXOR(en, key)
		assert.Equal(t, nil, err)
		t.Logf("input: %s, decrypted: %s", cases[i].input, de)
		assert.Equal(t, cases[i].input, de)
	}
}