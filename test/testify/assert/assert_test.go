package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// https://github.com/stretchr/testify
// Testify 可以帮助你简化在测试用例中编写断言的方式。
// Testify 还可用于模拟测试框架中的对象，以确保你在测试时不会调用产品端。
func TestCalculate(t *testing.T) {
	if Calculate(2) != 4 {
		t.Error("Expected 2 + 2 to equal 4")
	}
}

func TestCalculate2(t *testing.T) {
	assert.Equal(t, Calculate(2), 4)
}

func TestCalculate3(t *testing.T) {
	// assert.New(t) 初始化断言，然后可以多次调用 assert.Equal()，只需传入输入值和期望值，
	// 而不是每次都将 t 作为第一个参数传入。
	assert := assert.New(t)

	var tests = []struct {
		input    int
		expected int
	}{
		{2, 4},
		{-1, 1},
		{0, 2},
		{-5, -3},
		{99999, 100001},
	}

	for _, test := range tests {
		assert.Equal(Calculate(test.input), test.expected)
	}
}
