package upgradeagent_test

import (
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	core "github.com/nutanix/patrao/internal/app/upgradeagent"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

type mockDockerClient struct {
	t          *testing.T
	containers []core.Container
}

func NewMockDockerClient(testing *testing.T, containerList []core.Container) core.DockerClient {
	return &mockDockerClient{
		t:          testing,
		containers: containerList,
	}
}

func (mock mockDockerClient) ListContainers() ([]core.Container, error) {
	return mock.containers, nil
}

func (mock mockDockerClient) StopContainer(core.Container, time.Duration) error {
	return nil
}

func (mock mockDockerClient) ExecContainer(*core.Container, string) (int, error) {
	return 0, nil
}

func (mock mockDockerClient) InspectContainer(*core.Container) (types.ContainerJSON, error) {
	var containerJSON types.ContainerJSON
	return containerJSON, nil
}

func (mock mockDockerClient) GetContainerByName(string, string) (*core.Container, error) {
	return nil, nil
}

func TestContains(t *testing.T) {
	array := []string{"one", "two", "three"}

	value := "two"
	assert.True(t, core.Contains(&array, &value))

	value = "four"
	assert.False(t, core.Contains(&array, &value))
}

func TestGenUUID(t *testing.T) {
	v1 := core.GenUUID()
	assert.NotEmpty(t, v1)
	v2 := core.GenUUID()
	assert.NotEmpty(t, v2)
	assert.NotEqual(t, v1, v2)
}

func TestGenNodeUUID(t *testing.T) {
	v1 := core.GenNodeUUID()
	assert.NotEmpty(t, v1)
	v2 := core.GenNodeUUID()
	assert.Equal(t, v1, v2)
}

func TestParseLabels(t *testing.T) {
	app := cli.NewApp()
	app.Flags = core.SetupAppFlags()
	app.Action = func(context *cli.Context) {
		c := CreateTestContainer(t, containerInfo, imageInfo)

		deploymentClient, err := core.ParseLabels(context, nil, c.Labels())
		assert.NoError(t, err)
		assert.Equal(t, projectValue, deploymentClient.GetLocalSolutionInfo().GetName())
		assert.Equal(t, []string{"cache"}, deploymentClient.GetLocalSolutionInfo().GetServices())
		assert.Equal(t, core.DockerComposeDeployment, deploymentClient.GetLocalSolutionInfo().GetDeploymentKind())
		c1 := CreateTestContainer(t, containerInfoNoLabels, imageInfo)
		deploymentClient, err = core.ParseLabels(context, nil, c1.Labels())
		assert.Error(t, err)
		assert.Empty(t, deploymentClient)
	}
	args := []string{"/Projects/Nutanix/patrao/cmd/upgradeagent/__debug_bin", "--run-once"}
	app.Run(args)
}

func TestGetLocalSolutionList(t *testing.T) {
	app := cli.NewApp()
	app.Flags = core.SetupAppFlags()
	app.Action = func(context *cli.Context) {
		assert.Empty(t, core.GetLocalSolutionList(context, NewMockDockerClient(t, []core.Container{})))
		c := CreateTestContainer(t, containerInfoNoLabels, imageInfo)
		assert.Empty(t, core.GetLocalSolutionList(context, NewMockDockerClient(t, []core.Container{*c})))
		c = CreateTestContainer(t, containerInfo, imageInfo)
		assert.NotEmpty(t, core.GetLocalSolutionList(context, NewMockDockerClient(t, []core.Container{*c})))
		c1 := CreateTestContainer(t, containerInfoNewName, imageInfo)
		assert.NotEmpty(t, core.GetLocalSolutionList(context, NewMockDockerClient(t, []core.Container{*c, *c1})))
	}
	args := []string{"/Projects/Nutanix/patrao/cmd/upgradeagent/__debug_bin", "--run-once"}
	app.Run(args)
}
