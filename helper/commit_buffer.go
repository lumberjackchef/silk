package helper

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// ChangesNotInCommitBuffer returns a []FilesChange list of all changes that are not in the current _root_ buffer
func ChangesNotInCommitBuffer() []FileChange {
	var allFileChanges []FileChange
	cWarning := color.New(color.FgYellow).SprintFunc()

	// Get each line from every file, creating new FileChange{}s
	for _, file := range ListAllFiles() {
		bufferFile, err := os.Open(file)
		if err != nil {
			fmt.Println(cWarning("\n\tError") + ": unable to open file")
			fmt.Print("\n")
		}
		defer bufferFile.Close()

		scanner := bufio.NewScanner(bufferFile)
		line := 0
		for scanner.Scan() {
			line = line + 1

			fileChange := FileChange{
				FileName:   file,
				LineNumber: line,
				Text:       scanner.Text(),
			}

			allFileChanges = append(allFileChanges, fileChange)
		}
	}

	allFileChanges = append(allFileChanges, LatestCommit().Changes...)
	allFileChanges = append(allFileChanges, CommitBuffer().Changes...)

	unique := make(map[FileChange]bool)
	returnSlice := []FileChange{}
	if len(allFileChanges) > 0 {
		for _, change := range allFileChanges {
			if !unique[change] {
				returnSlice = append(returnSlice, change)
				unique[change] = true
			} else if unique[change] {
				var index int
				var changed bool
				for i, value := range returnSlice {
					if value == change {
						index = i
						changed = true
					}
				}
				if changed {
					returnSlice[index] = returnSlice[len(returnSlice)-1]
					returnSlice = returnSlice[:len(returnSlice)-1]
				}
			}
		}
	}
	// 		// Append when:
	// 		// if we hit a matching file & line number _and_ the text is different
	// 		// if we hit a matching file, the line numbers don't match, but the text is the same
	// 		// if we hit a matching file & the line number doesn't exist

	// 		// maybe consider adding ChangeType to FileChange (would be "line" or "text")
	// 		// 		"text" would be the type for additions
	// 		// how do we get only the lines that changed when it might be a multi-line change?
	// 		// are there additional new line numbers with text?

	// //

	return returnSlice
}

// UnstagedFilesList returns a list of all files with changes not in the commit buffer or the commit history
func UnstagedFilesList() []string {
	var files []string

	for _, change := range ChangesNotInCommitBuffer() {
		files = append(files, change.FileName)
	}

	// Sanitizes the return to only have unique elements
	files = UniqueNonEmptyElementsOf(files)
	return files
}

// TODO: combine the above & below into a more flexible function

// FilesInCommitBuffer returns a []string list of all files with changes in the current _root_ buffer
func FilesInCommitBuffer() []string {
	var files []string

	if len(CommitBuffer().Changes) < 0 {
		return []string{}
	}

	for _, change := range CommitBuffer().Changes {
		files = append(files, change.FileName)
	}

	// Sanitizes the return to only have unique elements
	files = UniqueNonEmptyElementsOf(files)
	return files
}

// FilesNotInCommitBuffer returns a list of files that are not currently in the commit buffer
func FilesNotInCommitBuffer() []string {
	var files []string
	var filePath string
	cWarning := color.New(color.FgYellow).SprintFunc()
	addFile := true
	index := 1

	if IsComponentOrRoot() == "component" {
		filePath = SilkComponentRoot()
	} else if IsComponentOrRoot() == "root" {
		filePath = SilkRoot()
	}

	err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		// first, find all files not currently in the commit buffer at all
		fileName := strings.Replace(path, SilkRoot()+"/", "", 1)

		// comparison to commit buffer files list
		if index != 1 && !info.IsDir() {
			for _, file := range FilesInCommitBuffer() {
				if file == fileName {
					addFile = false
					break
				} else {
					addFile = true
				}
			}

			if IsNotExcluded(fileName, info) && addFile {
				files = append(files, fileName)
			}
		}

		index++
		return nil
	})

	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to walk files not in the commit buffer")
		fmt.Print("\n")
	}

	return files
}
