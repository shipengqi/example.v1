package options

type CheckOptions struct {
	CertType       string
	Cert           string
	Key            string
	CACert         string
	CDFNamespace   string
	Namespace      string
}

func (c *CheckOptions) Check() error {
	return nil
}
