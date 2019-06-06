package main

import (
	"encoding/json"
	"errors"
	"fmt"
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
// TODO: implement better logging/error handling. Panic is not the only way to handle an error
//       need to implement recovers as well
// TODO: Move all error handling to an errors.go file/package?
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// SilkMetaFile provides project metadata in an easy to consume format
func SilkMetaFile() ProjectMeta {
	var fileData ProjectMeta

	// Open, check, & defer closing of the meta data file
	jsonFile, jsonFileErr := os.Open(SilkRoot() + "/.silk/meta.json")
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

// SilkRoot returns the project root directory path
func SilkRoot() string {
	currentWorkingDirectory, currentWorkingDirectoryErr := os.Getwd()
	check(currentWorkingDirectoryErr)

	returnPath, walkUpErr := walkUp(currentWorkingDirectory, rootDirectoryName)
	check(walkUpErr)

	return returnPath
}

// walkUp allows us to walk up the file tree looking for a certain file name as an anchor, returns the directory path of the anchor
func walkUp(currentPath string, directoryName string) (string, error) {
	readCurrentPath, readCurrentPathErr := os.Open(currentPath)
	check(readCurrentPathErr)
	defer readCurrentPath.Close()

	filesInCurrentDir, filesInCurrentDirErr := readCurrentPath.Readdir(-1)
	check(filesInCurrentDirErr)

	for _, file := range filesInCurrentDir {
		if file.Name() == directoryName {
			return currentPath, nil
		}
	}

	// Checks if we're at the root
	userRoot, userRootErr := filepath.Match("/", currentPath)
	check(userRootErr)

	// TODO: Return a better error without panic()-ing
	if userRoot {
		fmt.Println("warning: This is the root of the local machine or environment \n Please switch to the appropriate directory to continue")
		return "", errors.New("Root directory reached")
	}

	// Recursion
	recursiveWalk, recursiveWalkErr := walkUp(filepath.Dir(currentPath), directoryName)
	check(recursiveWalkErr)

	return recursiveWalk, nil
}
