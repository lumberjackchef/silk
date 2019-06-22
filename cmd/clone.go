package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/lumberjackchef/silk/helper"
	"github.com/urfave/cli"
)

// Clone copies down the project root from remote
func Clone() cli.Command {
	cNotice := color.New(color.FgGreen).SprintFunc()

	return cli.Command{
		Name:  "clone",
		Usage: "Copies down the project root from remote",
		Action: func(c *cli.Context) error {
			helper.CommandAction(func() { fmt.Printf("\t%s\n", cNotice("Coming Soon!")) })
			return nil
		},
	}
}
