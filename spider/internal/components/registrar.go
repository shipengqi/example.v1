package component

type Registrar interface {
	Register(c Component) (bool, error)
	UnRegister(cid CID) (bool, error)
	Get(cType Type) (Component, error)
	GetAllByType(cType Type) (map[CID]Component, error)
	GetAll() (map[CID]Component, error)
	Clear()
}
