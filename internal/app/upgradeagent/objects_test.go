package upgradeagent_test

import (
	"testing"

	core "github.com/nutanix/patrao/internal/app/upgradeagent"
	"github.com/stretchr/testify/assert"
)

func TestNewUpstreamResponseUpgradeInfo(t *testing.T) {
	v := core.NewUpstreamResponseUpgradeInfo()
	assert.NotNil(t, v)
	assert.Equal(t, core.Undefined, string(v.HealthCheckStatus))
	assert.Equal(t, "", string(v.Name))
	assert.Equal(t, "", string(v.Spec))
	assert.False(t, v.DeleteVolumes)
	assert.Equal(t, 0, v.ThresholdTimeS)
}

func TestNewNode(t *testing.T) {
	v := core.NewNode()
	assert.NotNil(t, v)
	assert.Equal(t, core.NodeKind, v.Kind)
	assert.NotEmpty(t, v.UUID)
	assert.NotEmpty(t, v.NodeUUID)
}

func TestNewAppTemplate(t *testing.T) {
	v := core.NewAppTemplate()
	assert.NotNil(t, v)
	assert.Equal(t, core.AppTemplateKind, string(v.Kind))
	assert.NotEmpty(t, v.UUID)
}

func TestNewDeployment(t *testing.T) {
	v := core.NewDeployment()
	assert.NotNil(t, v)
	assert.Equal(t, core.DeploymentKind, string(v.Kind))
	assert.NotEmpty(t, v.UUID)
}

func TestNewLocalSolutionInfo(t *testing.T) {
	info := core.NewLocalSolutionInfo()
	assert.NotNil(t, info)
	assert.Equal(t, info.GetDeploymentKind(), core.UndefinedDeployment)
}

func TestSetGetDeploymentKind(t *testing.T) {
	info := core.NewLocalSolutionInfo()
	deploymentKind := info.GetDeploymentKind()
	assert.Equal(t, core.UndefinedDeployment, deploymentKind)
	deploymentKind = core.DockerComposeDeployment
	info.SetDeploymentKind(deploymentKind)
	assert.Equal(t, deploymentKind, info.GetDeploymentKind())
}

func TestGetSetName(t *testing.T) {
	info := core.NewLocalSolutionInfo()
	assert.Empty(t, info.GetName())
	solutionName := "test"
	info.SetName(solutionName)
	assert.Equal(t, solutionName, info.GetName())
}

func TestGetServices(t *testing.T) {
	info := core.NewLocalSolutionInfo()
	assert.Empty(t, info.GetServices())
	info.AddService("test")
	assert.NotEmpty(t, info.GetServices())
}

func TestAddService(t *testing.T) {
	info := core.NewLocalSolutionInfo()
	info.AddService("test")
	assert.NotEmpty(t, info.GetServices())
	assert.Equal(t, []string{"test"}, info.GetServices())
	info.AddService("test1")
	assert.Equal(t, []string{"test", "test1"}, info.GetServices())
}

func TestAddServices(t *testing.T) {
	servicesArray := []string{"test1", "test2", "test3"}
	info := core.NewLocalSolutionInfo()
	assert.Empty(t, info.GetServices())
	info.AddServices(servicesArray)
	assert.Equal(t, servicesArray, info.GetServices())
}
