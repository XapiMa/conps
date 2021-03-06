package docker

import (
	"io"
	"os"
	"path"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/xapima/conps/pkg/util"
	"golang.org/x/net/context"
)

func (d *DockerApi) containerUp(repoTag string, containerName string) (string, error) {
	images, err := d.cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return "", util.ErrorWrapFunc(err)
	}
	log.Debug("Image list Got!")

	if !inImage(repoTag, images) {
		if err := d.imagePull(repoTag); err != nil {
			return "", util.ErrorWrapFunc(err)
		}
	}
	resp, err := d.cli.ContainerCreate(context.Background(), &container.Config{Image: repoTag, Cmd: []string{"/bin/sh"}, Tty: true}, nil, nil, containerName)
	if err != nil {
		return "", util.ErrorWrapFunc(err)
	}
	log.Debug("Container Created!")
	if err := d.cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", util.ErrorWrapFunc(err)
	}
	log.Debug("Container Started")
	return resp.ID, nil
}

func (d *DockerApi) containerDown(containerID string) error {
	log.Debug("in Container Down")
	statusCh, errCh := d.cli.ContainerWait(context.Background(), containerID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return util.ErrorWrapFunc(err)
		}
	case <-statusCh:
	}

	log.Debug("Container stoped!")

	if err := d.cli.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{RemoveVolumes: true, RemoveLinks: true, Force: true}); err != nil {
		return util.ErrorWrapFunc(err)
	}

	return nil
}

func (d *DockerApi) imagePull(repoTag string) error {
	reader, err := d.cli.ImagePull(context.Background(), path.Join("docker.io/library", repoTag), types.ImagePullOptions{})
	if err != nil {
		return util.ErrorWrapFunc(err)
	}
	io.Copy(os.Stdout, reader)
	log.Debug("image pulled!")

	return nil
}

func inImage(repoTag string, images []types.ImageSummary) bool {
	for _, v := range images {
		if inSlice(repoTag, v.RepoTags) {
			return true
		}
	}
	return false
}

func inSlice(target string, slice []string) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

func unknownCid(cid string) error {
	return errors.Wrap(UnknownCidError, cid)
}
