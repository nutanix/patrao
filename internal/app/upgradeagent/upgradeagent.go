package upgradeagent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"strings"

	"github.com/jasonlvhit/gocron"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

var (
	client   DockerClient
	rootPath string
)

// Main - Upgrade Agent entry point
func Main(context *cli.Context) error {
	ex, _ := os.Executable()
	rootPath = filepath.Dir(ex)
	client = NewClient(false)

	if context.GlobalBool(RunOnceName) {
		return runOnce(context)
	}
	return schedulePeriodicUpgrades(context)
}

func runOnce(context *cli.Context) error {
	log.Infoln("[+]runOnce()")
	containers, rc := client.ListContainers()
	if nil != rc {
		log.Error(rc)
		return rc
	}
	if len(containers) == 0 {
		log.Info("There are no launched containers on the host")
		return nil
	}
	upgradeInfoArray, rc := getLaunchedSolutionsList(context, &containers)
	if nil != rc {
		log.Error(rc)
		return rc
	}
	rc = doUpgradeSolutions(&upgradeInfoArray, &containers)
	log.Infoln("[-]runOnce()")
	return rc
}

func schedulePeriodicUpgrades(context *cli.Context) error {
	log.Infoln("[+]schedulePeriodicUpgrades()")
	{
		gocron.Every(uint64(context.GlobalInt64(UpgradeIntervalName))).Seconds().Do(runOnce, context)
		<-gocron.Start()
	}
	log.Infoln("[-]schedulePeriodicUpgrades()")
	return nil
}

func contains(names *[]string, item *string) bool {
	for _, it := range *names {
		if it == *item {
			return true
		}
	}
	return false
}

func getSolutionAndServiceName(containerName string) (string, string, error) {
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

func getLaunchedSolutionsList(context *cli.Context, containers *[]Container) ([]UpstreamResponseUpgradeInfo, error) {
	var (
		solutions []string
		rc        []UpstreamResponseUpgradeInfo
		err       error
	)
	getURLPath := context.GlobalString(UpstreamName) + UpstreamGetUpgrade

	for _, current := range *containers {
		currentSolution, _, err := getSolutionAndServiceName(current.Name())
		if nil != err {
			log.Error(err)
			continue
		}
		if contains(&solutions, &currentSolution) {
			continue
		}
		solutions = append(solutions, currentSolution)
		currentPath := getURLPath + currentSolution
		resp, err := http.Get(currentPath)
		if nil != err {
			log.Error(err)
			continue
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if nil != err {
			log.Error(err)
			continue
		}

		var toUpgrade UpstreamResponseUpgradeInfo
		err = json.Unmarshal([]byte(body), &toUpgrade)

		if nil != err {
			log.Error(err)
			continue
		}

		rc = append(rc, toUpgrade)
	}
	return rc, err
}

func doUpgradeSolutions(upgradeInfoList *[]UpstreamResponseUpgradeInfo, containers *[]Container) error {
	var rc error
	for _, item := range *upgradeInfoList {
		if !isNewVersion(&item, containers) {
			log.Infof("Solution [%s] is up-to-date.", item.Name)
			continue
		}
		for _, container := range *containers {
			name, _, _ := getSolutionAndServiceName(container.Name())
			if item.Name == name {
				err := client.StopContainer(container, StopContainerDefaultTimeoutS)
				if nil != err {
					log.Error(err)
				}
			}
		}
		err := client.LaunchSolution(&item)
		if err != nil {
			log.Error(err)
			continue
		}
		log.Infof("Solution [%s] is successful launched", item.Name)
	}
	return rc
}

// isNewVersion check if there are new version available.
func isNewVersion(upgradeInfo *UpstreamResponseUpgradeInfo, containers *[]Container) bool {
	var rc bool

	rc = false
	containersSpec := make(map[string]ContainerSpec)
	specMap := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(upgradeInfo.Spec), specMap)
	if nil != err {
		log.Error(err)
		return false
	}
	if val, exists := specMap[DockerComposeServicesName]; exists {
		servicesMap := val.(map[interface{}]interface{})
		for service := range servicesMap {
			var (
				serviceImageName string
				found            bool
				ok               bool
			)
			details := servicesMap[service].(map[interface{}]interface{})
			found = false
			for item := range details {
				if DockerComposeImageName == item {
					val, itemFound := details[item]
					if itemFound {
						serviceImageName, ok = val.(string)
						if ok {
							found = true
							break
						}
					}
				}
			}
			if found {

				containersSpec[service.(string)] = ContainerSpec{Name: service.(string), Image: serviceImageName}
			}
		}
		for _, container := range *containers {
			solutionName, serviceName, _ := getSolutionAndServiceName(container.Name())
			val, exist := containersSpec[serviceName]
			if (exist) && (upgradeInfo.Name == solutionName) {
				if container.ImageName() != val.Image {
					rc = true
					break
				}
			}
		}
	}
	return rc
}
