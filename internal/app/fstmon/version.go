package fstmon

type InitFlags struct {
	CommitHash string
	Version    string
	Repository string
}

func (inf InitFlags) GetCommitHash() string {
	return inf.CommitHash
}

func (inf InitFlags) GetVersion() string {
	return inf.Version
}

func (inf InitFlags) GetRepository() string {
	return inf.Version
}
