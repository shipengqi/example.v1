package action

type check struct {
	*action
}

func NewCheck(cfg *Configuration) Interface {
	return &check{&action{
		name: "check",
		cfg:  cfg,
	}}
}

func (a *check) Name() string {
	return a.name
}

func (a *check) Run() error {
	return nil
}

func (a *check) PostRun() error {
	return nil
}

func (a *check) Execute() error {
	err := a.PreRun()
	if err != nil {
		return err
	}
	err = a.Run()
	if err != nil {
		return err
	}
	return a.PostRun()
}
