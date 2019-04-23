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
	// for debug. don't forget to remove!
	log.Info(containers)
	//

	return rc
}

func schedulePeriodicUpgrades(context *cli.Context) error {
	log.Infoln("schedulePeriodicUpgrades()")

	return nil
}
