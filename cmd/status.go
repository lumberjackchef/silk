package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/lumberjackchef/silk/helper"
	"github.com/urfave/cli"
)

// StatusCommand returns the status of the current commit buffer
func StatusCommand(c *cli.Context) error {
	helper.CommandAction(func() {
		cNotice := color.New(color.FgGreen).SprintFunc()

		// Print status
		fmt.Printf("\t%s "+helper.SilkMetaFile().ProjectName+"\n\n", cNotice("Project:"))

		if helper.IsComponentOrRoot() == "component" {
			os.Chdir(helper.SilkComponentRoot())
		} else {
			os.Chdir(helper.SilkRoot())
		}

		currentWorkingDirectory, _ := os.Getwd()

		// File list
		files := helper.ComposeFileList(currentWorkingDirectory)

		// Print the file status
		helper.ListFilesInCommitBuffer(files)
	})
	return nil
}
