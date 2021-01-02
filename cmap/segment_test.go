package cmap

import (
	"fmt"
	"testing"
)

func TestSegment_New(t *testing.T) {
	s := newSegment(-1, nil)
	if s == nil {
		t.Fatal("Couldn't new segment!")
	}
}

func TestSegment_Put(t *testing.T) {
	number := 30
	cases := genTestingPairs(number)
	s := newSegment(-1, nil)
	var count uint64
	for _, p := range cases {
		ok, err := s.Put(p)
		if err != nil {
			t.Fatalf("An error occurs when putting a pair to the segment: %s (pair: %#v)",
				err, p)
		}
		if !ok {
			t.Fatalf("Couldn't put pair to the segment! (pair: %#v)",
				p)
		}
		actualPair := s.Get(p.Key())
		if actualPair == nil {
			t.Fatalf("Inconsistent pair: expected: %#v, actual: %#v",
				p.Element(), nil)
		}
		ok, err = s.Put(p)
		if err != nil {
			t.Fatalf("An error occurs when putting a repeated pair to the segment: %s (pair: %#v)",
				err, p)
		}
		if ok {
			t.Fatalf("Couldn't put repeated pair to the segment! (pair: %#v)",
				p)
		}
		count++
		if s.Size() != count {
			t.Fatalf("Inconsistent size: expected: %d, actual: %d",
				count, s.Size())
		}
	}
	if s.Size() != uint64(number) {
		t.Fatalf("Inconsistent size: expected: %d, actual: %d",
			number, s.Size())
	}
}

func TestSegment_PutInParallel(t *testing.T) {
	number := 30
	cases := genNoRepetitiveTestingPairs(number)
	s := newSegment(-1, nil)
	testFunc := func(p Pair, t *testing.T) func(t *testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			ok, err := s.Put(p)
			if err != nil {
				t.Fatalf("An error occurs when putting a pair to the segment: %s (pair: %#v)",
					err, p)
			}
			if !ok {
				t.Fatalf("Couldn't put a pair to the segment! (pair: %#v)", p)
			}
			actualPair := s.Get(p.Key())
			if actualPair == nil {
				t.Fatalf("Inconsistent pair: expected: %#v, actual: %#v",
					p.Element(), nil)
			}
			ok, err = s.Put(p)
			if err != nil {
				t.Fatalf("An error occurs when putting a repeated pair to the segment: %s (pair: %#v)",
					err, p)
			}
			if ok {
				t.Fatalf("Couldn't put a repeated pair to the segment! (pair: %#v)", p)
			}
		}
	}
	t.Run("Put in parallel", func(t *testing.T) {
		for _, p := range cases {
			t.Run(fmt.Sprintf("Key=%s", p.Key()), testFunc(p, t))
		}
	})
	if s.Size() != uint64(number) {
		t.Fatalf("Inconsistent size: expected: %d, actual: %d",
			number, s.Size())
	}
}

func TestSegment_GetInParallel(t *testing.T) {
	number := 30
	cases := genNoRepetitiveTestingPairs(number)
	s := newSegment(-1, nil)
	for _, p := range cases {
		s.Put(p)
	}
	testFunc := func(p Pair, t *testing.T) func(t *testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			actualPair := s.Get(p.Key())
			if actualPair == nil {
				t.Fatalf("Not found pair in segment! (key: %s)",
					p.Key())
			}
			if actualPair.Key() != p.Key() {
				t.Fatalf("Inconsistent key: expected: %s, actual: %s",
					p.Key(), actualPair.Key())
			}
			if actualPair.Hash() != p.Hash() {
				t.Fatalf("Inconsistent hash: expected: %d, actual: %d",
					p.Hash(), actualPair.Hash())
			}
			if actualPair.Element() != p.Element() {
				t.Fatalf("Inconsistent element: expected: %#v, actual: %#v",
					p.Element(), actualPair.Element())
			}
		}
	}
	t.Run("Get in parallel", func(t *testing.T) {
		t.Run("Put in parallel", func(t *testing.T) {
			for _, p := range cases {
				s.Put(p)
			}
		})
		for _, p := range cases {
			t.Run(fmt.Sprintf("Get: Key=%s", p.Key()), testFunc(p, t))
		}
	})
	if s.Size() != uint64(number) {
		t.Fatalf("Inconsistent size: expected: %d, actual: %d",
			number, s.Size())
	}
}

