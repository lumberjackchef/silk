package helper

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

// RootDirectoryName ...
const RootDirectoryName = ".silk"

// SilkRoot returns the project root directory path
func SilkRoot() string {
	cWarning := color.New(color.FgYellow).SprintFunc()
	currentWorkingDirectory, _ := os.Getwd()

	returnPath, err := walkUp(currentWorkingDirectory, RootDirectoryName)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to walk up current directory to get the project root path")
		fmt.Print("\n")
	}

	return returnPath
}

// IsComponentOrRoot returns component, root, or false
func IsComponentOrRoot() string {
	cWarning := color.New(color.FgYellow).SprintFunc()
	var partType string

	currentWorkingDirectory, _ := os.Getwd()

	checkReturnPath, err := CheckWalkUp(currentWorkingDirectory)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to recursively check up the tree to determine if project or component")
		fmt.Print("\n")
	}

	if checkReturnPath == "component" {
		partType = "component"
	} else if checkReturnPath == "root" {
		partType = "root"
	} else {
		partType = "false"
	}

	return partType
}

// SilkComponentRoot returns the component root directory path
func SilkComponentRoot() string {
	cWarning := color.New(color.FgYellow).SprintFunc()
	currentWorkingDirectory, _ := os.Getwd()

	returnPath, err := walkUp(currentWorkingDirectory, ".silk-component")
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to walk up current tree to get the component root path")
		fmt.Print("\n")
	}

	return returnPath
}

// CheckWalkUp is a separate function solely for recursion
func CheckWalkUp(currentPath string) (string, error) {
	cWarning := color.New(color.FgYellow).SprintFunc()

	readCurrentPath, err := os.Open(currentPath)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to read current path")
		fmt.Print("\n")
	}
	defer readCurrentPath.Close()

	filesInCurrentDir, err := readCurrentPath.Readdir(-1)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to read files in current directory")
		fmt.Print("\n")
	}

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
	userRoot, err := filepath.Match("/", currentPath)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to find system root")
		fmt.Print("\n")
	}

	if userRoot {
		return "", nil
	}

	// Recursion
	recursiveWalk, err := CheckWalkUp(filepath.Dir(currentPath))
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to recursively check up the tree")
		fmt.Print("\n")
	}

	return recursiveWalk, nil
}

// TODO: combine with global helper checkWalkUp
// walkUp allows us to walk up the file tree looking for a certain file name as an anchor, returns the directory path of the anchor
func walkUp(currentPath string, directoryName string) (string, error) {
	cWarning := color.New(color.FgYellow).SprintFunc()

	readCurrentPath, err := os.Open(currentPath)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to read current path")
		fmt.Print("\n")
	}
	defer readCurrentPath.Close()

	filesInCurrentDir, err := readCurrentPath.Readdir(-1)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to read files in current directory")
		fmt.Print("\n")
	}

	for _, file := range filesInCurrentDir {
		if file.Name() == directoryName {
			return currentPath, nil
		}
	}

	// Checks if we're at the root, returns an error if true
	// TODO: Make sure this works with all filesystem types including containerized environments
	// Mac: '/', Windows: 'C:\', Linux: '/', (Docker: '/'?): filepath.VolumeName(path)?
	userRoot, err := filepath.Match("/", currentPath)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to find system root")
		fmt.Print("\n")
	}

	if userRoot {
		// TODO: this should invoke the "not a silk project" line instead
		return "", errors.New("Root directory reached")
	}

	// Recursion
	recursiveWalk, err := walkUp(filepath.Dir(currentPath), directoryName)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to recursively check up the tree")
		fmt.Print("\n")
	}

	return recursiveWalk, nil
}
