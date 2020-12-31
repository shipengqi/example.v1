package cmap

import "sync/atomic"

// BucketStatus 代表散列桶状态的类型
type BucketStatus uint8

const (
	// BUCKET_STATUS_NORMAL 代表散列桶正常
	BUCKET_STATUS_NORMAL BucketStatus = iota
	// BUCKET_STATUS_UNDERWEIGHT 代表散列桶过轻
	BUCKET_STATUS_UNDERWEIGHT
	// BUCKET_STATUS_OVERWEIGHT 代表散列桶过重
	BUCKET_STATUS_OVERWEIGHT
)

// bucketCountTemplate 代表调试用散列桶状态信息模板
var bucketCountTemplate = `Bucket count: 
    pairTotal: %d
    bucketNumber: %d
    average: %f
    upperThreshold: %d
    emptyBucketCount: %d
`

// PairRedistributor 代表针对键-元素对的再分布器
// 用于当散列段内的键-元素对分布不均时进行重新分布
type PairRedistributor interface {
	//  UpdateThreshold 会根据键-元素对总数和散列桶总数计算并更新阈值
	UpdateThreshold(pairTotal uint64, bucketNumber int)
	// CheckBucketStatus 用于检查散列桶的状态。
	CheckBucketStatus(pairTotal uint64, bucketSize uint64) (bucketStatus BucketStatus)
	// Redistribe 用于实施键-元素对的再分布。
	Redistribe(bucketStatus BucketStatus, buckets []Bucket) (newBuckets []Bucket, changed bool)
}


type PairRedistributorImpl struct {
	// loadFactor 代表装载因子
	loadFactor float64
	// upperThreshold 代表散列桶重量的上阈限
	// 当某个散列桶的尺寸增至此值时会触发再散列
	upperThreshold uint64
	// overweightBucketCount 代表过重的散列桶的计数
	overweightBucketCount uint64
	// emptyBucketCount 代表空的散列桶的计数
	emptyBucketCount uint64
}

// newDefaultPairRedistributor 会创建一个PairRedistributor类型的实例
// 参数 loadFactor 代表散列桶的负载因子。
// 参数 bucketNumber 代表散列桶的数量。
func newDefaultPairRedistributor(loadFactor float64, bucketNumber int) PairRedistributor {
	if loadFactor <= 0 {
		loadFactor = DEFAULT_BUCKET_LOAD_FACTOR
	}
	pr := &PairRedistributorImpl{}
	pr.loadFactor = loadFactor
	pr.UpdateThreshold(0, bucketNumber)
	return pr
}

func (pr *PairRedistributorImpl) UpdateThreshold(pairTotal uint64, bucketNumber int) {
	var average float64
	average = float64(pairTotal / uint64(bucketNumber))
	if average < 100 {
		average = 100
	}
	atomic.StoreUint64(&pr.upperThreshold, uint64(average*pr.loadFactor))
}