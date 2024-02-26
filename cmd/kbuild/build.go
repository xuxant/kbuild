package kbuild

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xuxant/kbuild/pkg/build"
	k8s "github.com/xuxant/kbuild/pkg/kubernetes"
	"github.com/xuxant/kbuild/pkg/options"
	"github.com/xuxant/kbuild/pkg/output/log"
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
			tag, _ := cmd.Flags().GetString("tag")
			buildContext, _ := cmd.Flags().GetString("context")
			dockerfile, _ := cmd.Flags().GetString("dockerfile")
			err := build.InitializeBuild(tag, buildContext, dockerfile)
			if err != nil {
				log.Entry(context.TODO()).Infof("initialize kaniko build: %s", err)
			}
		},
	}

	buildCommand.PersistentFlags().StringP("tag", "t", "", "image tag")
	buildCommand.PersistentFlags().StringP("dockerfile", "f", "Dockerfile", "path of the dockerfile")
	buildCommand.PersistentFlags().String("context", ".", "context path")

	return buildCommand
}
