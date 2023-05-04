package splits

import (
	"strings"
	"testing"
)

func TestStringsSplit(t *testing.T) {
	sep := "[ hh hh"
	s1 := "jskdhjf [ hh hh ss]"
	s2 := "jskdhjf [ hh hh ss] [ hh hh ss]"
	tokens1 := strings.Split(s1, sep)
	t.Log(tokens1[0])
	tokens2 := strings.Split(s2, sep)
	t.Log(tokens2[0])
}
