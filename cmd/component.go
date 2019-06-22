package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/lumberjackchef/silk/helper"
	"github.com/urfave/cli"
)

// ComponentCommand creates new component if arg provided, shows the index if not
func ComponentCommand(c *cli.Context) error {
	helper.CommandAction(func() {
		cWarning := color.New(color.FgYellow).SprintFunc()
		cNotice := color.New(color.FgGreen).SprintFunc()

		if c.NArg() > 0 {
			// Parameterized & lower-cased version of the user input string
			componentName := fmt.Sprintf(strings.Join(strings.Split(strings.ToLower(c.Args().Get(0)), " "), "-"))
			componentConfigDirectory := helper.SilkRoot() + "/" + componentName + "/.silk-component"

			helper.CreateComponentsListFile(componentName, componentConfigDirectory)
		} else {
			// Lists index of components
			if len(GetComponentIndex()) > 0 {
				fmt.Println(cNotice("\tComponents in project " + helper.SilkMetaFile().ProjectName + ":"))
				for _, component := range GetComponentIndex() {
					fmt.Println("\t\t" + component)
				}
			} else {
				fmt.Printf("\t%s There are no components in the current project.\n", cWarning("Warning:"))
			}
		}
	})
	return nil
}

// GetComponentIndex returns a slice that lists all the components in .silk/components.json
func GetComponentIndex() []string {
	var componentFileData helper.ComponentList

	// Open, Check, & defer closing of the component list file
	componentJSONFile, componentJSONFileErr := os.Open(helper.SilkRoot() + "/.silk/components.json")
	helper.Check(componentJSONFileErr)
	defer componentJSONFile.Close()

	// Get the []byte version of the json data
	componentByteValue, componentByteValueErr := ioutil.ReadAll(componentJSONFile)
	helper.Check(componentByteValueErr)

	// Transform the []byte data into usable struct data
	componentJSONDataErr := json.Unmarshal(componentByteValue, &componentFileData)
	helper.Check(componentJSONDataErr)

	return componentFileData.ComponentList
}

// RemoveComponent removes the component from the Components List (obvi)
func RemoveComponent(c *cli.Context) error {
	helper.CommandAction(func() {
		// Colors setup
		cWarning := color.New(color.FgYellow).SprintFunc()
		cBold := color.New(color.Bold).SprintFunc()

		if c.NArg() > 0 {
			// Parameterized & lower-cased version of the user input string
			componentName := fmt.Sprintf(strings.Join(strings.Split(strings.ToLower(c.Args().Get(0)), " "), "-"))

			componentDirectory := helper.SilkRoot() + "/" + componentName
			_, componentErr := os.Stat(componentDirectory)

			if componentErr == nil {
				var cList helper.ComponentList

				// Remove the component files
				os.RemoveAll(componentDirectory)

				// remove the component from the components.json list
				// open & read components file
				componentsList, componentsListErr := os.Open(helper.SilkRoot() + "/.silk/components.json")
				helper.Check(componentsListErr)
				defer componentsList.Close()

				// get byte value of components file
				byteValue, byteValueErr := ioutil.ReadAll(componentsList)
				helper.Check(byteValueErr)

				// unmarshall the data into the ComponentList struct
				cListErr := json.Unmarshal(byteValue, &cList)
				helper.Check(cListErr)

				// remove the component from the list []string
				for index, value := range cList.ComponentList {
					if value == componentName {
						cList.ComponentList = append(cList.ComponentList[:index], cList.ComponentList[index+1:]...)
					}
				}

				// transform back to JSON
				componentsListJSON, componentsListJSONErr := json.MarshalIndent(cList, "", " ")
				helper.Check(componentsListJSONErr)

				// Write the version change to the file
				componentFileWriteErr := ioutil.WriteFile(helper.SilkRoot()+"/.silk/components.json", []byte(string(componentsListJSON)+"\n"), 0766)
				helper.Check(componentFileWriteErr)

				// Confirmation message
				fmt.Println("\tComponent " + cBold(componentName) + " successfully removed!")
			} else {
				fmt.Printf("\t%s No component matching that name exists.\n", cWarning("Error:"))
			}
		} else {
			fmt.Printf("\t%s No component name specified.\n", cWarning("Error:"))
		}
	})
	return nil
}
