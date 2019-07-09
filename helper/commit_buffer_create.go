package helper

import (
	"bufio"
	"encoding/json"
	"os"
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
		},
		"",
		"	",
	)

	_, commitBufferWriteErr := commitBuffer.WriteString(string(commitBufferData) + "\n")
	Check(commitBufferWriteErr)
}
