package component

// CalculateScore 代表用于计算组件评分的函数类型
type CalculateScore func(counts Counts) uint64
