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
)

var (
	client   Client
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

	log.Infoln("[-]runOnce()")

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
