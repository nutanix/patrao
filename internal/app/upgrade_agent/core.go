package core

import (
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
	log.Infoln("runOnce()")

	filter := BuildFilter(context.Args(), true)
	containers, rc := client.ListContainers(filter)

	if rc != nil {
		log.Fatal(rc)
	}

	for _, element := range containers {
		log.Infof("try to stop container %s", element.Name())

		if err := client.StopContainer(element, 0); err != nil {
			log.Fatal(err)
		}
	}

	for _, element := range containers {
		log.Infof("try to start container %s", element.Name())

		if err := client.StartContainer(element); nil != err {
			log.Fatal(err)
		}
	}

	return rc
}

func schedulePeriodicUpgrades(context *cli.Context) error {
	log.Infoln("schedulePeriodicUpgrades()")

	return nil
}
