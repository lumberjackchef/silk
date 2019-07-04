package helper

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// CreateInitialCommitBuffer initiates the commit buffer and writes a new file
func CreateInitialCommitBuffer() {
	var allFileChanges []FileChange
	os.Chdir(SilkRoot() + "/")

	// Creates the commit buffer file
	createPathErr := os.MkdirAll(RootDirectoryName+"/commit", 0744)
	Check(createPathErr)
	commitBuffer, err := os.Create(RootDirectoryName + "/commit/buffer")
	Check(err)
	defer commitBuffer.Close()

	// TODO: meta data should be transferred on pull/push
	// Corral all the project meta data
	// Should return a JSON string of all the data in .silk/*
	// projectMeta := ""

	// Get each line from every file, creating new FileChange{}s
	files := ListAllFiles()
	for _, file := range files {
		bufferFile, err := os.Open(file)
		Check(err)
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

	// Creates the commit buffer data & writes to the file
	commitBufferData, _ := json.MarshalIndent(
		&RootCommitBuffer{
			ProjectName: SilkMetaFile().ProjectName,
			Changes:     allFileChanges,
		}, "", "  ",
	)

	_, commitBufferWriteErr := commitBuffer.WriteString(string(commitBufferData) + "\n")
	Check(commitBufferWriteErr)
}

// ChangesNotInCommitBuffer returns a []FilesChange list of all changes that are not in the current _root_ buffer
func ChangesNotInCommitBuffer() []FileChange {
	var allFileChanges []FileChange

	// Get each line from every file, creating new FileChange{}s
	for _, file := range ListAllFiles() {
		bufferFile, err := os.Open(file)
		Check(err)
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
	files = uniqueNonEmptyElementsOf(files)
	return files
}

// TODO: combine the above & below into a more flexible function

// FilesInCommitBuffer returns a []string list of all files with changes in the current _root_ buffer
func FilesInCommitBuffer() []string {
	var files []string

	for _, change := range CommitBuffer().Changes {
		files = append(files, change.FileName)
	}

	// Sanitizes the return to only have unique elements
	files = uniqueNonEmptyElementsOf(files)
	return files
}

// FilesNotInCommitBuffer returns a list of files that are not currently in the commit buffer
func FilesNotInCommitBuffer() []string {
	var files []string
	var filePath string
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
	Check(err)

	return files
}
