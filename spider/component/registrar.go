package component

import (
	"fmt"
	"sync"
)

type Registrar interface {
	// Register 用于注册组件实例
	Register(c Component) (bool, error)
	// Unregister 用于注销组件实例
	UnRegister(cid CID) (bool, error)
	// Get 用于获取一个指定类型的组件的实例
	// 本函数应该基于负载均衡策略返回实例
	Get(cType Type) (Component, error)
	// GetAllByType 用于获取指定类型的所有组件实例
	GetAllByType(cType Type) (map[CID]Component, error)
	// GetAll 用于获取所有组件实例
	GetAll() map[CID]Component
	// Clear 会清除所有的组件注册记录
	Clear()
}

// NewRegistrar 用于创建一个组件注册器的实例。
func NewRegistrar() Registrar {
	return &RegistrarImpl{
		componentTypeMap: map[Type]map[CID]Component{},
	}
}

type RegistrarImpl struct {
	// componentTypeMap 代表组件类型与对应组件实例的映射
	componentTypeMap map[Type]map[CID]Component
	// rwlock 代表组件注册专用读写锁
	rwlock sync.RWMutex
}

func (r *RegistrarImpl) Register(c Component) (bool, error) {
	if c == nil {
		return false, errors.NewIllegalParameterError("nil component instance")
	}
	cid := c.ID()
	parts, err := SplitCID(cid)
	if err != nil {
		return false, err
	}
	cType := legalLetterTypeMap[parts[0]]
	if !CheckType(cType, c) {
		errMsg := fmt.Sprintf("incorrect module type: %s", cType)
		return false, errors.NewIllegalParameterError(errMsg)
	}
	r.rwlock.Lock()
	defer r.rwlock.Unlock()
	components := r.componentTypeMap[cType]
	if components == nil {
		components = map[CID]Component{}
	}
	if _, ok := components[cid]; ok {
		return false, nil
	}
	components[cid] = c
	r.componentTypeMap[cType] = components
	return true, nil
}


func (r *RegistrarImpl) UnRegister(cid CID) (bool, error) {
	parts, err := SplitCID(cid)
	if err != nil {
		return false, err
	}
	cType := legalLetterTypeMap[parts[0]]
	var deleted bool
	r.rwlock.Lock()
	defer r.rwlock.Unlock()
	if components, ok := r.componentTypeMap[cType]; ok {
		if _, ok := components[cid]; ok {
			delete(components, cid)
			deleted = true
		}
	}
	return deleted, nil
}

// Get 用于获取一个指定类型的组件的实例。
// 本函数会基于负载均衡策略返回实例。
func (r *RegistrarImpl) Get(cType Type) (Component, error) {
	components, err := r.GetAllByType(cType)
	if err != nil {
		return nil, err
	}
	minScore := uint64(0)
	var selectedComponent Component
	for _, c := range components {
		SetScore(c)
		score := c.Score()
		if minScore == 0 || score < minScore {
			selectedComponent = c
			minScore = score
		}
	}
	return selectedComponent, nil
}

// GetAllByType 用于获取指定类型的所有组件实例。
func (r *RegistrarImpl) GetAllByType(cType Type) (map[CID]Component, error) {
	if !LegalType(cType) {
		errMsg := fmt.Sprintf("illegal component type: %s", cType)
		return nil, errors.NewIllegalParameterError(errMsg)
	}
	r.rwlock.RLock()
	defer r.rwlock.RUnlock()
	modules := r.componentTypeMap[cType]
	if len(modules) == 0 {
		return nil, ErrNotFoundModuleInstance
	}
	result := map[CID]Component{}
	for cid, component := range components {
		result[cid] = component
	}
	return result, nil
}

// GetAll 用于获取所有组件实例。
func (r *RegistrarImpl) GetAll() map[CID]Component {
	result := map[CID]Component{}
	r.rwlock.RLock()
	defer r.rwlock.RUnlock()
	for _, components := range r.componentTypeMap {
		for cid, component := range components {
			result[cid] = component
		}
	}
	return result
}

// Clear 会清除所有的组件注册记录。
func (r *RegistrarImpl) Clear() {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()
	r.componentTypeMap = map[Type]map[CID]Component{}
}
