/*
Package cmd is a package for all the root commands for Silk
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/lumberjackchef/silk/helper"
	"github.com/urfave/cli"
)

// TODO: Create initial commit buffer
// TODO: Determine how to save/record project & component meta info silently (pass to remote & down to local, etc)

// New creates a new silk project
func New() cli.Command {
	return cli.Command{
		Name:    "new",
		Aliases: []string{"n"},
		Usage:   "Create a new silk project",
		Action: func(c *cli.Context) error {
			helper.CommandAction(func() {
				// Colors setup
				cNotice := color.New(color.FgGreen).SprintFunc()
				cWarning := color.New(color.FgYellow).SprintFunc()

				// Project tracking folder. This checks if the folder exists, creates it if not.
				// TODO: implement checking similar to SilkRoot() to ensure we're not within the bounds of another project
				if _, err := os.Stat(helper.RootDirectoryName); os.IsNotExist(err) {
					if c.NArg() > 0 {
						// Creates the silk directory
						os.Mkdir(helper.RootDirectoryName, 0766)

						projectName := fmt.Sprintf(c.Args().Get(0))

						// Create relevant project files
						helper.CreateRootMetaFile(projectName)
						helper.CreateComponentList(projectName)

						// Confirmation message
						fmt.Println("\tNew project " + cNotice(fmt.Sprintf(c.Args().Get(0))) + " created!")
					} else {
						fmt.Printf("\t%s No project name specified!\n", cWarning("Warning:"))
					}
				} else {
					fmt.Printf("\t%s this is an existing silk project!\n", cWarning("Warning:"))
				}
			})
			return nil
		},
	}
}
