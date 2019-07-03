package cmd

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/lumberjackchef/silk/helper"
	"github.com/urfave/cli"
)

// TODO: add remove command

// Remove allows for the addition of changes to the commit buffer
func Remove() cli.Command {
	cNotice := color.New(color.FgGreen).SprintFunc()
	cWarning := color.New(color.FgYellow).SprintFunc()

	return cli.Command{
		Name:  "remove",
		Usage: "Removes a file or files from the current commit buffer",
		Action: func(c *cli.Context) error {
			helper.CommandAction(func() {
				if c.NArg() > 0 {
					fmt.Printf("\t%s\n", cNotice("Coming Soon!"))
					// TODO: check against changes in the commit buffer
					// TODO: remove file/directory changes to commit buffer
				} else {
					err := errors.New("no files specified")
					fmt.Println("\n\t" + cWarning("Warning: ") + err.Error() + ". Please specify which files/changes you wish to remove.\n")
				}
			})
			return nil
		},
		Subcommands: []cli.Command{
			InteractiveAddition(),
		},
	}
}
