package helper

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// CustomExclusionList returns the list of files user has decided to exclude from status/commit based on contents of .silk-ignore
func CustomExclusionList() []string {
	var exclusionList []string

	// TODO: Fix this to be component/root agnostic
	if _, err := os.Stat(SilkRoot() + "/.silk-ignore"); !os.IsNotExist(err) {
		ignoreFile, ignoreFileErr := os.Open(SilkRoot() + "/.silk-ignore")
		Check(ignoreFileErr)
		defer ignoreFile.Close()

		scanner := bufio.NewScanner(ignoreFile)
		for scanner.Scan() {
			exclusionList = append(exclusionList, scanner.Text())
		}
	}

	return exclusionList
}

// IsNotExcluded takes a path & returns a boolean showing whether a file is excluded or not
func IsNotExcluded(path string, info os.FileInfo) bool {
	IsNotExcluded := true
	IsNotGit := !strings.Contains(filepath.Dir(path), ".git")
	IsNotSilkFiles := !strings.Contains(filepath.Dir(path), ".silk")
	PrimaryNotExcluded := IsNotGit && IsNotSilkFiles && !info.IsDir()

	if PrimaryNotExcluded {
		for _, matchPattern := range CustomExclusionList() {
			if IsNotExcluded {
				// NOTE: May need to fix this to be more specific
				// 			 This may exclude intentionally added files
				// 			 Need to support match types somehow
				// 			 Take a look at .gitignore patterns and see what we can glean from that
				IsNotExcluded = !strings.Contains(path, matchPattern)
			} else {
				break
			}
		}
	} else {
		IsNotExcluded = false
	}

	return IsNotExcluded
}

// ListAllFiles returns an array of file names based on project root
func ListAllFiles() []string {
	var files []string
	var currentWorkingDirectory string

	if IsComponentOrRoot() == "component" {
		currentWorkingDirectory = SilkComponentRoot()
	} else {
		currentWorkingDirectory = SilkRoot()
	}

	err := filepath.Walk(currentWorkingDirectory, func(path string, info os.FileInfo, err error) error {
		if IsNotExcluded(path, info) {
			files = append(files, strings.Replace(path, currentWorkingDirectory+"/", "", 1))
		}
		return nil
	})
	Check(err)

	return files
}
