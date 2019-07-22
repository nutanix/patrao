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
