package helper

import (
	"errors"
	"os"
	"path/filepath"
)

// RootDirectoryName ...
const RootDirectoryName = ".silk"

// SilkRoot returns the project root directory path
func SilkRoot() string {
	currentWorkingDirectory, currentWorkingDirectoryErr := os.Getwd()
	Check(currentWorkingDirectoryErr)

	returnPath, walkUpErr := walkUp(currentWorkingDirectory, RootDirectoryName)
	Check(walkUpErr)

	return returnPath
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

// SilkComponentRoot returns the component root directory path
func SilkComponentRoot() string {
	currentWorkingDirectory, currentWorkingDirectoryErr := os.Getwd()
	Check(currentWorkingDirectoryErr)

	returnPath, walkUpErr := walkUp(currentWorkingDirectory, ".silk-component")
	Check(walkUpErr)

	return returnPath
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

// TODO: combine with global helper checkWalkUp
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
