package common

//KeySignOpts :
type KeySignOpts struct {
	PipelineDir  string
	RegistryPath string
	ImageTag     string
	PipelineOpt  string
	KeyRef       string
}

//KeyVerifyOpts :
type KeyVerifyOpts struct {
	PipelineDir  string
	RegistryPath string
	ImageTag     string
	PipelineOpt  string
	KeyRef       string
}
