package kbuild

import (
	"fmt"
	"github.com/spf13/cobra"
	k8s "github.com/xuxant/kbuild/pkg/kubernetes"
	"github.com/xuxant/kbuild/pkg/options"
	"os"
)

func newBuildCommand() *cobra.Command {
	buildCommand := &cobra.Command{
		Use:   "build",
		Short: "Build Container using kaniko",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			options.SetConfig(options.RGlobalOptions{
				KubeConfig: kubeConfig,
				Namespace:  namespace,
			})

			connection, msg := k8s.CheckPermissions()
			if !connection {
				for _, v := range msg {
					fmt.Printf("Failed testing kubernetes connection: %s\n", v)
				}
				os.Exit(2)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Subcommand ran")
		},
	}

	buildCommand.PersistentFlags().StringP("tag", "t", "", "image tag")
	buildCommand.PersistentFlags().StringP("dockerfile", "f", "Dockerfile", "path of the dockerfile")
	buildCommand.PersistentFlags().String("context", ".", "context path")

	return buildCommand
}
