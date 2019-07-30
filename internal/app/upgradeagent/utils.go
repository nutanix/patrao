package upgradeagent

import (
	"fmt"
	"strings"

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

func parseLabels(labels map[string]string) (*LocalSolutionInfo, error) {
	if value, found := labels[DockerComposeProjectLabel]; found {
		info := NewLocalSolutionInfo()
		info.SetDeploymentKind(DockerComposeDeployment)
		info.SetName(value)
		info.AddService(labels[DockerComposeServiceLabel])

		return info, nil
	}
	return nil, fmt.Errorf("can't parse Labels [%s]", labels)
}

// GetLocalSolutionList return the list of running solutions
func GetLocalSolutionList(containers *[]Container) *[]LocalSolutionInfo {
	var (
		list             []LocalSolutionInfo
		alreadyProcessed []string
	)
	for i, current := range *containers {
		info, err := parseLabels(current.Labels())
		if err != nil {
			log.Error(err)
			continue
		}
		name := info.GetName()
		if Contains(&alreadyProcessed, &name) {
			continue
		}
		for _, item := range (*containers)[i+1:] {
			tempInfo, e := parseLabels(item.Labels())
			if e != nil {
				log.Error(e)
				continue
			}
			if strings.Compare(info.GetName(), tempInfo.GetName()) == 0 {
				info.AddServices(tempInfo.GetServices())
			}
		}
		list = append(list, *info)
		alreadyProcessed = append(alreadyProcessed, info.GetName())
	}
	return &list
}

// GenUUID generate UUID string
func GenUUID() string {
	u4, err := uuid.NewV4()
	if err != nil {
		log.Error(err)
	}
	return u4.String()
}

// GenNodeUUID generate node uuid
func GenNodeUUID() string {
	// TBD
	return "node-uuid"
}
