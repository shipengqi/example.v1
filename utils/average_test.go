package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAverage(t *testing.T) {

	tests := []struct {
		input1   uint
		input2   uint
		expected int
	}{
		{1333, 1335, 1334},
		{6, 9, 7},
	}

	for _, v := range tests {
		got := AverageV2016(v.input1, v.input2)
		assert.Equal(t, v.expected, int(got))
	}
	for _, v := range tests {
		got := AverageV1(v.input1, v.input2)
		assert.Equal(t, v.expected, int(got))
	}

}
