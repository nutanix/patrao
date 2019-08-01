package upgradeagent

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
)

// NewContainer returns a new Container instance instantiated with the
// specified ContainerInfo and ImageInfo structs.
func NewContainer(containerInfo *types.ContainerJSON, imageInfo *types.ImageInspect) *Container {
	return &Container{
		containerInfo: containerInfo,
		imageInfo:     imageInfo,
	}
}

// Container represents a running Docker container.
type Container struct {
	Stale bool

	containerInfo *types.ContainerJSON
	imageInfo     *types.ImageInspect
}

// ID returns the Docker container ID.
func (c Container) ID() string {
	return c.containerInfo.ID
}

// Name returns the Docker container name.
func (c Container) Name() string {
	return c.containerInfo.Name
}

// ImageName returns the name of the Docker image that was used to start the
// container. If the original image was specified without a particular tag, the
// "latest" tag is assumed.
func (c Container) ImageName() string {
	imageName := c.containerInfo.Config.Image

	if !strings.Contains(imageName, ":") {
		imageName = fmt.Sprintf("%s:latest", imageName)
	}

	return imageName
}

// Labels returns labels information for dedicated container
func (c Container) Labels() map[string]string {
	return c.containerInfo.Config.Labels
}

// GetProjectName returns project name for given container
func (c Container) GetProjectName() (string, bool) {
	name, found := c.containerInfo.Config.Labels[DockerComposeProjectLabel]
	return name, found
}

// GetServiceName returns service name for given container
func (c Container) GetServiceName() (string, bool) {
	name, found := c.containerInfo.Config.Labels[DockerComposeServiceLabel]
	return name, found
}
