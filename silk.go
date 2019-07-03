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
	// TODO: add system file edits & additions to enable full bash completion out of the gate
	app.EnableBashCompletion = true

	// TODO: find a way to add extensions without requiring the user to install go or rebuild or anything
	app.Commands = []cli.Command{
		cmd.New(),
		cmd.Clone(),
		cmd.Status(),
		cmd.Add(),
		cmd.Remove(),
		cmd.Commit(),
		cmd.Component(),
		cmd.Version(),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
