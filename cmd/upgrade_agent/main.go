package main

import (
	"os"

	core "github.com/nutanix/patrao/internal/app/upgrade_agent"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func before(context *cli.Context) error {
	if context.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
	}
	// TBD
	return nil
}

func after(context *cli.Context) error {
	// TBD
	return nil
}

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
		cli.BoolFlag{
			Name:  core.RunOnceName,
			Usage: core.RunOnceUsage,
		},
	}
}

func main() {
	log.SetLevel(log.InfoLevel)

	app := cli.NewApp()

	app.Name = "Patrao Upgrade Agent"
	app.Usage = "Upgrade service for automatically upgrade docker based solutions"
	app.Before = before
	app.After = after
	app.Action = start

	app.Flags = setupAppFlags()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
