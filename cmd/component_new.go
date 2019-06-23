package cmd

import (
	"fmt"
	"strings"

	"github.com/lumberjackchef/silk/helper"
	"github.com/urfave/cli"
)

// NewComponent removes the component from the Components List (obvi)
func NewComponent() cli.Command {
	return cli.Command{
		Name:  "new",
		Usage: "add a new component",
		Action: func(c *cli.Context) error {
			helper.CommandAction(func() {
				// Parameterized & lower-cased version of the user input string
				componentName := fmt.Sprintf(strings.Join(strings.Split(strings.ToLower(c.Args().Get(0)), " "), "-"))
				componentConfigDirectory := helper.SilkRoot() + "/" + componentName + "/.silk-component"

				helper.CreateComponentMetaFile(componentName, componentConfigDirectory)
			})
			return nil
		},
	}
}
