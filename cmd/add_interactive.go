package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/lumberjackchef/silk/helper"
	"github.com/urfave/cli"
)

// InteractiveAddition allows for the addition of a large number of changes to the commit buffer
func InteractiveAddition() cli.Command {
	cNotice := color.New(color.FgGreen).SprintFunc()

	return cli.Command{
		Name:    "interactive",
		Aliases: []string{"i"},
		Usage:   "Adds/removes a large amount of changes to the current commit buffer via an interactive interface",
		Action: func(c *cli.Context) error {
			helper.CommandAction(func() { fmt.Printf("\t%s\n", cNotice("Coming Soon!")) })
			return nil
		},
	}
}
