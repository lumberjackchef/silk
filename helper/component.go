package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/fatih/color"
)

/*
Components List
*/

// CreateComponentList creates the root project components list
func CreateComponentList(projectName string) {
	// Create a blank components list file
	componentsList, componentsListErr := os.Create(RootDirectoryName + "/components.json")
	Check(componentsListErr)
	defer componentsList.Close()

	// Creates the components data & writes to the file
	componentsListData, _ := json.MarshalIndent(
		&ComponentList{
			ProjectName:   projectName,
			ComponentList: []string{},
		},
		"",
		"  ",
	)

	_, componentsListWriteError := componentsList.WriteString(string(componentsListData) + "\n")
	Check(componentsListWriteError)
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

/*
Component Meta Data
*/

// CreateComponentMetaFile creates {silk component root}/.silk-component/meta.json file
func CreateComponentMetaFile(componentName string, componentConfigDirectory string) {
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

/*
Component Index
*/

// ComponentIndex returns a slice that lists all the components in .silk/components.json
func ComponentIndex() []string {
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
