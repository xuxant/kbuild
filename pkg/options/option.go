package options

type RGlobalOptions struct {
	Namespace   string
	KubeConfig  string
	KubeContext string
}

type RBuildOptions struct {
	Build      string
	Tag        string
	Context    string
	DockerFile string
	Command    string
}

var cfg RGlobalOptions

func SetConfig(config RGlobalOptions) {
	cfg = config
}

func GetConfig() RGlobalOptions {
	return cfg
}
