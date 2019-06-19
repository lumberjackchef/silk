package cmd

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// RootDirectoryName ...
var RootDirectoryName = ".silk"

// ProjectMeta ...
type ProjectMeta struct {
	ProjectName string `json:"project_name"`
	InitDate    string `json:"init_date"`
	Version     string `json:"version"`
	Description string `json:"description"`
	ProjectURL  string `json:"url"`
}

// Check provides basic error checking & logging
// TODO: implement better logging/error handling. Panic is not the only way to handle an error
//       need to implement recovers as well
// TODO: Move all error handling to an errors.go file/package?
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// SilkMetaFile provides project metadata in an easy to consume format
func SilkMetaFile() ProjectMeta {
	var fileData ProjectMeta

	// Open, check, & defer closing of the meta data file
	jsonFile, jsonFileErr := os.Open(SilkRoot() + "/.silk/meta.json")
	Check(jsonFileErr)
	defer jsonFile.Close()

	// Get the []byte version of the json data
	byteValue, byteValueErr := ioutil.ReadAll(jsonFile)
	Check(byteValueErr)

	// Transform the []byte data into usable struct data
	jsonDataErr := json.Unmarshal(byteValue, &fileData)
	Check(jsonDataErr)

	return fileData
}

// SilkRoot returns the project root directory path
func SilkRoot() string {
	currentWorkingDirectory, currentWorkingDirectoryErr := os.Getwd()
	Check(currentWorkingDirectoryErr)

	returnPath, walkUpErr := walkUp(currentWorkingDirectory, RootDirectoryName)
	Check(walkUpErr)

	return returnPath
}

// TODO: combine with global checkWalkUp
// walkUp allows us to walk up the file tree looking for a certain file name as an anchor, returns the directory path of the anchor
func walkUp(currentPath string, directoryName string) (string, error) {
	readCurrentPath, readCurrentPathErr := os.Open(currentPath)
	Check(readCurrentPathErr)
	defer readCurrentPath.Close()

	filesInCurrentDir, filesInCurrentDirErr := readCurrentPath.Readdir(-1)
	Check(filesInCurrentDirErr)

	for _, file := range filesInCurrentDir {
		if file.Name() == directoryName {
			return currentPath, nil
		}
	}

	// Checks if we're at the root, returns an error if true
	// TODO: Make sure this works with all filesystem types including containerized environments
	// Mac: '/', Windows: 'C:\', Linux: '/', (Docker: '/'?)
	userRoot, userRootErr := filepath.Match("/", currentPath)
	Check(userRootErr)

	if userRoot {
		// TODO: this should invoke the "not a silk project" line instead
		return "", errors.New("Root directory reached")
	}

	// Recursion
	recursiveWalk, recursiveWalkErr := walkUp(filepath.Dir(currentPath), directoryName)
	Check(recursiveWalkErr)

	return recursiveWalk, nil
}
