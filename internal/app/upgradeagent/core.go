package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
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
	upgradeInfoArray, rc := getLaunchedSolutionsList(&containers, context)
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
	gocron.Every(uint64(context.GlobalInt64(UpgradeIntervalName))).Seconds().Do(runOnce, context)
	<-gocron.Start()
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

func getLaunchedSolutionsList(containers *[]Container, context *cli.Context) ([]UpstreamResponseUpgradeInfo, error) {
	var (
		solutions []string
		rc        []UpstreamResponseUpgradeInfo
		err       error
	)

	log.Info("[+]getLaunchedSolutionsList")
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
	log.Info("[-]getLaunchedSolutionsList")
	return rc, err
}

func doUpgradeSolutions(upgradeInfoList *[]UpstreamResponseUpgradeInfo, containers *[]Container) error {
	var rc error

	log.Info("[+]doUpgradeSolutions")
	for _, item := range *upgradeInfoList {
		if !isNewVersion(&item, containers) {
			log.Infof("Solution [%s] is up-to-date.", item.Name)
			continue
		}
		for _, container := range *containers {
			name, _, _ := getSolutionAndServiceName(container.Name())
			if item.Name == name {
				err := client.StopContainer(container, 0)
				if nil != err {
					log.Error(err)
					// TBD
					continue
				}
			}
		}
		if _, isFileExist := os.Stat(item.Name); !os.IsNotExist(isFileExist) {
			os.RemoveAll(item.Name)
		}
		err := os.Mkdir(item.Name, os.ModePerm)
		if nil != err {
			log.Error(err)
			// TBD
			continue
		}
		defer os.Remove(item.Name)
		dockerComposeFileName := fmt.Sprintf("%s/%s", item.Name, DockerComposeFileName)

		f, err := os.Create(dockerComposeFileName)
		if err != nil {
			log.Error(err)
			continue
			//TBD
		}
		defer func() {
			f.Close()
			os.Remove(dockerComposeFileName)
		}()

		_, err = f.Write([]byte(item.Spec)[:len(item.Spec)])
		if nil != err {
			log.Error(err)
			// TBD
			continue
		}

		log.Infof("Launching solution [%s]", item.Name)

		cmd := exec.Command(DockerComposeCommand, "-f", fmt.Sprintf("%s/%s", rootPath, dockerComposeFileName), "up", "-d")
		err = cmd.Run()

		if nil != err {
			log.Error(err)
			// TBD
			continue
		}
		log.Infof("Solution [%s] is successful launched", item.Name)
	}
	log.Info("[-]doUpgradeSolutions")
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
