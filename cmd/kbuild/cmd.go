package kbuild

import (
	"github.com/spf13/cobra"
	"os"
	"os/user"
)

var kubeConfig, namespace string

func NewBuildCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "kbuild",
		Short: "Remote build with Kaniko",

		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.PersistentFlags().StringVarP(&kubeConfig, "kubeconfig", "k", getDefaultKubeConfig(), "path for kubernetes config file")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "namespace for building container")

	rootCmd.AddCommand(getVersionCmd())
	rootCmd.AddCommand(newBuildCommand())

	return rootCmd
}

func getDefaultKubeConfig() string {
	currentUser, err := user.Current()
	if err != nil {
		os.Exit(5)
	}
	homeDir := currentUser.HomeDir
	defaultKubeConfig := homeDir + "/.kube/config"
	return defaultKubeConfig
}
