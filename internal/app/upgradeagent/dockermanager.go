package upgradeagent

import (
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// DockerClient is interface which agent communicate with DockerAPI
type DockerClient interface {
	ListContainers() ([]Container, error)
	StopContainer(Container, time.Duration) error
	InspectContainer(*Container) (types.ContainerJSON, error)
	ExecContainer(*Container, string) (int, error)
	GetContainerByName(string, string) (*Container, error)
}

type dockerClient struct {
	api dockerclient.CommonAPIClient
}

// NewClient returns a new Client instance which can be used to interact with
// the Docker API.
// The client reads its configuration from the following environment variables:
//  * DOCKER_HOST			the docker-engine host to send api requests to
//  * DOCKER_TLS_VERIFY		whether to verify tls certificates
//  * DOCKER_API_VERSION	the minimum docker api version to work with
func NewClient() DockerClient {
	client, err := dockerclient.NewEnvClient()
	if err != nil {
		log.Fatalf("Error instantiating Docker client: %s", err)
	}
	return dockerClient{api: client}
}

func (client dockerClient) ListContainers() ([]Container, error) {
	cs := []Container{}
	bg := context.Background()
	runningContainers, err := client.api.ContainerList(
		bg,
		types.ContainerListOptions{})

	if err != nil {
		return nil, fmt.Errorf("DockerClient::ListContainers() [%v]", err)
	}
	for _, runningContainer := range runningContainers {
		containerInfo, err := client.api.ContainerInspect(bg, runningContainer.ID)
		if err != nil {
			return nil, fmt.Errorf("DockerClient::ListContainers() [%v]", err)
		}
		imageInfo, _, err := client.api.ImageInspectWithRaw(bg, containerInfo.Image)
		if err != nil {
			return nil, fmt.Errorf("DockerClient::ListContainers() [%v]", err)
		}
		c := Container{containerInfo: &containerInfo, imageInfo: &imageInfo}
		if PatraoAgentContainerName != c.Name() {
			cs = append(cs, c)
		}
	}
	return cs, nil
}

func (client dockerClient) StopContainer(c Container, timeout time.Duration) error {
	bg := context.Background()
	signal := DefaultStopSignal
	log.Infof("Stopping %s (%s) with %s", c.Name(), c.ID(), signal)
	if err := client.api.ContainerKill(bg, c.ID(), signal); err != nil {
		return fmt.Errorf("DockerClient::StopContainer() [%v]", err)
	}
	// Wait for container to exit, but proceed anyway after the timeout elapses
	client.waitForStop(c, timeout)
	if c.containerInfo.HostConfig.AutoRemove {
		log.Debugf("AutoRemove container %s, skipping ContainerRemove call.", c.ID())
	} else {
		log.Debugf("Removing container %s", c.ID())

		if err := client.api.ContainerRemove(bg, c.ID(),
			types.ContainerRemoveOptions{Force: true, RemoveVolumes: false}); err != nil {
			return fmt.Errorf("DockerClient::StopContainer() [%v]", err)
		}
	}
	// Wait for container to be removed. In this case an error is a good thing
	if err := client.waitForStop(c, timeout); err == nil {
		return fmt.Errorf("DockerClient::StopContainer() [Container %s (%s) could not be removed]", c.Name(), c.ID())
	}
	return nil
}

// waitForStop waits until container being stopped
func (client dockerClient) waitForStop(c Container, waitTime time.Duration) error {
	bg := context.Background()
	timeout := time.After(waitTime)
	for {
		select {
		case <-timeout:
			return nil
		default:
			if ci, err := client.api.ContainerInspect(bg, c.ID()); err != nil {
				return err
			} else if !ci.State.Running {
				return nil
			}
		}
		time.Sleep(1 * time.Second)
	}
}

//InspectContainer returns container configuration data structure
func (client dockerClient) InspectContainer(c *Container) (types.ContainerJSON, error) {
	return client.api.ContainerInspect(context.Background(), c.ID())
}

//waitForContainerExec waits while execution of the command is completed.
func (client dockerClient) waitForContainerExec(execID string, waitTime time.Duration) (int, error) {
	bg := context.Background()
	timeout := time.After(waitTime)
	for {
		select {
		case <-timeout:
			return DefaultExitCode, nil
		default:
			if ci, err := client.api.ContainerExecInspect(bg, execID); err != nil {
				return DefaultExitCode, err
			} else if !ci.Running {
				return ci.ExitCode, nil
			}
		}
		time.Sleep(1 * time.Second)
	}
}

// ExecContainer execute a command inside another container
func (client dockerClient) ExecContainer(c *Container, cmd string) (int, error) {
	ctx := context.Background()
	cmdWithParams := strings.Split(cmd, " ")
	config := types.ExecConfig{AttachStdin: false, AttachStdout: true, Cmd: cmdWithParams}
	execID, err := client.api.ContainerExecCreate(ctx, c.ID(), config)
	if err != nil {
		return DefaultExitCode, fmt.Errorf("DockerClient::ExecContainer() [%v]", err)
	}
	if _, err := client.api.ContainerExecAttach(ctx, execID.ID, types.ExecConfig{}); err != nil {
		return DefaultExitCode, fmt.Errorf("DockerClient::ExecContainer() [%v]", err)
	}
	err = client.api.ContainerExecStart(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return DefaultExitCode, fmt.Errorf("DockerClient::ExecContainer() [%v]", err)
	}
	return client.waitForContainerExec(execID.ID, DefaultTimeoutS*time.Second)
}

// GetContainerByName returns Container struct by solution name and container name
func (client dockerClient) GetContainerByName(solutionName string, containerName string) (*Container, error) {
	containers, err := client.ListContainers()
	if err != nil {
		return nil, fmt.Errorf("DockerClient::GetContainerByName() [%v]", err)
	}
	for _, item := range containers {
		if currSolutionName, found := item.GetProjectName(); found {
			if currServiceName, found := item.GetServiceName(); found && currSolutionName == solutionName && currServiceName == containerName {
				return &item, nil
			}
		}
	}
	return nil, fmt.Errorf("DockerClient::GetContainerByName() [%s]", "Container not found")
}
