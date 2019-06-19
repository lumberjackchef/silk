package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// ComponentList .silk/components.json file structure
type ComponentList struct {
	ProjectName   string   `json:"project_name"`
	ComponentList []string `json:"component_list"`
}

// ComponentMeta .silk/../{component}/.silk-component/meta.json file structure
type ComponentMeta struct {
	ProjectName   string `json:"project_name"`
	ComponentName string `json:"component_name"`
	InitDate      string `json:"init_date"`
	Version       string `json:"version"`
	Description   string `json:"description"`
}

// SilkComponentRoot returns the component root directory path
func SilkComponentRoot() string {
	currentWorkingDirectory, currentWorkingDirectoryErr := os.Getwd()
	Check(currentWorkingDirectoryErr)

	returnPath, walkUpErr := walkUp(currentWorkingDirectory, ".silk-component")
	Check(walkUpErr)

	return returnPath
}

// AddToSilkComponentList adds componentName to the list of components in .silk/components.json
func AddToSilkComponentList(componentName string) error {
	var componentFileData ComponentList

	// Open, Check, & defer closing of the component list file
	componentJSONFile, componentJSONFileErr := os.Open(SilkRoot() + "/.silk/components.json")
	Check(componentJSONFileErr)
	defer componentJSONFile.Close()

	// Get the []byte version of the json data
	componentByteValue, componentByteValueErr := ioutil.ReadAll(componentJSONFile)
	Check(componentByteValueErr)

	// Transform the []byte data into usable struct data
	componentJSONDataErr := json.Unmarshal(componentByteValue, &componentFileData)
	Check(componentJSONDataErr)

	// Append the component name
	componentFileData.ComponentList = append(componentFileData.ComponentList, componentName)
	componentFileJSONData, componentFileJSONDataErr := json.MarshalIndent(componentFileData, "", "  ")
	Check(componentFileJSONDataErr)

	// Write the component addition to the file
	componentFileJSONDataWriteErr := ioutil.WriteFile(SilkRoot()+"/.silk/components.json", []byte(string(componentFileJSONData)+"\n"), 0766)
	Check(componentFileJSONDataWriteErr)

	return nil
}

// ComponentCommand creates new component if arg provided, shows the index if not
func ComponentCommand(c *cli.Context) error {
	CommandAction(func() {
		cWarning := color.New(color.FgYellow).SprintFunc()
		cNotice := color.New(color.FgGreen).SprintFunc()

		if c.NArg() > 0 {
			// Parameterized & lower-cased version of the user input string
			componentName := fmt.Sprintf(strings.Join(strings.Split(strings.ToLower(c.Args().Get(0)), " "), "-"))
			componentConfigDirectory := SilkRoot() + "/" + componentName + "/.silk-component"

			CreateComponentsListFile(componentName, componentConfigDirectory)
		} else {
			// Lists index of components
			if len(GetComponentIndex()) > 0 {
				fmt.Println(cNotice("\tComponents in project " + SilkMetaFile().ProjectName + ":"))
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

// CreateComponentsListFile creates {silk root}/components.json file
func CreateComponentsListFile(componentName string, componentConfigDirectory string) {
	// Colors setup
	cWarning := color.New(color.FgYellow).SprintFunc()
	cBold := color.New(color.Bold).SprintFunc()

	// Component tracking directory. This Checks if the directory exists, creates it if not.
	_, componentConfigErr := os.Stat(componentConfigDirectory)
	if os.IsNotExist(componentConfigErr) {
		// creates the '{component}/.silk-component directory as well as the {component} directory if one or both don't exist
		os.MkdirAll(componentConfigDirectory, 0766)

		// Creates the project meta json file
		componentMeta, componentMetaErr := os.Create(componentConfigDirectory + "/meta.json")
		Check(componentMetaErr)
		defer componentMeta.Close()

		// Creates the project metadata & writes to the file
		dT := time.Now().String()
		componentMetaData, _ := json.MarshalIndent(&ComponentMeta{
			ProjectName:   SilkMetaFile().ProjectName,
			ComponentName: componentName,
			InitDate:      dT,
			Version:       "0.0.0",
		}, "", "  ")
		_, componentMetaWriteErr := componentMeta.WriteString(string(componentMetaData) + "\n")
		Check(componentMetaWriteErr)

		// Adds component to component list file
		AddToSilkComponentList(componentName)

		// Confirmation message
		fmt.Printf("\tNew component %s created!\n", cBold(componentName))
	} else {
		fmt.Printf("\t%s component %s already exists!\n", cWarning("Warning:"), cBold(componentName))
	}
}

// GetComponentIndex returns a slice that lists all the components in .silk/components.json
func GetComponentIndex() []string {
	var componentFileData ComponentList

	// Open, Check, & defer closing of the component list file
	componentJSONFile, componentJSONFileErr := os.Open(SilkRoot() + "/.silk/components.json")
	Check(componentJSONFileErr)
	defer componentJSONFile.Close()

	// Get the []byte version of the json data
	componentByteValue, componentByteValueErr := ioutil.ReadAll(componentJSONFile)
	Check(componentByteValueErr)

	// Transform the []byte data into usable struct data
	componentJSONDataErr := json.Unmarshal(componentByteValue, &componentFileData)
	Check(componentJSONDataErr)

	return componentFileData.ComponentList
}

// RemoveComponent removes the component from the Components List (obvi)
func RemoveComponent(c *cli.Context) error {
	CommandAction(func() {
		// Colors setup
		cWarning := color.New(color.FgYellow).SprintFunc()
		cBold := color.New(color.Bold).SprintFunc()

		if c.NArg() > 0 {
			// Parameterized & lower-cased version of the user input string
			componentName := fmt.Sprintf(strings.Join(strings.Split(strings.ToLower(c.Args().Get(0)), " "), "-"))

			componentDirectory := SilkRoot() + "/" + componentName
			_, componentErr := os.Stat(componentDirectory)

			if componentErr == nil {
				var cList ComponentList

				// Remove the component files
				os.RemoveAll(componentDirectory)

				// remove the component from the components.json list
				// open & read components file
				componentsList, componentsListErr := os.Open(SilkRoot() + "/.silk/components.json")
				Check(componentsListErr)
				defer componentsList.Close()

				// get byte value of components file
				byteValue, byteValueErr := ioutil.ReadAll(componentsList)
				Check(byteValueErr)

				// unmarshall the data into the ComponentList struct
				cListErr := json.Unmarshal(byteValue, &cList)
				Check(cListErr)

				// remove the component from the list []string
				for index, value := range cList.ComponentList {
					if value == componentName {
						cList.ComponentList = append(cList.ComponentList[:index], cList.ComponentList[index+1:]...)
					}
				}

				// transform back to JSON
				componentsListJSON, componentsListJSONErr := json.MarshalIndent(cList, "", " ")
				Check(componentsListJSONErr)

				// Write the version change to the file
				componentFileWriteErr := ioutil.WriteFile(SilkRoot()+"/.silk/components.json", []byte(string(componentsListJSON)+"\n"), 0766)
				Check(componentFileWriteErr)

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
