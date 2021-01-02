package cmap

import (
	"fmt"
	"testing"
)

func TestNewConcurrentMap(t *testing.T) {
	var concurrency int
	var pairRedistributor PairRedistributor
	cm, err := NewConcurrentMap(concurrency, pairRedistributor)
	if err == nil {
		t.Fatalf("No error when new a concurrent map with concurrency %d, but should not be the case!",
			concurrency)
	}
	concurrency = MAX_CONCURRENCY + 1
	concurrency = MAX_CONCURRENCY + 1
	cm, err = NewConcurrentMap(concurrency, pairRedistributor)
	if err == nil {
		t.Fatalf("No error when new a concurrent map with concurrency %d, but should not be the case!",
			concurrency)
	}
	concurrency = 16
	cm, err = NewConcurrentMap(concurrency, pairRedistributor)
	if err != nil {
		t.Fatalf("An error occurs when new a concurrent map: %s (concurrency: %d, pairRedistributor: %#v)",
			err, concurrency, pairRedistributor)
	}
	if cm == nil {
		t.Fatalf("Couldn't a new concurrent map! (concurrency: %d, pairRedistributor: %#v)",
			concurrency, pairRedistributor)
	}
	if cm.Concurrency() != concurrency {
		t.Fatalf("Inconsistent concurrency: expected: %d, actual: %d",
			concurrency, cm.Concurrency())
	}
}

func TestConcurrentMap_Put(t *testing.T) {
	number := 30
	cases := genTestingPairs(number)
	concurrency := 10
	var pairRedistributor PairRedistributor
	cm, _ := NewConcurrentMap(concurrency, pairRedistributor)
	var count uint64
	for _, p := range cases {
		key := p.Key()
		element := p.Element()
		ok, err := cm.Put(key, element)
		if err != nil {
			t.Fatalf("An error occurs when putting a key-element to the cmap: %s (key: %s, element: %#v)",
				err, key, element)
		}
		if !ok {
			t.Fatalf("Couldn't put key-element to the cmap! (key: %s, element: %#v)",
				key, element)
		}
		actualElement := cm.Get(key)
		if actualElement == nil {
			t.Fatalf("Inconsistent element: expected: %#v, actual: %#v",
				element, nil)
		}
		ok, err = cm.Put(key, element)
		if err != nil {
			t.Fatalf("An error occurs when putting a repeated key-element to the cmap! %s (key: %s, element: %#v)",
				err, key, element)
		}
		if ok {
			t.Fatalf("Couldn't put key-element to the cmap! (key: %s, element: %#v)",
				key, element)
		}
		count++
		if cm.Len() != uint64(count) {
			t.Fatalf("Inconsistent size: expected: %d, actual: %d",
				count, cm.Len())
		}
	}
	if cm.Len() != uint64(number) {
		t.Fatalf("Inconsistent size: expected: %d, actual: %d",
			number, cm.Len())
	}
}

func TestConcurrentMap_PutInParallel(t *testing.T) {
	number := 30
	cases := genNoRepetitiveTestingPairs(number)
	concurrency := number / 2
	cm, _ := NewConcurrentMap(concurrency, nil)
	testFunc := func(key string, element interface{}, t *testing.T) func(t *testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			ok, err := cm.Put(key, element)
			if err != nil {
				t.Fatalf("An error occurs when putting a key-element to the cmap: %s (key: %s, element: %#v)",
					err, key, element)
			}
			if !ok {
				t.Fatalf("Couldn't put key-element to the cmap! (key: %s, element: %#v)",
					key, element)
			}
			actualElement := cm.Get(key)
			if actualElement == nil {
				t.Fatalf("Inconsistent element: expected: %#v, actual: %#v",
					element, nil)
			}
			ok, err = cm.Put(key, element)
			if err != nil {
				t.Fatalf("An error occurs when putting a repeated key-element to the cmap! %s (key: %s, element: %#v)",
					err, key, element)
			}
			if ok {
				t.Fatalf("Couldn't put key-element to the cmap! (key: %s, element: %#v)",
					key, element)
			}
		}
	}
	t.Run("Put in parallel", func(t *testing.T) {
		for _, p := range cases {
			t.Run(fmt.Sprintf("Key=%s", p.Key()),
				testFunc(p.Key(), p.Element(), t))
		}
	})
	if cm.Len() != uint64(number) {
		t.Fatalf("Inconsistent size: expected: %d, actual: %d",
			number, cm.Len())
	}
}

func TestConcurrentMap_GetInParallel(t *testing.T) {
	number := 30
	cases := genNoRepetitiveTestingPairs(number)
	concurrency := number / 2
	cm, _ := NewConcurrentMap(concurrency, nil)
	for _, p := range cases {
		cm.Put(p.Key(), p.Element())
	}
	testFunc := func(key string, element interface{}, t *testing.T) func(t *testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			actualElement := cm.Get(key)
			if actualElement == nil {
				t.Fatalf("Inconsistent element: expected: %#v, actual: %#v",
					element, nil)
			}
			if actualElement != element {
				t.Fatalf("Inconsistent element: expected: %#v, actual: %#v",
					element, actualElement)
			}
		}
	}
	t.Run("Get in parallel", func(t *testing.T) {
		t.Run("Put in parallel", func(t *testing.T) {
			for _, p := range cases {
				cm.Put(p.Key(), p.Element())
			}
		})
		for _, p := range cases {
			t.Run(fmt.Sprintf("Get: Key=%s", p.Key()),
				testFunc(p.Key(), p.Element(), t))
		}
	})
	if cm.Len() != uint64(number) {
		t.Fatalf("Inconsistent size: expected: %d, actual: %d",
			number, cm.Len())
	}
}

