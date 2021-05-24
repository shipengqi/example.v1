package options

type CreateOptions struct {
	*Global

	CAKey          string
	NodeType       string
	Host           string
	KubeApiCertSan string
}

func (c *CreateOptions) Check() error {
	err := c.Global.Check()
	if err != nil {
		return err
	}
	return nil
}
