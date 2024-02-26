package constants

type Phase string

const (
	Init  = Phase("Init")
	Build = Phase("Build")

	DefaultDockerFilePath = "Dockerfile"
	RandomPodCharacter    = int64(6)
	DefaultKanikoImage    = "gcr.io/kaniko-project/executor"
)
