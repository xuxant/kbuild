package main

import "github.com/xuxant/kbuild/cmd/kbuild"

func main() {
	cmd := kbuild.NewBuildCommand()
	cmd.Execute()
}
