package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testify assert 和 require 都提供了断言的功能
// 这两个包的唯一区别就是：require 的函数断言失败会直接导致单元测试结束
// assert 只会标记为 case 失败，但单元测试会继续往下执行剩下的 case。

func TestCase1(t *testing.T) {
	name := "Bob"
	age := 10

	assert.Equal(t, "bob", name)
	assert.Equal(t, 20, age)
}

// Output: === RUN   TestCase1
//    TestCase1: require_test.go:18:
//        	Error Trace:
//        	Error:      	Not equal:
//        	            	expected: "bob"
//        	            	actual  : "Bob"
//
//        	            	Diff:
//        	            	--- Expected
//        	            	+++ Actual
//        	            	@@ -1 +1 @@
//        	            	-bob
//        	            	+Bob
//        	Test:       	TestCase1
//    TestCase1: require_test.go:19:
//        	Error Trace:
//        	Error:      	Not equal:
//        	            	expected: 20
//        	            	actual  : 10
//        	Test:       	TestCase1
//--- FAIL: TestCase1 (0.00s)
//
//
//Expected :20
//Actual   :10
//<Click to see difference>
//
//
//FAIL
//
//Process finished with exit code 1
// 两个 assert.Equal() 指令都被执行了



func TestCase2(t *testing.T) {
	name := "Bob"
	age := 10

	require.Equal(t, "bob", name)
	require.Equal(t, 20, age)
}

// Output: === RUN   TestCase2
//    TestCase2: require_test.go:61:
//        	Error Trace:
//        	Error:      	Not equal:
//        	            	expected: "bob"
//        	            	actual  : "Bob"
//
//        	            	Diff:
//        	            	--- Expected
//        	            	+++ Actual
//        	            	@@ -1 +1 @@
//        	            	-bob
//        	            	+Bob
//        	Test:       	TestCase2
//--- FAIL: TestCase2 (0.00s)
//
//
//Expected :bob
//Actual   :Bob
//<Click to see difference>
//
//
//FAIL
//
//
//Process finished with exit code 1
// 只有第一个 require.Equal() 指令被执行了，第二个 require.Equal() 没有被执行。
