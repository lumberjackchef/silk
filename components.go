package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
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

// GetComponentIndex duh
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
