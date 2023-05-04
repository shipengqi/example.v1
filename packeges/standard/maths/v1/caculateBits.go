package v1

func CalculateStrBits(str string) int {
	return len([]byte(str))*8
}
