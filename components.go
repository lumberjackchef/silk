package main

import (
	"encoding/json"
	// "fmt"
	"io/ioutil"
	"os"
)

// ComponentList ...
type ComponentList struct {
	ProjectName   string   `json:"project_name"`
	ComponentList []string `json:"component_list"`
}

// ComponentMeta ...
type ComponentMeta struct {
	ProjectName   string `json:"project_name"`
	ComponentName string `json:"component_name"`
	InitDate      string `json:"init_date"`
	Version       string `json:"version"`
	Description   string `json:"description"`
}

// AddToSilkComponentList Project metadata helper
func AddToSilkComponentList(componentName string) error {
	var componentFileData ComponentList

	// Open, check, & defer closing of the component list file
	// TODO: Add SilkRoot() here
	componentJSONFile, componentJSONFileErr := os.Open(".silk/components.json")
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
	// TODO: Add SilkRoot() here
	componentFileJSONDataWriteErr := ioutil.WriteFile(".silk/components.json", []byte(string(componentFileJSONData)+"\n"), 0766)
	check(componentFileJSONDataWriteErr)

	return nil
}

// GetComponentIndex ...
func GetComponentIndex() []string {
	var componentFileData ComponentList

	// Open, check, & defer closing of the component list file
	// TODO: Add SilkRoot() here
	componentJSONFile, componentJSONFileErr := os.Open(".silk/components.json")
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
