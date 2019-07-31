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
