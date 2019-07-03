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
		cmd.Status(),
		cmd.Clone(),
		cmd.Component(),
		cmd.Version(),
		cmd.Add(),
		cmd.Remove(),
		cmd.Commit(),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
