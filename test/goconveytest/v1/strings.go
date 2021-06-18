package v1

func StringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	// 两个字符串切片一个是 nil，一个不是 nil
	if (a == nil) != (b == nil) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
