package upgradeagent_test

import (
	"encoding/json"
	"testing"
	"time"

	core "github.com/nutanix/patrao/internal/app/upgradeagent"
	"github.com/stretchr/testify/assert"
)

const upstreamResponseUpgradeInfo = "{\"Name\": \"test\", \"Spec\": \"version: \\\"3\\\"\\nservices:\\n  db:\\n    image: postgres:10.8\\n    expose:\\n      - 5432\\n    environment:\\n      - POSTGRES_USER=db_user_name\\n      - POSTGRES_PASSWORD=P56FJXc\",\"DeleteVolumes\": false,\"ThresholdTimeS\": 60, \"HealthCheckCmds\": [{\"ContainerName\": \"db\", \"Cmd\": \"pg_isready -U postgres\"}]}"

func CreateUpstreamResponseUpgradeInfo(t *testing.T) *core.UpstreamResponseUpgradeInfo {
	var info core.UpstreamResponseUpgradeInfo
	err := json.Unmarshal([]byte(upstreamResponseUpgradeInfo), &info)
	assert.NoError(t, err)

	return &info
}

func TestNewClient(t *testing.T) {
	client := core.NewClient(false)
	assert.NotNil(t, client)
	client = core.NewClient(true)
	assert.NotNil(t, client)
}

func TestGetContainerByName(t *testing.T) {
	client := core.NewClient(false)
	assert.NotNil(t, client)
	c, err := client.GetContainerByName("test", "db")
	assert.Nil(t, c)
	assert.EqualError(t, err, "Container not found")
}

func TestListContainers(t *testing.T) {
	client := core.NewClient(false)
	assert.NotNil(t, client)
	list, err := client.ListContainers()
	assert.Nil(t, err)
	assert.Empty(t, list)
}

func TestStopContainer(t *testing.T) {
	client := core.NewClient(false)
	assert.NotNil(t, client)

	c := CreateTestContainer(t)
	err := client.StopContainer(*c, core.DefaultTimeoutS*time.Second)
	assert.Error(t, err)
}

func TestInspectContainer(t *testing.T) {
	client := core.NewClient(false)
	assert.NotNil(t, client)

	c := CreateTestContainer(t)
	_, err := client.InspectContainer(c)
	assert.Error(t, err)
}

func TestExecContainer(t *testing.T) {
	client := core.NewClient(false)
	assert.NotNil(t, client)

	c := CreateTestContainer(t)
	_, err := client.ExecContainer(c, "/bin/bash")
	assert.Error(t, err)
}

func TestLaunchSolution(t *testing.T) {
	client := core.NewClient(false)
	assert.NotNil(t, client)

	info := CreateUpstreamResponseUpgradeInfo(t)
	assert.NotNil(t, info)

	err := client.LaunchSolution(info)
	assert.Error(t, err)
}
