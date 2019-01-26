package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/rancher/dapper/file"
)

var (
	VERSION = "0.0.0"
)

func main() {
	exit := func(err error) {
		if err == file.ErrSkipBuild {
			logrus.Infof("Build not supported on this architecture")
			os.Exit(42)
		}
		if err != nil {
			logrus.Fatal(err)
		}
	}

	app := cli.NewApp()
	app.Author = "Rancher Labs"
	app.EnableBashCompletion = true
	app.Version = VERSION
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
		cli.BoolFlag{
			Name:  "no-out, O",
			Usage: "Do not copy the output back (in --mode cp)",
		},
		cli.BoolFlag{
			Name:  "build",
			Usage: "Perform Dapperfile build",
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
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "Make Docker build quieter",
		},
		cli.BoolFlag{
			Name:  "keep",
			Usage: "Don't remove the container that was used to build",
		},
		cli.BoolFlag{
			Name:   "no-context, X",
			Usage:  "send Dockerfile via stdin to docker build command",
			EnvVar: "DAPPER_NO_CONTEXT",
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
	shell := c.Bool("shell")
	build := c.Bool("build")

	if err := os.Chdir(dir); err != nil {
		return fmt.Errorf("Failed to change to directory %s: %v", dir, err)
	}

	dapperFile, err := file.Lookup(c.String("file"))
	if err != nil {
		return err
	}

	dapperFile.Mode = c.String("mode")
	dapperFile.Socket = c.Bool("socket")
	dapperFile.NoOut = c.Bool("no-out")
	dapperFile.Quiet = c.Bool("quiet")
	dapperFile.Keep = c.Bool("keep")
	dapperFile.NoContext = c.Bool("no-context")

	if shell {
		return dapperFile.Shell(c.Args())
	}

	if build {
		return dapperFile.Build(c.Args())
	}

	return dapperFile.Run(c.Args())
}
