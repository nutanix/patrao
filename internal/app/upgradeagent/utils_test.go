package upgradeagent_test

import (
	"testing"

	core "github.com/nutanix/patrao/internal/app/upgradeagent"
	"github.com/stretchr/testify/assert"
)

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
	c := CreateTestContainer(t, containerInfo, imageInfo)
	info, err := core.ParseLabels(c.Labels())
	assert.NoError(t, err)
	assert.Equal(t, projectValue, info.GetName())
	assert.Equal(t, []string{"cache"}, info.GetServices())
	assert.Equal(t, core.DockerComposeDeployment, info.GetDeploymentKind())
	c1 := CreateTestContainer(t, containerInfoNoLabels, imageInfo)
	info, err = core.ParseLabels(c1.Labels())
	assert.Error(t, err)
	assert.Empty(t, info)
}
