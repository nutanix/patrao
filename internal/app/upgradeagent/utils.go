package upgradeagent

import (
	"fmt"
	"strings"
	"time"

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

// GetSolutionAndServiceName returns solition name, container name by original container name
func GetSolutionAndServiceName(containerName string) (string, string, error) {
	var (
		err error
		rc  string
	)
	nameParts := strings.Split(containerName, "_")
	length := len(nameParts)

	if length == 0 || length < 2 {
		err = SolutionNameNotFound{
			time.Date(1989, 3, 15, 22, 30, 0, 0, time.UTC),
			fmt.Sprintf("getSolutionAndServiceName(): can't identify solution name [%s]", containerName),
		}
		return rc, rc, err
	}

	return nameParts[0][1:], nameParts[1], err
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
