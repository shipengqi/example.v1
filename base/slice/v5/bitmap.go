package main

// BitMap implement a bitmap
type BitMap struct {
	Bits     []byte `json:"bits"`
	Count    uint   `json:"count"`
	Capacity uint   `json:"cap"`
}

// NewBitMap Create a new BitMap
func NewBitMap(cap uint) *BitMap {
	bits := make([]byte, (cap>>3)+1)
	return &BitMap{Bits: bits, Capacity: cap, Count: 0}
}
