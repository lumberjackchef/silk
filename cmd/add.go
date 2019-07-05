package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/lumberjackchef/silk/helper"
	"github.com/urfave/cli"
)

// Add allows for the addition of changes to the commit buffer
func Add() cli.Command {
	cNotice := color.New(color.FgGreen).SprintFunc()
	cWarning := color.New(color.FgYellow).SprintFunc()

	return cli.Command{
		Name:  "add",
		Usage: "Adds a file or files to the current commit buffer (note: must be project/component root relative)",
		Action: func(c *cli.Context) error {
			helper.CommandAction(func() {
				if c.NArg() > 0 {
					var addIt bool
					var isDir bool
					fileName := fmt.Sprintf(c.Args().Get(0))

					fi, err := os.Stat(fileName)
					helper.Check(err)

					if fi.Mode().IsDir() {
						isDir = true
					}

					if !isDir {
						for _, file := range helper.UnstagedFilesList() {
							if fileName == file {
								addIt = true
							}
						}
					}
					if addIt && !isDir {
						helper.AddFileToCommitBuffer(fileName)
						fmt.Printf("\n\t%s added to the commit buffer!\n", cNotice(fileName))
					} else if isDir {
						var files []string

						err := filepath.Walk(fileName, func(path string, info os.FileInfo, err error) error {
							// TODO: need to get clean `path` with `SilkRoot()`
							for _, file := range helper.UnstagedFilesList() {
								if path == file {
									// TODO: add the file to commit buffer
									helper.AddFileToCommitBuffer(path)
									files = append(files, path)
								}
							}
							return nil
						})
						helper.Check(err)

						if len(files) > 0 {
							fmt.Printf("\n\t Files added to the commit buffer: ")
							// Print the names of the files added via `files`
						} else {
							fmt.Println("\n\t" + cWarning("Warning: ") + "no files with changes present\n")
						}
					} else {
						var inBuffer bool
						// TODO: Log an error here
						for _, file := range helper.FilesInCommitBuffer() {
							if fileName == file {
								inBuffer = true
							}
						}

						if inBuffer {
							fmt.Printf("\n\tFile(s) have already been addedto the commit buffer!\n\n")
						} else {
							fmt.Printf("\n\tNot a valid file\n\n")
						}
					}
				} else {
					err := errors.New("no files specified")
					fmt.Println("\n\t" + cWarning("Warning: ") + err.Error() + ". Please specify which files/changes you wish to add.\n")
				}
			})
			return nil
		},
		Subcommands: []cli.Command{
			InteractiveAddition(),
		},
	}
}
