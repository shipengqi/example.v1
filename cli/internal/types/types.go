package types

const (
	NodeTypeWorker       = "worker"
	NodeTypeControlPlane = "controlplane"
)

const (
	SourceTypeConfigMap = "cm"
	SourceTypeSecret    = "secret"
)

const (
	CertTypeExternal   = "external"
	CertTypeInternal   = "internal"
	VaultPkiPathRE     = "RE"
	VaultPkiPathRIC    = "RIC"
	VaultPkiPathRID    = "RID"
	VaultPkiPathCUS    = "CUS"
	CertUnitTimeDay    = "d"
	CertUnitTimeMinute = "m"
)

const (
	DepWorker int = iota
	DepMaster
	DepMasterAndWorker
)
