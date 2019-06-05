package main

import (
	"os"

	core "github.com/nutanix/patrao/internal/app/upgradeagent"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func start(context *cli.Context) error {
	return core.Main(context)
}

func setupAppFlags() []cli.Flag {

	return []cli.Flag{
		cli.StringFlag{
			Name:   core.HostName,
			Usage:  core.HostUsage,
			Value:  core.HostValue,
			EnvVar: core.HostEnvVar,
		},
		cli.StringFlag{
			Name:   core.UpstreamName,
			Usage:  core.UpstreamUsage,
			Value:  core.UpstreamValue,
			EnvVar: core.UpstreamEnvVar,
		},
		cli.StringFlag{
			Name:   core.UpgradeIntervalName,
			Usage:  core.UpgradeIntervalUsage,
			Value:  core.UpgradeIntervalValue,
			EnvVar: core.UpgradeIntervalValueEnvVar,
		},
		cli.BoolFlag{
			Name:  core.RunOnceName,
			Usage: core.RunOnceUsage,
		},
	}
}

func createApp() *cli.App {
	app := cli.NewApp()
	app.Name = core.ApplicationName
	app.Usage = core.ApplicationUsage
	app.Action = start
	app.Flags = setupAppFlags()
	return app
}

func main() {
	log.SetLevel(log.InfoLevel)
	if err := createApp().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
