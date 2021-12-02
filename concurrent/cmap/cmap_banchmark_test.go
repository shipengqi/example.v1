package cmap

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func BenchmarkConcurrentMap_Put(b *testing.B) {
	var number = 20
	var cases = genNoRepetitiveTestingPairs(number)
	concurrency := number / 4
	cm, _ := NewConcurrentMap(concurrency, nil)
	b.ResetTimer()
	for _, tc := range cases {
		key := tc.Key()
		element := tc.Element()
		b.Run(key, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cm.Put(key, element)
			}
		})
	}
}

func BenchmarkConcurrentMap_PutPresent(b *testing.B) {
	var number = 20
	concurrency := number / 4
	cm, _ := NewConcurrentMap(concurrency, nil)
	key := "invariable key"
	b.ResetTimer()
	for i := 0; i < number; i++ {
		element := strconv.Itoa(i)
		b.Run(key, func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				cm.Put(key, element)
			}
		})
	}
}

func BenchmarkMap_Put(b *testing.B) {
	var number = 10
	var cases = genNoRepetitiveTestingPairs(number)
	m := make(map[string]interface{})
	b.ResetTimer()
	for _, tc := range cases {
		key := tc.Key()
		element := tc.Element()
		b.Run(key, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m[key] = element
			}
		})
	}
}

// -- Get -- //

func BenchmarkConcurrentMap_Get(b *testing.B) {
	var number = 100000
	var cases = genNoRepetitiveTestingPairs(number)
	concurrency := number / 4
	cm, _ := NewConcurrentMap(concurrency, nil)
	for _, p := range cases {
		cm.Put(p.Key(), p.Element())
	}
	b.ResetTimer()
	for i := 0; i < 10; i++ {
		key := cases[rand.Intn(number)].Key()
		b.Run(key, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cm.Get(key)
			}
		})
	}
}

func BenchmarkMap_Get(b *testing.B) {
	var number = 100000
	var cases = genNoRepetitiveTestingPairs(number)
	m := make(map[string]interface{})
	for _, p := range cases {
		m[p.Key()] = p.Element()
	}
	b.ResetTimer()
	for i := 0; i < 10; i++ {
		key := cases[rand.Intn(number)].Key()
		b.Run(key, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = m[key]
			}
		})
	}
}

// -- Delete -- //

func BenchmarkConcurrentMap_Delete(b *testing.B) {
	var number = 100000
	var cases = genNoRepetitiveTestingPairs(number)
	concurrency := number / 4
	cm, _ := NewConcurrentMap(concurrency, nil)
	for _, p := range cases {
		cm.Put(p.Key(), p.Element())
	}
	b.ResetTimer()
	for i := 0; i < 20; i++ {
		key := cases[rand.Intn(number)].Key()
		b.Run(key, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cm.Delete(key)
			}
		})
	}
}

func BenchmarkMap_Delete(b *testing.B) {
	var number = 100000
	var cases = genNoRepetitiveTestingPairs(number)
	m := make(map[string]interface{})
	for _, p := range cases {
		m[p.Key()] = p.Element()
	}
	b.ResetTimer()
	for i := 0; i < 20; i++ {
		key := cases[rand.Intn(number)].Key()
		b.Run(key, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				delete(m, key)
			}
		})
	}
}

// -- Len -- //

func BenchmarkConcurrentMap_Len(b *testing.B) {
	var number = 100000
	var cases = genNoRepetitiveTestingPairs(number)
	concurrency := number / 4
	cm, _ := NewConcurrentMap(concurrency, nil)
	for _, p := range cases {
		cm.Put(p.Key(), p.Element())
	}
	b.ResetTimer()
	for i := 0; i < 5; i++ {
		b.Run(fmt.Sprintf("Len%d", i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cm.Len()
			}
		})
	}
}

func BenchmarkMap_Len(b *testing.B) {
	var number = 100000
	var cases = genNoRepetitiveTestingPairs(number)
	m := make(map[string]interface{})
	for _, p := range cases {
		m[p.Key()] = p.Element()
	}
	b.ResetTimer()
	for i := 0; i < 5; i++ {
		b.Run(fmt.Sprintf("Len%d", i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = len(m)
			}
		})
	}
}
