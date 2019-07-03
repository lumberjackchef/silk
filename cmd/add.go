package cmd

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/lumberjackchef/silk/helper"
	"github.com/urfave/cli"
)

// Add allows for the addition of changes to the commit buffer
func Add() cli.Command {
	cNotice := color.New(color.FgGreen).SprintFunc()
	cWarning := color.New(color.FgYellow).SprintFunc()

	return cli.Command{
		Name:  "add",
		Usage: "Adds a file or files to the current commit buffer",
		Action: func(c *cli.Context) error {
			helper.CommandAction(func() {
				if c.NArg() > 0 {
					fmt.Printf("\t%s\n", cNotice("Coming Soon!"))
					// TODO: check against changes not in the commit buffer
					// TODO: see if the listed item is a valid file/directory
					// TODO: add file/directory changes to commit buffer
				} else {
					err := errors.New("no files specified")
					fmt.Println("\n\t" + cWarning("Warning: ") + err.Error() + ". Please specify which files/changes you wish to add.\n")
				}
			})
			return nil
		},
		Subcommands: []cli.Command{
			InteractiveAddition(),
		},
	}
}
