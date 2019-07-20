package helper

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
)

// CreateInitialCommitBuffer initiates the commit buffer and writes a new file
func CreateInitialCommitBuffer() {
	var allFileChanges []FileChange
	cWarning := color.New(color.FgYellow).SprintFunc()
	os.Chdir(SilkRoot() + "/")

	// Creates the commit buffer file
	err := os.MkdirAll(RootDirectoryName+"/commit", 0744)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to create project commit directory")
		fmt.Print("\n")
	}

	commitBuffer, err := os.Create(RootDirectoryName + "/commit/buffer")
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to create initial project commit buffer")
		fmt.Print("\n")
	}
	defer commitBuffer.Close()

	// TODO: meta data should be transferred on pull/push
	// Corral all the project meta data
	// Should return a JSON string of all the data in .silk/*
	// projectMeta := ""

	// Get each line from every file, creating new FileChange{}s
	files := ListAllFiles()
	for _, file := range files {
		bufferFile, err := os.Open(file)
		if err != nil {
			fmt.Println(cWarning("\n\tError") + ": unable to open file to add to initial commit buffer")
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

	// Creates the commit buffer data & writes to the file
	commitBufferData, _ := json.MarshalIndent(
		&RootCommitBuffer{
			ProjectName: SilkMetaFile().ProjectName,
			Changes:     allFileChanges,
		},
		"",
		"	",
	)

	_, err = commitBuffer.WriteString(string(commitBufferData) + "\n")
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to write initial commit buffer file data")
		fmt.Print("\n")
	}
}
