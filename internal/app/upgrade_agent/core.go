package core

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// Main - Upgrade Agent entry point
func Main(context *cli.Context) error {

	if context.GlobalBool("run-once") {
		return runOnce(context)
	}

	return schedulePeriodicUpgrades(context)
}

func runOnce(context *cli.Context) error {
	log.Infoln("runOnce()")

	return nil
}

func schedulePeriodicUpgrades(context *cli.Context) error {
	log.Infoln("schedulePeriodicUpgrades()")

	return nil
}