func TestConcurrentMap_Delete(t *testing.T) {
	number := 30
	cases := genTestingPairs(number)
	concurrency := number / 2
	cm, _ := NewConcurrentMap(concurrency, nil)
	for _, p := range cases {
		cm.Put(p.Key(), p.Element())
	}
	count := uint64(number)
	for _, p := range cases {
		done := cm.Delete(p.Key())
		if !done {
			t.Fatalf("Couldn't delete a key-element from cmap! (key: %s, element: %#v)",
				p.Key(), p.Element())
		}
		actualElement := cm.Get(p.Key())
		if actualElement != nil {
			t.Fatalf("Inconsistent key-element: expected: %#v, actual: %#v",
				nil, actualElement)
		}
		done = cm.Delete(p.Key())
		if done {
			t.Fatalf("Couldn't delete a key-element from cmap again! (key: %s, element: %#v)",
				p.Key(), p.Element())
		}
		if count > 0 {
			count--
		}
		if cm.Len() != count {
			t.Fatalf("Inconsistent size: expected: %d, actual: %d",
				count, cm.Len())
		}
	}
	if cm.Len() != 0 {
		t.Fatalf("Inconsistent size: expected: %d, actual: %d",
			0, cm.Len())
	}
}

func TestConcurrentMap_DeleteInParallel(t *testing.T) {
	number := 30
	cases := genNoRepetitiveTestingPairs(number)
	concurrency := number / 2
	cm, _ := NewConcurrentMap(concurrency, nil)
	for _, p := range cases {
		cm.Put(p.Key(), p.Element())
	}
	testFunc := func(key string, element interface{}, t *testing.T) func(t *testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			done := cm.Delete(key)
			if !done {
				t.Fatalf("Couldn't delete a key-element from cmap! (key: %s, element: %#v)",
					key, element)
			}
			actualElement := cm.Get(key)
			if actualElement != nil {
				t.Fatalf("Inconsistent key-element: expected: %#v, actual: %#v",
					nil, actualElement)
			}
			done = cm.Delete(key)
			if done {
				t.Fatalf("Couldn't delete a key-element from cmap again! (key: %s, element: %#v)",
					key, element)
			}
		}
	}
	t.Run("Delete in parallel", func(t *testing.T) {
		for _, p := range cases {
			t.Run(fmt.Sprintf("Key=%s", p.Key()),
				testFunc(p.Key(), p.Element(), t))
		}
	})
	if cm.Len() != 0 {
		t.Fatalf("Inconsistent size: expected: %d, actual: %d",
			0, cm.Len())
	}
}

var testCaseNumberForCmapTest = 200000
var casesForCmapTest = genNoRepetitiveTestingPairs(testCaseNumberForCmapTest)
var cases1ForCmapTest = casesForCmapTest[:testCaseNumberForCmapTest/2]
var cases2ForCmapTest = casesForCmapTest[testCaseNumberForCmapTest/2:]

func TestConcurrentMap_AllInParallel(t *testing.T) {
	cases1 := cases1ForCmapTest
	cases2 := cases2ForCmapTest
	concurrency := testCaseNumberForCmapTest / 4
	cm, _ := NewConcurrentMap(concurrency, nil)
	t.Run("All in parallel", func(t *testing.T) {
		t.Run("Put1", func(t *testing.T) {
			t.Parallel()
			for _, p := range cases1 {
				_, err := cm.Put(p.Key(), p.Element())
				if err != nil {
					t.Fatalf("An error occurs when putting a key-element to the cmap: %s (key: %s, element: %#v)",
						err, p.Key(), p.Element())
				}
			}
		})
		t.Run("Put2", func(t *testing.T) {
			t.Parallel()
			for _, p := range cases2 {
				_, err := cm.Put(p.Key(), p.Element())
				if err != nil {
					t.Fatalf("An error occurs when putting a key-element to the cmap: %s (key: %s, element: %#v)",
						err, p.Key(), p.Element())
				}
			}
		})
		t.Run("Get1", func(t *testing.T) {
			t.Parallel()
			for _, p := range cases1 {
				actualElement := cm.Get(p.Key())
				if actualElement == nil {
					continue
				}
				if actualElement != p.Element() {
					t.Fatalf("Inconsistent element: expected: %#v, actual: %#v",
						p.Element(), actualElement)
				}
			}
		})
		t.Run("Get2", func(t *testing.T) {
			t.Parallel()
			for _, p := range cases1 {
				actualElement := cm.Get(p.Key())
				if actualElement == nil {
					continue
				}
				if actualElement != p.Element() {
					t.Fatalf("Inconsistent element: expected: %#v, actual: %#v",
						p.Element(), actualElement)
				}
			}
		})
		t.Run("Delete1", func(t *testing.T) {
			t.Parallel()
			for _, p := range cases1 {
				cm.Delete(p.Key())
			}
		})
		t.Run("Delete2", func(t *testing.T) {
			t.Parallel()
			for _, p := range cases2 {
				cm.Delete(p.Key())
			}
		})
	})
}
