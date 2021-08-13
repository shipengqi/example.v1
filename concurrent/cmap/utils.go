package cmap

// hash 用于计算给定字符串的哈希值的整数形式
// BKDR 哈希算法
func hash(str string) uint64 {
	seed := uint64(13131)
	var hash uint64
	for i := 0; i < len(str); i++ {
		hash = hash*seed + uint64(str[i])
	}
	return hash & 0x7FFFFFFFFFFFFFFF
}

// hash 用于计算给定字符串的哈希值的整数形式
// func hash(str string) uint64 {
// 	h := md5.Sum([]byte(str))
// 	var num uint64
// 	binary.Read(bytes.NewReader(h[:]), binary.LittleEndian, &num)
// 	return num
// }


// for test
func BKDRHash(str string) uint64 {
	seed := uint64(13131)
	var hash uint64
	for i := 0; i < len(str); i++ {
		// 遍历字符串输出的是 unicode 编码，str[i]
		hash = hash*seed + uint64(str[i])
	}
	return hash & 0x7FFFFFFFFFFFFFFF
}

// BKDR 哈希算法推导
// 由一个字符串（比方：ad）得到其哈希值。为了降低碰撞，应该使该字符串中每一个字符都參与哈希值计算。使其符合雪崩效应，也就是说即使改变字符串中的
// 一个字节，也会对终于的哈希值造成较大的影响。
// 最直接的办法就是让字符串中的每一个字符相加。得到其和 SUM，让 SUM 作为哈希值，如 SUM（ad）= a+d;
// 但是依据 ascii 码表得知 a(97)+d(100)=b(98)+c(99)。那么发生了碰撞，也就是直接求和的话会非常容易发生碰撞，
// 所以要对字符间的差距进行放大，乘以一个系数。