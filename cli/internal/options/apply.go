package options

type ApplyOptions struct {
	Remote         bool
}

func (a *ApplyOptions) Check() error {
	return nil
}
