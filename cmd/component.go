package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/lumberjackchef/silk/helper"
	"github.com/urfave/cli"
)

// TODO: Change new component command to `$ silk component new {name}`?

// Component either lists the component index or creates a new component if an arg is provided
func Component() cli.Command {
	return cli.Command{
		Name:    "component",
		Aliases: []string{"c"},
		Usage:   "If no arguments, lists all components in the current project. If a name is supplied, this will either move to the component, clone from remote & move to the component, or it will create a new component of name [name]",
		Action: func(c *cli.Context) error {
			helper.CommandAction(func() {
				cWarning := color.New(color.FgYellow).SprintFunc()
				cNotice := color.New(color.FgGreen).SprintFunc()

				// Lists index of components
				if len(helper.ComponentIndex()) > 0 {
					fmt.Println(cNotice("\tComponents in project " + helper.SilkMetaFile().ProjectName + ":"))
					for _, component := range helper.ComponentIndex() {
						fmt.Println("\t\t" + component)
					}
				} else {
					fmt.Printf("\t%s There are no components in the current project.\n", cWarning("Warning:"))
				}
			})
			return nil
		},
		Subcommands: []cli.Command{
			NewComponent(),
			RemoveComponent(),
		},
	}
}
