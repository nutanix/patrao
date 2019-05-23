package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"strings"

	"github.com/jasonlvhit/gocron"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	client Client
)

// Main - Upgrade Agent entry point
func Main(context *cli.Context) error {

	client = NewClient(false)

	if context.GlobalBool(RunOnceName) {
		return runOnce(context)
	}

	return schedulePeriodicUpgrades(context)
}

func runOnce(context *cli.Context) error {
	log.Infoln("[+]runOnce()")

	filter := BuildFilter(make([]string, 0), false)
	containers, rc := client.ListContainers(filter)

	if nil != rc {
		log.Fatal(rc)
	}

	getURLPath := context.GlobalString(UpstreamName) + UpstreamGetUpgrade

	log.Infof("URL path: %s", getURLPath)

	var solutions []string
	var upgradeInfoArray []UpstreamResponseUpgradeInfo

	for _, current := range containers {
		currentSolution, err := getSolutionName(current.Name())

		if nil != err {
			log.Error(err)
			continue
		}

		if contains(&solutions, &currentSolution) {
			continue
		}
		solutions = append(solutions, currentSolution)

		currentPath := getURLPath + currentSolution
		log.Infof("GET to Upstream: %s", currentPath)

		resp, rc := http.Get(currentPath)

		if nil != rc {
			log.Error(rc)
			continue
		}

		defer resp.Body.Close()

		body, rc := ioutil.ReadAll(resp.Body)

		if nil != rc {
			log.Error(rc)
			continue
		}

		log.Infof("Response from upstream service: %s", string(body))

		var toUpgrade UpstreamResponseUpgradeInfo
		rc = json.Unmarshal([]byte(body), &toUpgrade)

		if nil != rc {
			log.Error(rc)
			continue
		}

		upgradeInfoArray = append(upgradeInfoArray, toUpgrade)
	}

	for _, item := range upgradeInfoArray {
		for _, container := range containers {
			name, _ := getSolutionName(container.Name())
			if item.Name == name {
				err := client.StopContainer(container, 0)
				if nil != err {
					// TBD
				}
			}
		}

	}

	/*	var newContainers []NewContainerVersion

		rc = json.Unmarshal([]byte(body), &newContainers)

		log.Info(newContainers)

			for _, element := range containers {
				log.Infof("try to stop container [%s]", element.Name())

				if err := client.StopContainer(element, 0); err != nil {
					log.Infoln(err)
					continue
				}

				log.Infof("container is stopped [%s]", element.Name())
			}

			//
			log.Info("wait for 10 seconds.................")
			time.Sleep(10 * time.Second)
			//

			for _, element := range containers {
				log.Infof("try to start container [%s]", element.Name())

				if err := client.StartContainer(element); nil != err {
					log.Fatal(err)
					continue
				}
				log.Infof("container is started [%s]", element.Name())
			}

			log.Infoln("[-]runOnce()")
	*/

	return rc
}

func schedulePeriodicUpgrades(context *cli.Context) error {
	log.Infoln("[+]schedulePeriodicUpgrades()")

	gocron.Every(30).Seconds().Do(runOnce, context)

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

func getSolutionName(containerName string) (string, error) {
	var (
		err error
		rc  string
	)

	slice := strings.Split(containerName, "_")

	if 0 == len(slice) {
		err = SolutionNameNotFound{
			time.Date(1989, 3, 15, 22, 30, 0, 0, time.UTC),
			fmt.Sprintf("getSolutionName(): can't identify solution name [%s]", containerName),
		}
		return rc, err
	}

	return slice[0][1:], err
}

/*
func createRequest(containers *[]Container) ([]CurrentContainerVersion, error) {
	var (
		err error
		rc  []CurrentContainerVersion
	)

	for _, item := range *containers {
		current := new(CurrentContainerVersion)

		current.ID = "container-id"    // TBD will be changed to ID according specifications
		current.CREATED = "DD.MM.YYYY" // TBD will be changed to real created data according specifications
		current.NAME = item.Name()[1:]
		current.IMAGE = item.ImageName()

		rc = append(rc, *current)

	}

	return rc, err
}
*/
