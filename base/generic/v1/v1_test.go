package v1

import (
	"fmt"
	"testing"
)

// Initialize a map for the integer values
var ints = map[string]int64{
	"first":  34,
	"second": 12,
}

// Initialize a map for the float values
var floats = map[string]float64{
	"first":  35.98,
	"second": 26.99,
}

func TestNonGeneric(t *testing.T) {
	fmt.Printf("Non-Generic Sums: %v and %v\n", SumInts(ints), SumFloats(floats))
}

func TestSumIntsOrFloats(t *testing.T) {

	fmt.Printf("Generic Sums: %v and %v\n", SumIntsOrFloats[string, int64](ints), SumIntsOrFloats[string, float64](floats))
	fmt.Printf("Generic Sums, type parameters inferred: %v and %v\n", SumIntsOrFloats(ints), SumIntsOrFloats(floats))
}

func TestSumNumbers(t *testing.T) {
	fmt.Printf("Generic Sums with Constraint: %v and %v\n", SumNumbers(ints), SumNumbers(floats))
}
