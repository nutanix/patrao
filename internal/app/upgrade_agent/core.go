package core

import (
	"github.com/jasonlvhit/gocron"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"time"
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

	if rc != nil {
		log.Fatal(rc)
	}

	for _, element := range containers {
		log.Infof("try to stop container [%s]", element.Name())

		if err := client.StopContainer(element, 0); err != nil {
			log.Fatal(err)
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

	/*
		for _, element := range containers {
			log.Infof("try to stop container [%s]", element.Name())

			if err := client.StopContainer(element, 0); err != nil {
				log.Fatal(err)
				continue
			}

			log.Infof("container is stopped [%s]", element.Name())
			log.Infof("try to start container [%s]", element.Name())

			if err := client.StartContainer(element); nil != err {
				log.Fatal(err)
				continue
			}
			log.Infof("container is started [%s]", element.Name())
		}
	*/
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
