package upgradeagent

import (
	"fmt"

	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
)

// Contains returns true in case if there is item in array. otherwise return false
func Contains(names *[]string, item *string) bool {
	for _, it := range *names {
		if it == *item {
			return true
		}
	}
	return false
}

// ParseLabels parse labels map to LocalSolutionInfo data structure
func ParseLabels(labels map[string]string) (*LocalSolutionInfo, error) {
	if value, found := labels[DockerComposeProjectLabel]; found {
		info := NewLocalSolutionInfo()
		info.SetDeploymentKind(DockerComposeDeployment)
		info.SetName(value)
		info.AddServices(labels[DockerComposeServiceLabel])

		return info, nil
	}
	return nil, fmt.Errorf("Cannot read labels [%s]", labels)
}

// GetLocalSolutionList return the list of running solutions
func GetLocalSolutionList(containers []Container) map[string]*LocalSolutionInfo {
	projectMap := make(map[string]*LocalSolutionInfo)
	if containers != nil {
		for _, current := range containers {
			info, err := ParseLabels(current.Labels())
			if err != nil {
				log.Error(err)
				continue
			}
			if _, ok := projectMap[info.GetName()]; ok {
				projectMap[info.GetName()].AddServices(info.GetServices()...)
			} else {
				projectMap[info.GetName()] = info
			}
		}
	}
	return projectMap
}

// GenUUID generate UUID string
func GenUUID() string {
	u4, _ := uuid.NewV4()
	return u4.String()
}

// GenNodeUUID generate node uuid
func GenNodeUUID() string {
	// TBD
	return "node-uuid"
}
