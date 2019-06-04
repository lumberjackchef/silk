package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

var rootDirectoryName = ".silk"

// ProjectMeta ...
type ProjectMeta struct {
	ProjectName string `json:"project_name"`
	InitDate    string `json:"init_date"`
	Version     string `json:"version"`
	Description string `json:"description"`
	ProjectURL  string `json:"url"`
}

// Error checking & logging
// TODO: impolement better logging/error handling. Panic is not the only way to handle an error
//       need to implement recovers as well
// TODO: Move all error handling to an errors.go file/package?
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// SilkMetaFile Project metadata helper
func SilkMetaFile() ProjectMeta {
	var fileData ProjectMeta

	// Open, check, & defer closing of the meta data file
	// TODO: Add SilkRoot() here
	jsonFile, jsonFileErr := os.Open(".silk/meta.json")
	check(jsonFileErr)
	defer jsonFile.Close()

	// Get the []byte version of the json data
	byteValue, byteValueErr := ioutil.ReadAll(jsonFile)
	check(byteValueErr)

	// Transform the []byte data into usable struct data
	jsonDataErr := json.Unmarshal(byteValue, &fileData)
	check(jsonDataErr)

	return fileData
}

// SilkRoot Project root helper
func SilkRoot() string {
	var returnPath string
	currentWorkingDirectory, currentWorkingDirectoryErr := os.Getwd()
	check(currentWorkingDirectoryErr)

	returnPath, walkUpErr := walkUp(currentWorkingDirectory)
	check(walkUpErr)

	return returnPath
}

func walkUp(currentPath string) (string, error) {
	var returnP, nextUp string

	walkErr := filepath.Walk(currentPath, func(path string, info os.FileInfo, err error) error {
		// Escape early if there's an error
		if err != nil {
			return err
		}

		// Abstract checking if we're at the root of the file system
		// TODO: need to ensure this works on containerized environments as well
		// TODO: loops infintiely if no silk project found
		systemRoot, systemRootErr := filepath.Match("/", currentPath)
		check(systemRootErr)

		// Checks if we're currently in the project root
		if info.Name() == rootDirectoryName {
			returnP = filepath.Dir(path)
			return nil

			// Checks if we're at the system root.
		} else if systemRoot {
			return errors.New("warning: This is the root of the local machine or environment. Please switch to the appropriate directory to continue")

			// Gives us the next level up in the tree to check.
		} else {
			nextUp = filepath.Dir(path)
			return nil
		}
	})

	// Handle recursion outside of the walk function
	if returnP == "" && walkErr == nil && nextUp != "" {
		returnP, walkErr := walkUp(nextUp)
		return returnP, walkErr
	}

	// Fallthrough. This is the only return likely to have an error state
	return returnP, walkErr
}
