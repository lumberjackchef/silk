package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/fatih/color"
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
	check(currentWorkingDirectoryErr)

	returnPath, walkUpErr := walkUp(currentWorkingDirectory, ".silk-component")
	check(walkUpErr)

	return returnPath
}

// AddToSilkComponentList adds componentName to the list of components in .silk/components.json
func AddToSilkComponentList(componentName string) error {
	var componentFileData ComponentList

	// Open, check, & defer closing of the component list file
	componentJSONFile, componentJSONFileErr := os.Open(SilkRoot() + "/.silk/components.json")
	check(componentJSONFileErr)
	defer componentJSONFile.Close()

	// Get the []byte version of the json data
	componentByteValue, componentByteValueErr := ioutil.ReadAll(componentJSONFile)
	check(componentByteValueErr)

	// Transform the []byte data into usable struct data
	componentJSONDataErr := json.Unmarshal(componentByteValue, &componentFileData)
	check(componentJSONDataErr)

	// Append the component name
	componentFileData.ComponentList = append(componentFileData.ComponentList, componentName)
	componentFileJSONData, componentFileJSONDataErr := json.MarshalIndent(componentFileData, "", "  ")
	check(componentFileJSONDataErr)

	// Write the component addition to the file
	componentFileJSONDataWriteErr := ioutil.WriteFile(SilkRoot()+"/.silk/components.json", []byte(string(componentFileJSONData)+"\n"), 0766)
	check(componentFileJSONDataWriteErr)

	return nil
}

// CreateComponentsListFile creates {silk root}/components.json file
func CreateComponentsListFile(componentName string, componentConfigDirectory string) {
	// Colors setup
	cWarning := color.New(color.FgYellow).SprintFunc()
	cBold := color.New(color.Bold).SprintFunc()

	// Component tracking directory. This checks if the directory exists, creates it if not.
	_, componentConfigErr := os.Stat(componentConfigDirectory)
	if os.IsNotExist(componentConfigErr) {
		// creates the '{component}/.silk-component directory as well as the {component} directory if one or both don't exist
		os.MkdirAll(componentConfigDirectory, 0766)

		// Creates the project meta json file
		componentMeta, componentMetaErr := os.Create(componentConfigDirectory + "/meta.json")
		check(componentMetaErr)
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
		check(componentMetaWriteErr)

		// Adds component to component list file
		AddToSilkComponentList(componentName)

		// Confirmation message
		fmt.Printf("\tNew component %s created!\n", cBold(componentName))
	} else {
		fmt.Printf("\t%s component %s already exists!\n", cWarning("Warning:"), cBold(componentName))
	}
}

// RemoveComponent removes the component from the Components List (obvi)
func RemoveComponent(componentName string) {
	// Colors setup
	cWarning := color.New(color.FgYellow).SprintFunc()
	cBold := color.New(color.Bold).SprintFunc()

	componentDirectory := SilkRoot() + "/" + componentName
	_, componentErr := os.Stat(componentDirectory)

	if componentErr == nil {
		var cList ComponentList

		// Remove the component files
		os.RemoveAll(componentDirectory)

		// remove the component from the components.json list
		// open & read components file
		componentsList, componentsListErr := os.Open(SilkRoot() + "/.silk/components.json")
		check(componentsListErr)
		defer componentsList.Close()

		// get byte value of components file
		byteValue, byteValueErr := ioutil.ReadAll(componentsList)
		check(byteValueErr)

		// unmarshall the data into the ComponentList struct
		cListErr := json.Unmarshal(byteValue, &cList)
		check(cListErr)

		// remove the component from the list []string
		for index, value := range cList.ComponentList {
			if value == componentName {
				cList.ComponentList = append(cList.ComponentList[:index], cList.ComponentList[index+1:]...)
			}
		}

		// transform back to JSON
		componentsListJSON, componentsListJSONErr := json.MarshalIndent(cList, "", " ")
		check(componentsListJSONErr)

		// Write the version change to the file
		componentFileWriteErr := ioutil.WriteFile(SilkRoot()+"/.silk/components.json", []byte(string(componentsListJSON)+"\n"), 0766)
		check(componentFileWriteErr)

		// Confirmation message
		fmt.Println("\tComponent " + cBold(componentName) + " successfully removed!")
	} else {
		fmt.Printf("\t%s No component matching that name exists.\n", cWarning("Error:"))
	}
}

// GetComponentIndex returns a slice that lists all the components in .silk/components.json
func GetComponentIndex() []string {
	var componentFileData ComponentList

	// Open, check, & defer closing of the component list file
	componentJSONFile, componentJSONFileErr := os.Open(SilkRoot() + "/.silk/components.json")
	check(componentJSONFileErr)
	defer componentJSONFile.Close()

	// Get the []byte version of the json data
	componentByteValue, componentByteValueErr := ioutil.ReadAll(componentJSONFile)
	check(componentByteValueErr)

	// Transform the []byte data into usable struct data
	componentJSONDataErr := json.Unmarshal(componentByteValue, &componentFileData)
	check(componentJSONDataErr)

	return componentFileData.ComponentList
}
