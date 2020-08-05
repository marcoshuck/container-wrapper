package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/marcoshuck/container-wrapper/containers"
	"io"
	"os"
)

// Docker represents a docker client.
type Docker interface {
	Pull(ctx context.Context, image string) error
	Create(ctx context.Context, input CreateContainerInput) (containers.Container, error)
	Remove(ctx context.Context, id string) error
}

// CreateContainerInput is used to group a set of required fields when creating a dockerContainer.
type CreateContainerInput struct {
	Name         string
	Image        string
	Ports 		 nat.PortMap
	EnvVars      containers.EnvVars
}

// docker is a Docker implementation.
type docker struct {
	CLI *client.Client
}

// Pull pulls the given image from docker images repository.
func (d docker) Pull(ctx context.Context, image string) error {
	r, err := d.CLI.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, r)
	defer r.Close()
	return nil
}

// Create creates a new docker dockerContainer.
func (d docker) Create(ctx context.Context, input CreateContainerInput) (containers.Container, error) {
	body, err := d.CLI.ContainerCreate(
		ctx,
		&container.Config{
			Image: input.Image,
			Env:   input.EnvVars.ToSlice(),
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{},
		},
		nil,
		nil,
		input.Name,
	)
	if err != nil {
		return nil, err
	}

	return NewContainer(body.ID, input.Name, input.Image, input.Ports, input.EnvVars, d.CLI), nil
}

// Remove removes the dockerContainer that matches the given id.
func (d docker) Remove(ctx context.Context, id string) error {
	return d.CLI.ContainerRemove(ctx, id, types.ContainerRemoveOptions{
		RemoveVolumes: false,
		RemoveLinks:   false,
		Force:         true,
	})
}

// NewClient initializes a new Docker client.
func NewClient() (Docker, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return &docker{
		CLI: cli,
	}, nil
}
