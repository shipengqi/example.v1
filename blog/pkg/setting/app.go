package setting

type app struct {
	singingKey   string
	pageSize     int
	isPrintStack bool
}

func (a *app) IsPrintStack() bool {
	return a.isPrintStack
}

func (a *app) SetIsPrintStack(isPrintStack bool) {
	a.isPrintStack = isPrintStack
}

func (a *app) PageSize() int {
	return a.pageSize
}

func (a *app) SetPageSize(pageSize int) {
	a.pageSize = pageSize
}

func (a *app) SingingKey() string {
	return a.singingKey
}

func (a *app) SetSingingKey(singingKey string) {
	a.singingKey = singingKey
}
