package flags

type CheckFlags struct {
	CertType       string
	Cert           string
	Key            string
	CACert         string
	CDFNamespace   string
	Namespace      string
}

func (c *CheckFlags) Check() error {
	return nil
}
