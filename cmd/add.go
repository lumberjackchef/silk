package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/lumberjackchef/silk/helper"
	"github.com/urfave/cli"
)

// Add allows for the addition of changes to the commit buffer
func Add() cli.Command {
	cNotice := color.New(color.FgGreen).SprintFunc()

	return cli.Command{
		Name:  "add",
		Usage: "Adds a file or files to the current commit buffer",
		Action: func(c *cli.Context) error {
			// TODO: add interactive addition for multifile additions
			// TODO: add simple, single file/folder change additions to the commit buffer
			helper.CommandAction(func() { fmt.Printf("\t%s\n", cNotice("Coming Soon!")) })
			return nil
		},
	}
}
