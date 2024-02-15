package kbuild

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xuxant/kbuild/pkg/version"
)

func getVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "print app version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.Get().Version)
		},
	}
	return versionCmd
}
