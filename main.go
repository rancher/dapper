package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/ibuildthecloud/dapper/file"
)

func main() {
	exit := func(err error) {
		if err != nil {
			logrus.Fatal(err)
		}
	}

	app := cli.NewApp()
	app.Author = "@ibuildthecloud"
	app.EnableBashCompletion = true
	app.Usage = `Docker build wrapper

	Dockerfile variables

	DAPPER_SOURCE          The destination directory in the container to bind/copy the source
	DAPPER_CP              The location in the host to find the source
	DAPPER_OUTPUT          The files you want copied to the host in CP mode
	DAPPER_DOCKER_SOCKET   Whether the Docker socket should be bound in
	DAPPER_RUN_ARGS        Args to add to the docker run command when building
	DAPPER_ENV             Env vars that should be copied into the build`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Value: "Dockerfile.dapper",
			Usage: "Dockerfile to build from",
		},
		cli.BoolFlag{
			Name:  "socket, k",
			Usage: "Bind in the Docker socket",
		},
		cli.StringFlag{
			Name:   "mode, m",
			Value:  "auto",
			Usage:  "Execution mode for Dapper bind/cp/auto",
			EnvVar: "DAPPER_MODE",
		},
		cli.StringFlag{
			Name:  "directory, C",
			Value: ".",
			Usage: "The directory in which to run, --file is relative to this",
		},
		cli.BoolFlag{
			Name:  "shell, s",
			Usage: "Launch a shell",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Print debugging",
		},
	}
	app.Action = func(c *cli.Context) {
		exit(run(c))
	}

	exit(app.Run(os.Args))
}

func run(c *cli.Context) error {
	if c.Bool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	dir := c.String("directory")
	mode := c.String("mode")
	shell := c.Bool("shell")
	socket := c.Bool("socket")

	if err := os.Chdir(dir); err != nil {
		return fmt.Errorf("Failed to change to directory %s: %v", dir, err)
	}

	dapperFile, err := file.Lookup(c.String("file"))
	if err != nil {
		return err
	}

	dapperFile.SetSocket(socket)

	if shell {
		return dapperFile.Shell(mode, c.Args())
	}

	return dapperFile.Run(mode, c.Args())
}
