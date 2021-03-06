package cache

import (
	"strconv"
	"strings"
)

const (
	CACHE_KEY_PREFIX_TAG     = "TAG"
)

type Tag struct {
	Name     string
	ID       int
	State    int
	PageNum  int
	PageSize int
}

func (t *Tag) GetTagsCacheKey() string {
	keys := []string{
		CACHE_KEY_PREFIX_TAG,
		"LIST",
	}

	if t.Name != "" {
		keys = append(keys, t.Name)
	}
	if t.State >= 0 {
		keys = append(keys, strconv.Itoa(t.State))
	}
	if t.PageNum > 0 {
		keys = append(keys, strconv.Itoa(t.PageNum))
	}
	if t.PageSize > 0 {
		keys = append(keys, strconv.Itoa(t.PageSize))
	}

	return strings.Join(keys, "_")
}
