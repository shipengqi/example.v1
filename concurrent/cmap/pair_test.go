package cmap

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"
)

type keyElement struct {
	key     string
	element interface{}
}

func TestPair_New(t *testing.T) {
	cases := genTestKeyElementSlice(100)
	for _, c := range cases {
		t.Run(fmt.Sprintf("Key=%s,Element=%#v", c.key, c.element),
			func(t *testing.T) {
				p, err := newPair(c.key, c.element)
				if err != nil {
					t.Fatalf("An error occurs when new a pair: %s (key: %s, element: %#v)",
						err, c.key, c.element)
				}
				if p == nil {
					t.Fatalf("Couldn't new pair! (key: %s, element: %#v)",
						c.key, c.element)
				}
			})
	}
	nilkey := randString()
	t.Run(fmt.Sprintf("Key=%s,Element=nil error", nilkey),
		func(t *testing.T) {
			_, err := newPair(nilkey, nil)
			if err == nil {
				t.Fatalf("Should occurs an error: element is nil")
			}
		})
}

func TestPair_KeyHashAndElement(t *testing.T) {
	cases := genTestKeyElementSlice(50)
	for _, c := range cases {
		t.Run(fmt.Sprintf("Key=%s,Element=%#v", c.key, c.element),
			func(t *testing.T) {
				p, err := newPair(c.key, c.element)
				if err != nil {
					t.Fatalf("An error occurs when new a pair: %s (key: %s, element: %#v)",
						err, c.key, c.element)
				}
				if p == nil {
					t.Fatalf("Couldn't new pair! (key: %s, element: %#v)",
						c.key, c.element)
				}
				if p.Key() != c.key {
					t.Fatalf("Inconsistent key: expected: %s, actual: %s",
						c.key, p.Key())
				}
				expectedHash := hash(c.key)
				if p.Hash() != expectedHash {
					t.Fatalf("Inconsistent hash: expected: %d, actual: %d",
						expectedHash, p.Hash())
				}
				if p.Element() != c.element {
					t.Fatalf("Inconsistent element: expected: %#v, actual: %#v",
						c.element, p.Element())
				}
			})
	}
}

func TestPair_Set(t *testing.T) {
	cases := genTestKeyElementSlice(50)
	for _, c := range cases {
		t.Run(fmt.Sprintf("Key=%s,Element=%#v", c.key, c.element),
			func(t *testing.T) {
				p, err := newPair(c.key, c.element)
				if err != nil {
					t.Fatalf("An error occurs when new a pair: %s (key: %s, element: %#v)",
						err, c.key, c.element)
				}
				newElement := randString()
				_ = p.SetElement(newElement)
				if p.Element() != newElement {
					t.Fatalf("Inconsistent element: expected: %#v, actual: %#v",
						newElement, p.Element())
				}
			})
	}
}

func TestPair_Next(t *testing.T) {
	number := 50
	cases := genTestKeyElementSlice(number)
	var current Pair
	var prev Pair
	var err error
	for _, c := range cases {
		current, err = newPair(c.key, c.element)
		if err != nil {
			t.Fatalf("An error occurs when new a pair: %s (key: %s, element: %#v)",
				err, c.key, c.element)
		}
		if prev != nil {
			_ = current.SetNext(prev)
		}
		prev = current
	}
	for i := number - 1; i >= 0; i-- {
		next := current.Next()
		if i == 0 {
			if next != nil {
				t.Fatalf("Next is not nil! (pair: %#v, index: %d)",
					current, i)
			}
		} else {
			if next == nil {
				t.Fatalf("Next is nil! (pair: %#v, index: %d)",
					current, i)
			}
			expectedNext := cases[i-1]
			if next.Key() != expectedNext.key {
				t.Fatalf("Inconsistent next key: expected: %s, actual: %s, index: %d",
					expectedNext.key, next.Key(), i)
			}
			if next.Element() != expectedNext.element {
				t.Fatalf("Inconsistent element: expected: %#v, actual: %#v, index: %d",
					expectedNext.element, next.Element(), i)
			}
		}
		current = next
	}
}

func TestPair_Copy(t *testing.T) {
	cases := genTestKeyElementSlice(50)
	for _, c := range cases {
		t.Run(fmt.Sprintf("Key=%s,Element=%#v", c.key, c.element),
			func(t *testing.T) {
				p, err := newPair(c.key, c.element)
				if err != nil {
					t.Fatalf("An error occurs when new a pair: %s (key: %s, element: %#v)",
						err, c.key, c.element)
				}
				pCopy := p.Copy()
				if pCopy.Key() != p.Key() {
					t.Fatalf("Inconsistent key: expected: %s, actual: %s",
						p.Key(), pCopy.Key())
				}
				if pCopy.Hash() != p.Hash() {
					t.Fatalf("Inconsistent hash: expected: %d, actual: %d",
						p.Hash(), pCopy.Hash())
				}
				if pCopy.Element() != p.Element() {
					t.Fatalf("Inconsistent element: expected: %#v, actual: %#v",
						p.Element(), pCopy.Element())
				}
			})
	}
}

func genTestKeyElementSlice(number int) []*keyElement {
	cases := make([]*keyElement, number)
	for i := 0; i < number; i++ {
		cases[i] = &keyElement{
			key:     randString(),
			element: randElement(),
		}
	}
	return cases
}

// randElement 会生成并返回一个伪随机元素值。
func randElement() interface{} {
	// 调多次 rand.Seed(x)，但每次 x 保证不一样：每次使用纳秒级别的种子。
	// 强烈不推荐这种，因为高并发的情况下纳秒也可能重复。
	// rand.Seed(time.Now().UnixNano())
	if i := rand.Int31(); i%3 != 0 {
		return i
	}
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, rand.Int31())
	return hex.EncodeToString(buf.Bytes())
}

// randString 会生成并返回一个伪随机字符串。
func randString() string {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, rand.Int31())
	return hex.EncodeToString(buf.Bytes())
}
