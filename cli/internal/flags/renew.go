package flags

type RenewFlags struct {
	*Global
}

func (r *RenewFlags) Check() error {
	err := r.Global.Check()
	if err != nil {
		return err
	}
	return nil
}
