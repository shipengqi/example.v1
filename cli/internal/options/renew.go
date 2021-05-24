package options

type RenewOptions struct {
	*Global
}

func (r *RenewOptions) Check() error {
	err := r.Global.Check()
	if err != nil {
		return err
	}
	return nil
}
