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

	log.Infoln("before action")

	// TBD
	return nil
}

func after(context *cli.Context) error {
	log.Infoln("after action")

	// TBD
	return nil
}

func start(context *cli.Context) error {
	log.Println("entry point -> begin")

	rc := core.Main(context)

	log.Printf("entry point -> end")
	return rc
}

func setupAppFlags() []cli.Flag {

	return []cli.Flag{
		cli.StringFlag{
			Name:   core.HostName,
			Usage:  core.HostUsage,
			Value:  core.HostValue,
			EnvVar: core.HostEnvVar,
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
