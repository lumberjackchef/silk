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

// ComponentRemove removes the component from the Components List (obvi)
func ComponentRemove() cli.Command {
	return cli.Command{
		Name:  "remove",
		Usage: "remove an existing component",
		Action: func(c *cli.Context) error {
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
		},
	}
}