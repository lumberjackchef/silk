package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ComposeFileList gathers a list of all relevant files for showing current status
func ComposeFileList(currentWorkingDirectory string) []string {
	var files []string
	var exclusionList []string

	if _, err := os.Stat(SilkRoot() + "/.silk-ignore"); !os.IsNotExist(err) {
		// Read the file
		// open the file for reading
		ignoreFile, ignoreFileErr := os.Open(SilkRoot() + "/.silk-ignore")
		check(ignoreFileErr)
		defer ignoreFile.Close()

		scanner := bufio.NewScanner(ignoreFile)
		for scanner.Scan() {
			exclusionList = append(exclusionList, scanner.Text())
		}
	}

	err := filepath.Walk(currentWorkingDirectory, func(path string, info os.FileInfo, err error) error {
		// Ignore non-project related files
		IsNotExcluded := true
		IsNotDotFile := !strings.HasPrefix(filepath.Base(path), ".")
		IsNotGit := !strings.Contains(filepath.Dir(path), ".git")
		IsNotSilkFiles := !strings.Contains(filepath.Dir(path), ".silk")
		DoesNotContainDotPath := IsNotDotFile && IsNotGit && IsNotSilkFiles

		if DoesNotContainDotPath {
			for _, matchPattern := range exclusionList {
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
	check(err)

	return files
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
