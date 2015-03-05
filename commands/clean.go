package commands

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/edgard/goutil"
	"github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
)

var (
	cmdClean           = &cobra.Command{Use: "clean", Short: "Clean resources", Run: help}
	cmdCleanAll        = &cobra.Command{Use: "all", Short: "Clean all resources", Run: cleanAll}
	cmdCleanContainers = &cobra.Command{Use: "containers", Short: "Clean inactive containers", Run: cleanContainers}
	cmdCleanImages     = &cobra.Command{Use: "images", Short: "Clean untagged images", Run: cleanImages}
	cmdCleanVolumes    = &cobra.Command{Use: "volumes", Short: "Clean orphaned volumes", Run: cleanVolumes}
)

func init() {
	cmdRoot.AddCommand(cmdClean)
	cmdClean.AddCommand(cmdCleanAll, cmdCleanContainers, cmdCleanImages, cmdCleanVolumes)
}

func cleanAll(cmd *cobra.Command, args []string) {
	cleanContainers(cmd, args)
	cleanImages(cmd, args)
	cleanVolumes(cmd, args)
}

func cleanContainers(cmd *cobra.Command, args []string) {
	client, err := docker.NewClient(flEndpoint)
	if err != nil {
		logrus.Error(err)
		return
	}
	list, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		logrus.Error(err)
		return
	}
	for _, container := range list {
		if strings.HasPrefix(container.Status, "Exited") {
			logrus.Infoln("Removing container:", container.ID[:12])
			if !flDryrun {
				err = client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID})
				if err != nil {
					logrus.Error(err)
					return
				}
			}
		}
	}
}

func cleanImages(cmd *cobra.Command, args []string) {
	client, err := docker.NewClient(flEndpoint)
	if err != nil {
		logrus.Error(err)
		return
	}
	list, err := client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		logrus.Error(err)
		return
	}
	for _, image := range list {
		if image.RepoTags[0] == "<none>:<none>" {
			logrus.Infoln("Removing image:", image.ID[:12])
			if !flDryrun {
				err = client.RemoveImage(image.ID)
				if err != nil {
					logrus.Error(err)
					return
				}
			}
		}
	}
}

func cleanVolumes(cmd *cobra.Command, args []string) {
	client, err := docker.NewClient(flEndpoint)
	if err != nil {
		logrus.Error(err)
		return
	}

	info, err := client.Info()
	if err != nil {
		logrus.Error(err)
		return
	}
	rootDir := info.Get("DockerRootDir")
	volumesDir := path.Join(rootDir, "volumes")
	vfsDir := path.Join(rootDir, "vfs", "dir")

	list, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		logrus.Error(err)
		return
	}

	var active []string
	for _, container := range list {
		inspect, err := client.InspectContainer(container.ID)
		if err != nil {
			logrus.Error(err)
			return
		}
		for _, volume := range inspect.Volumes {
			if strings.HasPrefix(volume, volumesDir) || strings.HasPrefix(volume, vfsDir) {
				_, dir := path.Split(volume)
				active = append(active, dir)
			}
		}
	}

	var existing []string
	fl, err := ioutil.ReadDir(volumesDir)
	if err != nil {
		logrus.Error(err)
		return
	}
	for _, dir := range fl {
		existing = append(existing, dir.Name())
	}

	diff := goutil.DiffInStringSlice(active, existing)
	for _, dir := range diff {
		logrus.Infoln("Removing volume:", dir)
		if !flDryrun {
			err = os.RemoveAll(path.Join(vfsDir, dir))
			if err != nil {
				logrus.Warn(err)
			}
			err = os.RemoveAll(path.Join(volumesDir, dir))
			if err != nil {
				logrus.Warn(err)
			}
		}
	}
}
