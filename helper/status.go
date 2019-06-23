package helper

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CustomExclusionList returns the list of files user has decided to exclude from status/commit based on contents of .silk-ignore
func CustomExclusionList() []string {
	var exclusionList []string

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

// ListFilesInCommitBuffer lists all files in the project in the commit buffer
func ListFilesInCommitBuffer(files []string) {
	// TODO: needs to be updated in some way to indicate status of being in the current buffer
	// AWAITING: creation of a commit buffer
	for index, file := range files {
		fmt.Println("\t\t" + file)

		if index == len(files)-1 {
			fmt.Print("\n")
		}
	}
}

// ListAllFiles returns an array of file names based on project root
func ListAllFiles() []string {
	var files []string

	if IsComponentOrRoot() == "component" {
		os.Chdir(SilkComponentRoot())
	} else {
		os.Chdir(SilkRoot())
	}

	currentWorkingDirectory, _ := os.Getwd()

	err := filepath.Walk(currentWorkingDirectory, func(path string, info os.FileInfo, err error) error {
		// Ignore non-project related files
		IsNotExcluded := true
		IsNotGit := !strings.Contains(filepath.Dir(path), ".git")
		IsNotSilkFiles := !strings.Contains(filepath.Dir(path), ".silk")
		DoesNotContainDotPath := IsNotGit && IsNotSilkFiles

		if DoesNotContainDotPath {
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
		}

		if !info.IsDir() && DoesNotContainDotPath && IsNotExcluded {
			files = append(files, strings.Replace(path, SilkRoot()+"/", "", 1))
		}
		return nil
	})
	Check(err)

	return files
}
