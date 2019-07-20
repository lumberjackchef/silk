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
					if err != nil {
						fmt.Println(cWarning("\n\tError") + ": unable to read file information")
						fmt.Print("\n")
					}

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
					if err != nil {
						fmt.Println(cWarning("\n\tError") + ": unable to open commit buffer file")
						fmt.Print("\n")
					}
					defer bufferFile.Close()

					byteValue, err := ioutil.ReadAll(bufferFile)
					if err != nil {
						fmt.Println(cWarning("\n\tError") + ": unable to read buffer file data")
						fmt.Print("\n")
					}

					err = json.Unmarshal(byteValue, &commitBuffer)
					if err != nil {
						fmt.Println(cWarning("\n\tError") + ": unable to unmarshal commit buffer data")
						fmt.Print("\n")
					}

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
						if err != nil {
							fmt.Println(cWarning("\n\tError") + ": unable to walk directory")
							fmt.Print("\n")
						}
					}

					// override oldChanges & add newChanges to buffer
					commitBuffer.Changes = newChanges

					commitBufferJSON, err := json.MarshalIndent(commitBuffer, "", "	")
					if err != nil {
						fmt.Println(cWarning("\n\tError") + ": unable to marshal commit buffer data")
						fmt.Print("\n")
					}

					err = ioutil.WriteFile(helper.SilkRoot()+"/.silk/commit/buffer", []byte(string(commitBufferJSON)+"\n"), 0766)
					if err != nil {
						fmt.Println(cWarning("\n\tError") + ": unable to write to commit buffer")
						fmt.Print("\n")
					}

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
