package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/lumberjackchef/silk/cmd"
	"github.com/lumberjackchef/silk/helper"
	"github.com/urfave/cli"
)

func main() {
	// Colors setup
	cNotice := color.New(color.FgGreen).SprintFunc()

	// Application setup
	app := cli.NewApp()
	app.Name = "silk"
	app.Usage = "A modern version control paradigm for service oriented architectures."

	// Allows for bash completion of commands & subcommands
	app.EnableBashCompletion = true

	app.Commands = []cli.Command{
		// TODO: redo this so that cmd.New, etc returns all this in one line
		// 			 these should be named like the command line args (cmd.New, cmd.Status, etc)
		{
			Name:    "new",
			Aliases: []string{"n"},
			Usage:   "Create a new silk project",
			Action:  cmd.CreateNewProject,
		},
		{
			Name:    "status",
			Aliases: []string{"s"},
			Usage:   "Get the status of the current project and/or component.",
			Action:  cmd.StatusCommand,
		},
		{
			Name:  "clone",
			Usage: "Copies down the project root from remote and sets up all default branches & remotes.",
			Action: func(c *cli.Context) error {
				helper.CommandAction(func() { fmt.Printf("\t%s\n", cNotice("Coming Soon!")) })
				return nil
			},
		},
		{
			Name:    "component",
			Aliases: []string{"c"},
			Usage:   "If no arguments, lists all components in the current project. If a name is supplied, this will either move to the component, clone from remote & move to the component, or it will create a new component of name [name]",
			Action:  cmd.ComponentCommand,
			Subcommands: []cli.Command{
				{
					Name:   "remove",
					Usage:  "remove an existing component",
					Action: cmd.RemoveComponent,
				},
			},
		},
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "Lists or edits the current version of the project",
			Action:  cmd.PrintOrChangeVersion,
		},
		cmd.Add(),
		cmd.Commit(),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
