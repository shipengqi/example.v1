package flags

type ApplyFlags struct {
	Remote         bool
}

func (a *ApplyFlags) Check() error {
	return nil
}
