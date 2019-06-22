package main

import (
	"log"
	"os"

	"github.com/lumberjackchef/silk/cmd"
	"github.com/urfave/cli"
)

func main() {
	// Application setup
	app := cli.NewApp()
	app.Name = "silk"
	app.Usage = "A modern version control paradigm for service oriented architectures."

	// Allows for bash completion of commands & subcommands
	app.EnableBashCompletion = true

	app.Commands = []cli.Command{
		cmd.New(),
		cmd.Status(),
		cmd.Clone(),
		cmd.Component(),
		cmd.Version(),
		cmd.Add(),
		cmd.Commit(),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
