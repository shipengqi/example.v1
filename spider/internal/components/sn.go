package component

type SNGenerator interface {
	// Start 用于获取预设的最小序列号
	Start() uint64
	// Max 用于获取预设的最大序列号
	Max() uint64
	// Next 用于获取下一个序列号
	Next() uint64
	// CycleCount 用于获取循环计数
	CycleCount() uint64
	// Get 用于获得一个序列号并准备下一个序列号
	Get() uint64
}
