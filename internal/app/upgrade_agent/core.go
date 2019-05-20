package core

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

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

	if 0 == len(containers) {
		log.Info("There are no running containers at the moment")
		return rc
	}

	containerInfo, err := createRequest(&containers)

	if nil != err {
		log.Fatal(err)
	}

	request, _ := json.Marshal(containerInfo)
	log.Infof("Request to upstream service: %s", string(request))

	resp, rc := http.Post(context.GlobalString(UpstreamName), "json", bytes.NewBuffer(request))

	if nil != rc {
		log.Fatal(rc)
	}

	defer resp.Body.Close()

	body, rc := ioutil.ReadAll(resp.Body)

	if nil != rc {
		log.Fatal(rc)
	}

	log.Infof("Response from upstream service: %s", string(body))

	var newContainers []NewContainerVersion

	rc = json.Unmarshal([]byte(body), &newContainers)

	log.Info(newContainers)

	/*	for _, element := range containers {
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
