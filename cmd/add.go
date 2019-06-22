package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/lumberjackchef/silk/helper"
	"github.com/urfave/cli"
)

// Add allows for the addition of files to the commit buffer
func Add() cli.Command {
	cNotice := color.New(color.FgGreen).SprintFunc()

	return cli.Command{
		Name:  "add",
		Usage: "Adds a file or files to the current commit buffer",
		Action: func(c *cli.Context) error {
			// TODO: Add a list of changes to files to a commit buffer
			// Check if the root commit exists
			// diff whole files first to eliminate unchanged files
			// diff changed files line by line
			// add changes to a commit buffer file
			// should be file name, line number, & actual code changes
			helper.CommandAction(func() { fmt.Printf("\t%s\n", cNotice("Coming Soon!")) })
			return nil
		},
	}
}
