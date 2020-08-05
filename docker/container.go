package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/marcoshuck/container-wrapper/containers"
	"time"
)

type dockerContainer struct {
	id    string
	name  string
	image string
	ports nat.PortMap
	cli   *client.Client
	envs  map[string]string
}

// ID returns the dockerContainer's ID.
func (d dockerContainer) ID() string {
	return d.id
}

// Name returns the dockerContainer's name.
func (d dockerContainer) Name() string {
	return d.name
}

// Image returns the image's name that were used to create the dockerContainer.
func (d dockerContainer) Image() string {
	return d.image
}

// EnvVars returns the dockerContainer's environment variables.
func (d dockerContainer) EnvVars() containers.EnvVars {
	return d.envs
}

// Start starts the dockerContainer.
func (d dockerContainer) Start() error {
	return d.cli.ContainerStart(context.TODO(), d.id, types.ContainerStartOptions{})
}

// Stop stops the dockerContainer.
func (d dockerContainer) Stop() error {
	timeout := time.Second
	return d.cli.ContainerStop(context.TODO(), d.id, &timeout)
}

// Remove removes the dockerContainer.
func (d dockerContainer) Remove() error {
	return d.cli.ContainerRemove(context.TODO(), d.id, types.ContainerRemoveOptions{})
}

// NewContainer initializes a new containers.Container implementation using Docker.
func NewContainer(ID, name, image string, ports nat.PortMap, envVars containers.EnvVars, cli *client.Client) containers.Container {
	return &dockerContainer{
		id:    ID,
		name:  name,
		image: image,
		ports: ports,
		cli:   cli,
		envs:  envVars,
	}
}

