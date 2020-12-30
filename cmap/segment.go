package cmap

type Segment interface {
	Put(p Pair) (bool, error)
	Get(key string) Pair
}

type SegmentImpl struct {

}
