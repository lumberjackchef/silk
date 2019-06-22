/*
Package cmd is a package for all the root commands for Silk
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/lumberjackchef/silk/helper"
	"github.com/urfave/cli"
)

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

						// Creates the project meta json file
						projectMeta, projectMetaErr := os.Create(helper.RootDirectoryName + "/meta.json")
						helper.Check(projectMetaErr)
						defer projectMeta.Close()

						// Creates the project metadata & writes to the file
						dT := time.Now().String()
						projectMetaData, _ := json.MarshalIndent(
							&helper.ProjectMeta{
								ProjectName: fmt.Sprintf(c.Args().Get(0)),
								InitDate:    dT,
								Version:     "0.0.0",
							},
							"",
							"  ",
						)

						_, projectMetaWriteErr := projectMeta.WriteString(string(projectMetaData) + "\n")
						helper.Check(projectMetaWriteErr)

						// Create a blank components list file
						componentsList, componentsListErr := os.Create(helper.RootDirectoryName + "/components.json")
						helper.Check(componentsListErr)
						defer componentsList.Close()

						// Creates the components data & writes to the file
						componentsListData, _ := json.MarshalIndent(
							&helper.ComponentList{
								ProjectName:   fmt.Sprintf(c.Args().Get(0)),
								ComponentList: []string{},
							},
							"",
							"  ",
						)

						_, componentsListWriteError := componentsList.WriteString(string(componentsListData) + "\n")
						helper.Check(componentsListWriteError)

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

// CreateNewProject creates new project if arg provided & not already a project, errors if not
