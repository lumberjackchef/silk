package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

// Checks if this is a silk project before running a command
func commandAction(f func()) string {
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
	check(currentWorkingDirectoryErr)

	checkReturnPath, checkReturnPathErr := checkWalkUp(currentWorkingDirectory)
	check(checkReturnPathErr)

	if checkReturnPath == "component" {
		partType = "component"
	} else if checkReturnPath == "root" {
		partType = "root"
	} else {
		partType = "false"
	}

	return partType
}

// need a separate function solely for recursion
func checkWalkUp(currentPath string) (string, error) {
	readCurrentPath, readCurrentPathErr := os.Open(currentPath)
	check(readCurrentPathErr)
	defer readCurrentPath.Close()

	filesInCurrentDir, filesInCurrentDirErr := readCurrentPath.Readdir(-1)
	check(filesInCurrentDirErr)

	for _, file := range filesInCurrentDir {
		if file.Name() == rootDirectoryName {
			return "root", nil
		} else if file.Name() == ".silk-component" {
			return "component", nil
		}
	}

	// Checks if we're at the root, returns an error if true
	// TODO: Make sure this works with all filesystem types including containerized environments
	// Mac: '/', Windows: 'C:\', Linux: '/', (Docker: '/'?)
	userRoot, userRootErr := filepath.Match("/", currentPath)
	check(userRootErr)

	if userRoot {
		return "", nil
	}

	// Recursion
	recursiveWalk, recursiveWalkErr := checkWalkUp(filepath.Dir(currentPath))
	check(recursiveWalkErr)

	return recursiveWalk, nil
}
