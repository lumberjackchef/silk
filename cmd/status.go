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
				cError := color.New(color.FgRed).SprintFunc()
				cTitle := color.New(color.FgCyan).SprintFunc()
				cItalic := color.New(color.Italic).SprintFunc()

				// Print status
				fmt.Printf("\n\t%s    "+helper.SilkMetaFile().ProjectName+"\n", cTitle("Project:"))

				if helper.IsComponentOrRoot() == "component" {
					fmt.Printf("\t%s  %s\n", cTitle("Component:"), cItalic("Coming Soon!"))
				}

				fmt.Printf("\t%s	    %s\n\n", cTitle("Branch:"), cItalic("Coming Soon!"))
				// Print all file names in the commit buffer
				files := helper.FilesInCommitBuffer()
				if len(files) > 0 {
					fmt.Println("Changes to be committed:")
				}
				for index, file := range files {
					fmt.Println("\t\t" + cNotice(file))

					if index == len(files)-1 {
						fmt.Print("\n")
					}
				}

				// Print files with changes that are _not_ in the commit buffer
				files = helper.UnstagedFilesList()
				if len(files) > 0 {
					fmt.Println("Changes not staged for commit:")
				}
				for index, file := range files {
					fmt.Println("\t\t" + cError(file))

					if index == len(files)-1 {
						fmt.Print("\n")
					}
				}
			})
			return nil
		},
	}
}
