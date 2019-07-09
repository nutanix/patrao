package upgradeagent_test

import (
	"testing"

	core "github.com/nutanix/patrao/internal/app/upgradeagent"
	"github.com/stretchr/testify/assert"
)

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
