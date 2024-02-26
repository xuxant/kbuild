package kaniko

func Args(tag, context string) []string {
	args := []string{
		"--destination", tag,
		"--dockerfile", "Dockerfile",
	}

	if context != "" {
		args = append(args, "--context", context)
	}

	return args
}
