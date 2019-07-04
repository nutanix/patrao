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

func TestGetSolutionAndServiceName(t *testing.T) {
	solution_name, service_name, err := core.GetSolutionAndServiceName("/solution_service_1")
	assert.Equal(t, "solution", solution_name)
	assert.Equal(t, "service", service_name)
	assert.NoError(t, err)

	solution_name, service_name, err = core.GetSolutionAndServiceName("badsolutionname")
	assert.Equal(t, "", solution_name)
	assert.Equal(t, "", service_name)
	assert.Error(t, err)

	solution_name, service_name, err = core.GetSolutionAndServiceName("")
	assert.Equal(t, "", solution_name)
	assert.Equal(t, "", service_name)
	assert.Error(t, err)

	solution_name, service_name, err = core.GetSolutionAndServiceName("/solution_1")
	assert.Equal(t, "solution", solution_name)
	assert.Equal(t, "1", service_name)
	assert.NoError(t, err)
}
