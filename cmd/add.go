package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/lumberjackchef/silk/helper"
	"github.com/urfave/cli"
)

// Add allows for the addition of changes to the commit buffer
func Add() cli.Command {
	cWarning := color.New(color.FgYellow).SprintFunc()

	return cli.Command{
		Name:  "add",
		Usage: "Adds a file or files to the current commit buffer (note: must be project/component root relative)",
		Action: func(c *cli.Context) error {
			helper.CommandAction(func() {
				if c.NArg() > 0 {
					var addIt bool
					var commitBuffer helper.RootCommitBuffer
					var changes []helper.FileChange
					fileName := fmt.Sprintf(c.Args().Get(0))
					oldChanges := helper.CommitBuffer().Changes

					for _, file := range helper.UnstagedFilesList() {
						if fileName == file {
							addIt = true
						}
					}

					if addIt {
						// TODO: currently this only works for a single file addition
						// 			 update to accomodate for a whole directory of changes
						bufferFile, err := os.Open(helper.SilkRoot() + "/" + helper.RootDirectoryName + "/commit/buffer")
						helper.Check(err)
						defer bufferFile.Close()

						byteValue, err := ioutil.ReadAll(bufferFile)
						helper.Check(err)

						err = json.Unmarshal(byteValue, &commitBuffer)
						helper.Check(err)

						changedFile, err := os.Open(helper.SilkRoot() + "/" + fileName)
						helper.Check(err)
						defer changedFile.Close()

						scanner := bufio.NewScanner(changedFile)
						line := 0
						for scanner.Scan() {
							line = line + 1

							fileChange := helper.FileChange{
								FileName:   fileName,
								LineNumber: line,
								Text:       scanner.Text(),
							}

							changes = append(changes, fileChange)
						}
						changes = append(changes, oldChanges...)
						commitBuffer.Changes = changes

						commitBufferJSON, err := json.MarshalIndent(commitBuffer, " ", "")
						helper.Check(err)

						err = ioutil.WriteFile(helper.SilkRoot()+"/.silk/commit/buffer", []byte(string(commitBufferJSON)+"\n"), 0766)
						helper.Check(err)

						fmt.Println("\n\tAdded new file!")
					} else {
						// TODO: Print an error here
						fmt.Println("File has already been added or does not exist!")
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
