package v3

import "testing"

func TestEmptyMap(t *testing.T) {
	var map1 map[string]string
	map2 := make(map[string]string)
	map3 := map[string]string{
		"key": "value",
	}

	if map1 == nil {
		t.Log("map1 is nil")
	}

	if len(map1) == 0 {
		t.Log("map1 length is zero")
	}

	if map2 == nil {
		t.Log("map2 is nil")
	}

	if len(map2) == 0 {
		t.Log("map2 length is zero")
	}

	if map3 == nil {
		t.Log("map3 is nil")
	}

	if len(map3) == 0 {
		t.Log("map3 length is zero")
	}
}

// Output:
// map1 is nil
// map1 length is zero
// map2 length is zero
