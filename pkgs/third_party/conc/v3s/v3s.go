package v2s

import (
	"github.com/sourcegraph/conc/iter"
)

func process(values []int) {
	iter.ForEach(values, handle)
}

func handle(e *int) {

}
