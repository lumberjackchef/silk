package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/lumberjackchef/silk/helper"
	"github.com/urfave/cli"
)

// Remove allows for the addition of changes to the commit buffer
func Remove() cli.Command {
	// cNotice := color.New(color.FgGreen).SprintFunc()
	cWarning := color.New(color.FgYellow).SprintFunc()

	return cli.Command{
		Name:  "remove",
		Usage: "Removes a file or files from the current commit buffer",
		Action: func(c *cli.Context) error {
			helper.CommandAction(func() {
				if c.NArg() > 0 {
					var removeFile bool
					var isDir bool
					fileName := fmt.Sprintf(c.Args().Get(0))
					fi, err := os.Stat(fileName)
					helper.Check(err)

					if fi.Mode().IsDir() {
						isDir = true
					}

					if !isDir {
						for _, file := range helper.FilesInCommitBuffer() {
							if fileName == file {
								removeFile = true
							}
						}
					}

					var commitBuffer helper.RootCommitBuffer
					var newChanges []helper.FileChange
					var removedFiles []string
					oldChanges := helper.CommitBuffer().Changes

					// open buffer file & unmarshal
					bufferFile, err := os.Open(helper.SilkRoot() + "/" + helper.RootDirectoryName + "/commit/buffer")
					helper.Check(err)
					defer bufferFile.Close()

					byteValue, err := ioutil.ReadAll(bufferFile)
					helper.Check(err)

					err = json.Unmarshal(byteValue, &commitBuffer)
					helper.Check(err)

					if removeFile {
						for _, chg := range oldChanges {
							if chg.FileName != fileName {
								newChanges = append(newChanges, chg)
							}
						}
					} else if isDir {
						err = filepath.Walk(fileName, func(path string, info os.FileInfo, err error) error {
							if fi, _ := os.Stat(path); !fi.Mode().IsDir() {
								for _, file := range oldChanges {
									if file.FileName != path {
										newChanges = append(newChanges, file)
										removedFiles = append(removedFiles, path)
									}
								}
							}
							return nil
						})
						helper.Check(err)
					}

					// override oldChanges & add newChanges to buffer
					commitBuffer.Changes = newChanges

					commitBufferJSON, err := json.MarshalIndent(commitBuffer, "", "	")
					helper.Check(err)

					err = ioutil.WriteFile(helper.SilkRoot()+"/.silk/commit/buffer", []byte(string(commitBufferJSON)+"\n"), 0766)
					helper.Check(err)

					if removeFile {
						fmt.Println("\n\t" + cWarning(fileName) + " has been removed from the commit buffer!")
					} else if isDir && len(removedFiles) > 0 {
						removedFiles = helper.UniqueNonEmptyElementsOf(removedFiles)
						// TODO: order this list
						fmt.Println("\n\tFiles removed from the commit buffer: ")

						for _, file := range removedFiles {
							fmt.Println("\t\t" + cWarning(file))
						}
						fmt.Print("\n")
					}
				} else {
					err := errors.New("no files specified")
					fmt.Println("\n\t" + cWarning("Warning: ") + err.Error() + ". Please specify which files/changes you wish to remove.\n")
				}
			})
			return nil
		},
		Subcommands: []cli.Command{
			InteractiveAddition(),
		},
	}
}