func TestSegment_Delete(t *testing.T) {
	number := 30
	cases := genTestingPairs(number)
	s := newSegment(-1, nil)
	for _, p := range cases {
		s.Put(p)
	}
	count := uint64(number)
	for _, p := range cases {
		done := s.Delete(p.Key())
		if !done {
			t.Fatalf("Couldn't delete a pair from segment! (pair: %#v)", p)
		}
		actualPair := s.Get(p.Key())
		if actualPair != nil {
			t.Fatalf("Inconsistent pair: expected: %#v, actual: %#v",
				nil, actualPair)
		}
		done = s.Delete(p.Key())
		if done {
			t.Fatalf("Couldn't delete a pair from segment again! (pair: %#v)", p)
		}
		if count > 0 {
			count--
		}
		if s.Size() != count {
			t.Fatalf("Inconsistent size: expected: %d, actual: %d",
				count, s.Size())
		}
	}
	if s.Size() != 0 {
		t.Fatalf("Inconsistent size: expected: %d, actual: %d",
			0, s.Size())
	}
}

func TestSegment_DeleteInParallel(t *testing.T) {
	number := 30
	cases := genNoRepetitiveTestingPairs(number)
	s := newSegment(-1, nil)
	for _, p := range cases {
		s.Put(p)
	}
	testFunc := func(p Pair, t *testing.T) func(t *testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			done := s.Delete(p.Key())
			if !done {
				t.Fatalf("Couldn't delete a pair from segment! (pair: %#v)", p)
			}
			actualPair := s.Get(p.Key())
			if actualPair != nil {
				t.Fatalf("Inconsistent pair: expected: %#v, actual: %#v",
					nil, actualPair)
			}
			done = s.Delete(p.Key())
			if done {
				t.Fatalf("Couldn't delete a pair from segment again! (pair: %#v)", p)
			}
		}
	}
	t.Run("Delete in parallel", func(t *testing.T) {
		for _, p := range cases {
			t.Run(fmt.Sprintf("Key=%s", p.Key()), testFunc(p, t))
		}
	})
	if s.Size() != 0 {
		t.Fatalf("Inconsistent size: expected: %d, actual: %d",
			0, s.Size())
	}
}

var testCaseNumberForSegmentTest = 200000
var casesForSegmentTest = genNoRepetitiveTestingPairs(testCaseNumberForSegmentTest)
var cases1ForSegmentTest = casesForSegmentTest[:testCaseNumberForSegmentTest/2]
var cases2ForSegmentTest = casesForSegmentTest[testCaseNumberForSegmentTest/2:]

func TestSegment_AllInParallel(t *testing.T) {
	cases1 := cases1ForSegmentTest
	cases2 := cases2ForSegmentTest
	s := newSegment(-1, nil)
	t.Run("All in parallel", func(t *testing.T) {
		t.Run("Put1", func(t *testing.T) {
			t.Parallel()
			for _, p := range cases1 {
				_, err := s.Put(p)
				if err != nil {
					t.Fatalf("An error occurs when putting a pair to the segment: %s (pair: %#v)",
						err, p)
				}
			}
		})
		t.Run("Put2", func(t *testing.T) {
			t.Parallel()
			for _, p := range cases2 {
				_, err := s.Put(p)
				if err != nil {
					t.Fatalf("An error occurs when putting a pair to the segment: %s (pair: %#v)",
						err, p)
				}
			}
		})
		t.Run("Get1", func(t *testing.T) {
			t.Parallel()
			for _, p := range cases1 {
				actualPair := s.Get(p.Key())
				if actualPair == nil {
					continue
				}
				if actualPair.Key() != p.Key() {
					t.Fatalf("Inconsistent key: expected: %s, actual: %s",
						p.Key(), actualPair.Key())
				}
				if actualPair.Hash() != p.Hash() {
					t.Fatalf("Inconsistent hash: expected: %d, actual: %d (key=%s)",
						p.Hash(), actualPair.Hash(), p.Key())
				}
				if actualPair.Element() != p.Element() {
					t.Fatalf("Inconsistent element: expected: %#v, actual: %#v (key=%s)",
						p.Element(), actualPair.Element(), p.Key())
				}
			}
		})
		t.Run("Get2", func(t *testing.T) {
			t.Parallel()
			for _, p := range cases2 {
				actualPair := s.Get(p.Key())
				if actualPair == nil {
					continue
				}
				if actualPair.Key() != p.Key() {
					t.Fatalf("Inconsistent key: expected: %s, actual: %s",
						p.Key(), actualPair.Key())
				}
				if actualPair.Hash() != p.Hash() {
					t.Fatalf("Inconsistent hash: expected: %d, actual: %d (key=%s)",
						p.Hash(), actualPair.Hash(), p.Key())
				}
				if actualPair.Element() != p.Element() {
					t.Fatalf("Inconsistent element: expected: %#v, actual: %#v (key=%s)",
						p.Element(), actualPair.Element(), p.Key())
				}
			}
		})
		t.Run("Delete1", func(t *testing.T) {
			t.Parallel()
			for _, p := range cases1 {
				s.Delete(p.Key())
			}
		})
		t.Run("Delete2", func(t *testing.T) {
			t.Parallel()
			for _, p := range cases2 {
				s.Delete(p.Key())
			}
		})
	})
}