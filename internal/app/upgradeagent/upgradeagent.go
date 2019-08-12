package upgradeagent

import (
	"github.com/jasonlvhit/gocron"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// Main - Upgrade Agent entry point
func Main(context *cli.Context) error {
	if context.GlobalBool(RunOnceName) {
		return runOnce(context)
	}
	return schedulePeriodicUpgrades(context)
}

// runOnce do check launched solutions and do upgrade them if there are new versions available
func runOnce(context *cli.Context) error {
	log.Infoln("[+]runOnce()")
	localSolutionList := GetLocalSolutionList(context, NewClient())
	for _, current := range localSolutionList {
		if current.CheckUpgrade() == true {
			err := current.DoUpgrade()
			if err != nil {
				log.Error(err)
				current.DoRollback()
				continue
			}
			if current.CheckHealth() == false {
				current.DoRollback()
			}
		}
	}
	log.Infoln("[-]runOnce()")
	return nil
}

// schedulePeriodicUpgrades schedules pereodic upgrade check using upgrade interval command line parameter
func schedulePeriodicUpgrades(context *cli.Context) error {
	log.Infoln("[+]schedulePeriodicUpgrades()")
	{
		gocron.Every(uint64(context.GlobalInt64(UpgradeIntervalName))).Seconds().Do(runOnce, context)
		<-gocron.Start()
	}
	log.Infoln("[-]schedulePeriodicUpgrades()")
	return nil
}
