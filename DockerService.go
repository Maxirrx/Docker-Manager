package main

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func NewDockerClient() (*client.Client, error) {

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return cli, nil
}

func RunDockerImage(path string) {
}

func StartDocker(uuid string) {
	cli, err := NewDockerClient()
	if err != nil {
		panic(err)
	}
	err = cli.ContainerStart(context.Background(), uuid, container.StartOptions{})
	if err != nil {
		panic(err)
	}
}

func RestartDocker(uuid string) {
	cli, err := NewDockerClient()
	if err != nil {
		panic(err)
	}
	err = cli.ContainerRestart(context.Background(), uuid, container.StopOptions{})
	if err != nil {
		panic(err)
	}
}

func StopDocker(uuid string) {
	cli, err := NewDockerClient()
	if err != nil {
		panic(err)
	}
	err = cli.ContainerStop(context.Background(), uuid, container.StopOptions{})
	if err != nil {
		panic(err)
	}
}

func DeleteDocker(uuid string) {
	StopDocker(uuid)
	cli, err := NewDockerClient()
	if err != nil {
		panic(err)
	}
	err = cli.ContainerRemove(context.Background(), uuid, container.RemoveOptions{})
	if err != nil {
		panic(err)
	}
}

func GetAllDocker() []types.Container {
	cli, err := NewDockerClient()
	if err != nil {
		panic(err)
	}
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		panic(err)
	}

	return containers
}

func FindDocker(uuid string) (string, types.Container) {
	containers := GetAllDocker()
	for _, container := range containers {
		if container.ID == uuid {
			return "", container
		}
	}
	return "y'a pas", types.Container{}
}
