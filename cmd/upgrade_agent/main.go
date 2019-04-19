package main //import github.com/nutanix/patrao/

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

func setupAppFlags(app *cli.App) {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "host",
			Usage:  "daemon socket to connect to docker",
			Value:  "unix:///var/run/docker.sock",
			EnvVar: "PATRAO_DOCKER_HOST",
		},
		cli.BoolFlag{
			Name:  "run-once",
			Usage: "Run once now and exit",
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

	setupAppFlags(app)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
