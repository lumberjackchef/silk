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

// RemoveComponent removes the component from the Components List (obvi)
func RemoveComponent() cli.Command {
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
						componentsList, err := os.Open(helper.SilkRoot() + "/.silk/components.json")
						if err != nil {
							fmt.Println(cWarning("\n\tError") + ": unable to open components list")
							fmt.Print("\n")
						}
						defer componentsList.Close()

						// get byte value of components file
						byteValue, err := ioutil.ReadAll(componentsList)
						if err != nil {
							fmt.Println(cWarning("\n\tError") + ": unable to read components list")
							fmt.Print("\n")
						}

						// unmarshall the data into the ComponentList struct
						err = json.Unmarshal(byteValue, &cList)
						if err != nil {
							fmt.Println(cWarning("\n\tError") + ": unable to unmarshall data into ComponentsList struct")
							fmt.Print("\n")
						}

						// remove the component from the list []string
						for index, value := range cList.ComponentList {
							if value == componentName {
								cList.ComponentList = append(cList.ComponentList[:index], cList.ComponentList[index+1:]...)
							}
						}

						// transform back to JSON
						componentsListJSON, err := json.MarshalIndent(cList, "", " ")
						if err != nil {
							fmt.Println(cWarning("\n\tError") + ": unable to marshall component list data")
							fmt.Print("\n")
						}

						// Write the version change to the file
						err = ioutil.WriteFile(helper.SilkRoot()+"/.silk/components.json", []byte(string(componentsListJSON)+"\n"), 0766)
						if err != nil {
							fmt.Println(cWarning("\n\tError") + ": unable to write to components file")
							fmt.Print("\n")
						}

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
