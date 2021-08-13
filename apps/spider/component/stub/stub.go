package stub

import (
	"fmt"
	component2 "github.com/shipengqi/example.v1/apps/spider/component"
	"sync/atomic"
)

// ComponentInternal 代表组件的内部基础接口类型。
type ComponentInternal interface {
	component2.Component
	// IncrCalledCount 会把调用计数增 1
	IncrCalledCount()
	// IncrAcceptedCount 会把接受计数增 1
	IncrAcceptedCount()
	// IncrCompletedCount 会把成功完成计数增 1
	IncrCompletedCount()
	// IncrHandlingNumber 会把实时处理数增 1
	IncrHandlingNumber()
	// DecrHandlingNumber 会把实时处理数减 1
	DecrHandlingNumber()
	// Clear 用于清空所有计数。
	Clear()
}

type ComponentInternalImpl struct {
	// cid 代表组件 ID。
	cid component2.CID
	// addr 代表组件的网络地址
	addr string
	// score 代表组件评分
	score uint64
	// scoreCalculator 代表评分计算器
	scoreCalculator component2.CalculateScore
	// calledCount 代表调用计数
	calledCount uint64
	// acceptedCount 代表接受计数
	acceptedCount uint64
	// completedCount 代表成功完成计数
	completedCount uint64
	// handlingNumber 代表实时处理数
	handlingNumber uint64
}

func NewComponentInternal(
	cid component2.CID,
	scoreCalculator component2.CalculateScore) (ComponentInternal, error) {
	parts, err := component2.SplitCID(cid)
	if err != nil {
		return nil, errors.NewIllegalParameterError(
			fmt.Sprintf("illegal ID %q: %s", mid, err))
	}
	return &ComponentInternalImpl{
		cid:             cid,
		addr:            parts[2],
		scoreCalculator: scoreCalculator,
	}, nil
}

func (c *ComponentInternalImpl) ID() component2.CID {
	return c.cid
}

func (c *ComponentInternalImpl) Addr() string {
	return c.addr
}


func (c *ComponentInternalImpl) Score() uint64 {
	return atomic.LoadUint64(&c.score)
}

func (c *ComponentInternalImpl) SetScore(score uint64) {
	atomic.StoreUint64(&c.score, score)
}
