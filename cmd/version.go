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

// PrintOrChangeVersion ...
func PrintOrChangeVersion(c *cli.Context) string {
	// Colors setup
	cNotice := color.New(color.FgGreen).SprintFunc()

	var metaData helper.ProjectMeta

	// Open, check, & defer closing of the meta data file
	metaFile, metaFileErr := os.Open(helper.SilkRoot() + "/.silk/meta.json")
	helper.Check(metaFileErr)
	defer metaFile.Close()

	// Get the []byte version of the json data
	byteValue, byteValueErr := ioutil.ReadAll(metaFile)
	helper.Check(byteValueErr)

	// Transform the []byte data into usable struct data
	metaDataErr := json.Unmarshal(byteValue, &metaData)
	helper.Check(metaDataErr)

	if c.NArg() > 0 {
		// Change the version & transform back to actual json
		metaData.Version = fmt.Sprintf(c.Args().Get(0))
		metaDataJSON, metaDataJSONErr := json.MarshalIndent(metaData, "", " ")
		helper.Check(metaDataJSONErr)

		// Write the version change to the file
		metaFileWriteErr := ioutil.WriteFile(helper.SilkRoot()+"/.silk/meta.json", []byte(string(metaDataJSON)+"\n"), 0766)
		helper.Check(metaFileWriteErr)

		// Confirmation message
		fmt.Println("\tVersion successfully updated to " + cNotice(metaData.Version) + "!")
	} else {
		// If the user just wants to check the version and not change it
		fmt.Println("\t" + cNotice(metaData.Version))
	}

	return ""
}
