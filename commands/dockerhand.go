package commands

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	flDryrun   bool
	flEndpoint string

	cmdRoot = &cobra.Command{Use: "dockerhand", Short: "A command line tool Docker resources manager"}
)

func init() {
	cmdRoot.PersistentFlags().BoolVarP(&flDryrun, "dryrun", "d", false, "enable dryrun mode")
	cmdRoot.PersistentFlags().StringVarP(&flEndpoint, "endpoint", "e", "unix:///var/run/docker.sock", "docker endpoint uri")
	cobra.OnInitialize(initialize)
}

// Execute is the app entry point
func Execute() {
	cmdRoot.Execute()
}

func initialize() {
	if flEndpoint == "" {
		logrus.Fatal("no endpoint configured")
	}
}

func help(cmd *cobra.Command, args []string) {
	cmd.Help()
}
