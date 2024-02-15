package main

import (
	"github.com/xuxant/kbuild/cmd/kbuild"
	"log"
)

func main() {
	cmd := kbuild.NewBuildCommand()
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
