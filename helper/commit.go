package helper

import (
	"bufio"
	"encoding/json"
	"os"
)

// CreateInitialCommitBuffer initiates the commit buffer and writes a new file
func CreateInitialCommitBuffer() {
	var allFileChanges []FileChange

	// TODO: Add SilkRoot()-like chacking here
	os.Chdir(SilkRoot()) // Not 100% that this will work

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

	files := ListAllFiles()

	// Get each line from every file, creating new FileChange{}s
	for _, file := range files {
		bufferFile, err := os.Open(file)
		Check(err)
		defer bufferFile.Close()

		scanner := bufio.NewScanner(bufferFile)
		line := 1
		for scanner.Scan() {
			fileChange := FileChange{
				FileName:   file,
				LineNumber: line,
				Text:       scanner.Text(),
			}

			allFileChanges = append(allFileChanges, fileChange)
			line++
		}
	}

	// Creates the commit buffer data & writes to the file
	commitBufferData, _ := json.MarshalIndent(
		&RootCommitBuffer{
			ProjectName: SilkMetaFile().ProjectName,
			Changes:     allFileChanges,
		},
		"",
		"  ",
	)

	_, commitBufferWriteErr := commitBuffer.WriteString(string(commitBufferData) + "\n")
	Check(commitBufferWriteErr)
}

// FilesInCommitBuffer returns a []string list of all files in the current _root_ buffer
func FilesInCommitBuffer() []string {
	// TODO:
	var files []string
	changes := CommitBufferFile().Changes

	for _, change := range changes {
		for _, file := range files {
			if file != change.FileName {
				files = append(files, change.FileName)
			}
		}
	}

	return files
}
