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
	cWarning := color.New(color.FgYellow).SprintFunc()

	// Create a blank components list file
	componentsList, err := os.Create(RootDirectoryName + "/components.json")
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to create components file")
		fmt.Print("\n")
	}
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

	_, err = componentsList.WriteString(string(componentsListData) + "\n")
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to write components file data")
		fmt.Print("\n")
	}
}

// AddToSilkComponentList adds componentName to the list of components in .silk/components.json
func AddToSilkComponentList(componentName string) error {
	var componentFileData ComponentList
	cWarning := color.New(color.FgYellow).SprintFunc()

	// Open, Check, & defer closing of the component list file
	componentJSONFile, err := os.Open(SilkRoot() + "/.silk/components.json")
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to open components file")
		fmt.Print("\n")
	}
	defer componentJSONFile.Close()

	// Get the []byte version of the json data
	componentByteValue, err := ioutil.ReadAll(componentJSONFile)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to read components file")
		fmt.Print("\n")
	}

	// Transform the []byte data into usable struct data
	err = json.Unmarshal(componentByteValue, &componentFileData)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to unmarshal component file byte values")
		fmt.Print("\n")
	}

	// Append the component name
	componentFileData.ComponentList = append(componentFileData.ComponentList, componentName)
	componentFileJSONData, err := json.MarshalIndent(componentFileData, "", "  ")
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to unable to marshal component file json")
		fmt.Print("\n")
	}

	// Write the component addition to the file
	err = ioutil.WriteFile(SilkRoot()+"/.silk/components.json", []byte(string(componentFileJSONData)+"\n"), 0766)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to write components file data")
		fmt.Print("\n")
	}

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
	_, err := os.Stat(componentConfigDirectory)
	if os.IsNotExist(err) {
		// creates the '{component}/.silk-component directory as well as the {component} directory if one or both don't exist
		os.MkdirAll(componentConfigDirectory, 0766)

		// Creates the project meta json file
		componentMeta, err := os.Create(componentConfigDirectory + "/meta.json")
		if err != nil {
			fmt.Println(cWarning("\n\tError") + ": unable to create component meta file")
			fmt.Print("\n")
		}
		defer componentMeta.Close()

		// Creates the project metadata & writes to the file
		dT := time.Now().String()
		componentMetaData, _ := json.MarshalIndent(&ComponentMeta{
			ProjectName:   SilkMetaFile().ProjectName,
			ComponentName: componentName,
			InitDate:      dT,
			Version:       "0.0.0",
		}, "", "  ")
		_, err = componentMeta.WriteString(string(componentMetaData) + "\n")
		if err != nil {
			fmt.Println(cWarning("\n\tError") + ": unable to write to component meta file")
			fmt.Print("\n")
		}

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
	cWarning := color.New(color.FgYellow).SprintFunc()

	// Open, Check, & defer closing of the component list file
	componentJSONFile, err := os.Open(SilkRoot() + "/.silk/components.json")
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to open components file")
		fmt.Print("\n")
	}
	defer componentJSONFile.Close()

	// Get the []byte version of the json data
	componentByteValue, err := ioutil.ReadAll(componentJSONFile)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to read components file data")
		fmt.Print("\n")
	}

	// Transform the []byte data into usable struct data
	err = json.Unmarshal(componentByteValue, &componentFileData)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to unmarshal components file data")
		fmt.Print("\n")
	}

	return componentFileData.ComponentList
}
