package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/lumberjackchef/silk/helper"
	"github.com/urfave/cli"
)

// Status shows the state of the current project or component
func Status() cli.Command {
	return cli.Command{
		Name:    "status",
		Aliases: []string{"s"},
		Usage:   "Get the status of the current project and/or component.",
		Action: func(c *cli.Context) error {
			helper.CommandAction(func() {
				cNotice := color.New(color.FgGreen).SprintFunc()

				// Print status
				fmt.Printf("\t%s "+helper.SilkMetaFile().ProjectName+"\n\n", cNotice("Project:"))

				// Print all file names in the commit buffer
				files := helper.FilesInCommitBuffer()
				for index, file := range files {
					fmt.Println("\t\t" + cNotice(file))

					if index == len(files)-1 {
						fmt.Print("\n")
					}
				}

				// TODO: Print files with changes that are _not_ in the commit buffer
			})
			return nil
		},
	}
}
