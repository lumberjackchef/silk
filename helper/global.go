/*
Package helper provides global helper funcs for Silk core and all related packages
*/
package helper

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

// RootDirectoryName ...
const RootDirectoryName = ".silk"

// Check provides basic error checking & logging
// TODO: implement better logging/error handling. Panic is not the only way to handle an error
//       need to implement recovers as well
// TODO: Move all error handling to an errors.go file/package?
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// CommandAction checks if this is a silk project before running a command
func CommandAction(f func()) string {
	cWarning := color.New(color.FgYellow).SprintFunc()

	if IsComponentOrRoot() == "false" {
		fmt.Printf("\t%s this is not a silk project! To create a new silk project, run `$ silk new`\n", cWarning("Warning:"))
	} else {
		f()
	}
	return ""
}

// IsComponentOrRoot returns component, root, or false
func IsComponentOrRoot() string {
	var partType string

	currentWorkingDirectory, currentWorkingDirectoryErr := os.Getwd()
	Check(currentWorkingDirectoryErr)

	checkReturnPath, checkReturnPathErr := CheckWalkUp(currentWorkingDirectory)
	Check(checkReturnPathErr)

	if checkReturnPath == "component" {
		partType = "component"
	} else if checkReturnPath == "root" {
		partType = "root"
	} else {
		partType = "false"
	}

	return partType
}

// CheckWalkUp is a separate function solely for recursion
func CheckWalkUp(currentPath string) (string, error) {
	readCurrentPath, readCurrentPathErr := os.Open(currentPath)
	Check(readCurrentPathErr)
	defer readCurrentPath.Close()

	filesInCurrentDir, filesInCurrentDirErr := readCurrentPath.Readdir(-1)
	Check(filesInCurrentDirErr)

	for _, file := range filesInCurrentDir {
		if file.Name() == RootDirectoryName {
			return "root", nil
		} else if file.Name() == ".silk-component" {
			return "component", nil
		}
	}

	// Checks if we're at the root, returns an error if true
	// TODO: Make sure this works with all filesystem types including containerized environments
	// Mac: '/', Windows: 'C:\', Linux: '/', (Docker: '/'?)
	userRoot, userRootErr := filepath.Match("/", currentPath)
	Check(userRootErr)

	if userRoot {
		return "", nil
	}

	// Recursion
	recursiveWalk, recursiveWalkErr := CheckWalkUp(filepath.Dir(currentPath))
	Check(recursiveWalkErr)

	return recursiveWalk, nil
}
