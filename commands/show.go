package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
)

var (
	cmdShow             = &cobra.Command{Use: "show", Short: "Show resources", Run: help}
	cmdShowContainers   = &cobra.Command{Use: "containers", Short: "Show containers resources", Run: help}
	cmdShowContainersIP = &cobra.Command{Use: "ip", Short: "Show containers ip list", Run: showContainersIP}
)

func init() {
	cmdRoot.AddCommand(cmdShow)
	cmdShow.AddCommand(cmdShowContainers)
	cmdShowContainers.AddCommand(cmdShowContainersIP)
}

func showContainersIP(cmd *cobra.Command, args []string) {
	client, err := docker.NewClient(flEndpoint)
	if err != nil {
		logrus.Error(err)
		return
	}
	list, err := client.ListContainers(docker.ListContainersOptions{All: false})
	if err != nil {
		logrus.Error(err)
		return
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "CONTAINER ID\tIP\tNAMES")
	for _, container := range list {
		insp, err := client.InspectContainer(container.ID)
		if err != nil {
			logrus.Error(err)
			return
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", container.ID[:12], insp.NetworkSettings.IPAddress, insp.Name[1:])
	}
	w.Flush()
}
