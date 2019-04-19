package Core

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// Main - Upgrade Agent entry point
func Main(context *cli.Context) error {

	if context.GlobalBool("run-once") {
		return runOnce(context)
	}

	return shchedulePereodicUpgrades(context)
}

func runOnce(context *cli.Context) error {
	log.Infoln("runOnce()")

	return nil
}

func shchedulePereodicUpgrades(context *cli.Context) error {
	log.Infoln("shchedulePereodicUpgrades()")

	return nil
}
