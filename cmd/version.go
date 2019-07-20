package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/lumberjackchef/silk/helper"
	"github.com/urfave/cli"
)

// Version lists or edits the current version of the project
func Version() cli.Command {
	return cli.Command{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Lists or edits the current version of the project",
		Action: func(c *cli.Context) error {
			helper.CommandAction(func() {
				// Colors setup
				cNotice := color.New(color.FgGreen).SprintFunc()
				cWarning := color.New(color.FgYellow).SprintFunc()
				var metaData helper.ProjectMeta

				// Open, check, & defer closing of the meta data file
				metaFile, err := os.Open(helper.SilkRoot() + "/.silk/meta.json")
				if err != nil {
					fmt.Println(cWarning("\n\tError") + ": unable to open project meta file")
					fmt.Print("\n")
				}
				defer metaFile.Close()

				// Get the []byte version of the json data
				byteValue, err := ioutil.ReadAll(metaFile)
				if err != nil {
					fmt.Println(cWarning("\n\tError") + ": unable to read project meta file")
					fmt.Print("\n")
				}

				// Transform the []byte data into usable struct data
				err = json.Unmarshal(byteValue, &metaData)
				if err != nil {
					fmt.Println(cWarning("\n\tError") + ": unable to unmarshal meta file information")
					fmt.Print("\n")
				}

				if c.NArg() > 0 {
					// Change the version & transform back to actual json
					metaData.Version = fmt.Sprintf(c.Args().Get(0))
					metaDataJSON, err := json.MarshalIndent(metaData, "", " ")
					if err != nil {
						fmt.Println(cWarning("\n\tError") + ": unable to marshal project meta data")
						fmt.Print("\n")
					}

					// Write the version change to the file
					err = ioutil.WriteFile(helper.SilkRoot()+"/.silk/meta.json", []byte(string(metaDataJSON)+"\n"), 0766)
					if err != nil {
						fmt.Println(cWarning("\n\tError") + ": unable to write project meta file")
						fmt.Print("\n")
					}

					// Confirmation message
					fmt.Println("\tVersion successfully updated to " + cNotice(metaData.Version) + "!")
				} else {
					// If the user just wants to check the version and not change it
					fmt.Println("\t" + cNotice(metaData.Version))
				}
			})
			return nil
		},
	}
}
