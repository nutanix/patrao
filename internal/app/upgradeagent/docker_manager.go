package upgradeagent

import (
	"fmt"
	"os"
	"os/exec"
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
	LaunchSolution(*UpstreamResponseUpgradeInfo) error
}

type dockerClient struct {
	api        dockerclient.CommonAPIClient
	pullImages bool
}

// NewClient returns a new Client instance which can be used to interact with
// the Docker API.
// The client reads its configuration from the following environment variables:
//  * DOCKER_HOST			the docker-engine host to send api requests to
//  * DOCKER_TLS_VERIFY		whether to verify tls certificates
//  * DOCKER_API_VERSION	the minimum docker api version to work with
func NewClient(pullImages bool) DockerClient {
	client, err := dockerclient.NewEnvClient()
	if err != nil {
		log.Fatalf("Error instantiating Docker client: %s", err)
	}
	return dockerClient{api: client, pullImages: pullImages}
}

func (client dockerClient) ListContainers() ([]Container, error) {
	cs := []Container{}
	bg := context.Background()
	runningContainers, err := client.api.ContainerList(
		bg,
		types.ContainerListOptions{})

	if err != nil {
		return nil, err
	}
	for _, runningContainer := range runningContainers {
		containerInfo, err := client.api.ContainerInspect(bg, runningContainer.ID)
		if err != nil {
			return nil, err
		}
		imageInfo, _, err := client.api.ImageInspectWithRaw(bg, containerInfo.Image)
		if err != nil {
			return nil, err
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
		return err
	}
	// Wait for container to exit, but proceed anyway after the timeout elapses
	client.waitForStop(c, timeout)
	if c.containerInfo.HostConfig.AutoRemove {
		log.Debugf("AutoRemove container %s, skipping ContainerRemove call.", c.ID())
	} else {
		log.Debugf("Removing container %s", c.ID())

		if err := client.api.ContainerRemove(bg, c.ID(),
			types.ContainerRemoveOptions{Force: true, RemoveVolumes: false}); err != nil {
			return err
		}
	}
	// Wait for container to be removed. In this case an error is a good thing
	if err := client.waitForStop(c, timeout); err == nil {
		return fmt.Errorf("Container %s (%s) could not be removed", c.Name(), c.ID())
	}
	return nil
}

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

func (client dockerClient) LaunchSolution(info *UpstreamResponseUpgradeInfo) error {
	if _, isFileExist := os.Stat(info.Name); !os.IsNotExist(isFileExist) {
		os.RemoveAll(info.Name)
	}
	err := os.Mkdir(info.Name, os.ModePerm)
	if nil != err {
		log.Error(err)
		return err
	}
	defer os.Remove(info.Name)
	dockerComposeFileName := fmt.Sprintf("%s/%s", info.Name, DockerComposeFileName)
	f, err := os.Create(dockerComposeFileName)
	if err != nil {
		log.Error(err)
		return err
	}

	defer func() {
		f.Close()
		os.Remove(dockerComposeFileName)
	}()

	_, err = f.Write([]byte(info.Spec)[:len(info.Spec)])
	if err != nil {
		log.Error(err)
		return err
	}
	log.Infof("Launching solution [%s]", info.Name)
	cmd := exec.Command(DockerComposeCommand, "-f", fmt.Sprintf("%s/%s", rootPath, dockerComposeFileName), "up", "-d")
	err = cmd.Run()
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
