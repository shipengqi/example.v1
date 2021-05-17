package flags

type CreateFlags struct {
	*Global

	CAKey          string
	NodeType       string
	Host           string
	KubeApiCertSan string
}

func (c *CreateFlags) Check() error {
	err := c.Global.Check()
	if err != nil {
		return err
	}
	return nil
}
