package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/lumberjackchef/silk/helper"
	"github.com/urfave/cli"
)

// Commit adds commit buffer data to a commit file
func Commit() cli.Command {
	cNotice := color.New(color.FgGreen).SprintFunc()

	return cli.Command{
		Name:  "commit",
		Usage: "Tags the current commit buffer and resets all file statuses",
		Action: func(c *cli.Context) error {
			// TODO: add commit buffer to some sort of hash function
			// TODO: return as a commit hash of some kind when committing
			// TODO: add commit message ability
			helper.CommandAction(func() { fmt.Printf("\t%s\n", cNotice("Coming Soon!")) })
			return nil
		},
	}
}
