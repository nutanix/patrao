package main //import github.com/nutanix/patrao/

import (
	"log"
	"os"

	core "github.com/nutanix/patrao/internal/app/upgrade_agent"
	"github.com/urfave/cli"
)

func before(context *cli.Context) error {
	log.Println("before action")

	// TBD
	return nil
}

func after(context *cli.Context) error {
	log.Println("after action")

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
	}
}

func main() {
	app := cli.NewApp()

	app.Name = "Patrao Upgrade Agent"
	app.Usage = "Upgrade service for automatically upgrade docker based solutions"
	app.Before = before
	app.After = after
	app.Action = start

	setupAppFlags(app)

	if err := app.Run(os.Args); err != nil {
		log.Fatal("fatal error!")
	}
}
