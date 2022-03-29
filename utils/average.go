package utils

func AverageV1(low, high uint) uint {
	return low + (high-low)/2
}

// AverageV2016 Average 2016
// 对两个无符号整数进行除法，同时通过 按位与 修正低位数字，保证两个证书都为奇数时，仍然正确
func AverageV2016(a, b uint) uint {
	return (a / 2) + (b / 2) + (a & b & 1)
}
