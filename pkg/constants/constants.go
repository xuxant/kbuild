package constants

type Phase string

const (
	Init  = Phase("Init")
	Build = Phase("Build")

	DefaultDockerFilePath = "Dockerfile"
)
