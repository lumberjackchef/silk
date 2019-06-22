package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/fatih/color"
)

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
