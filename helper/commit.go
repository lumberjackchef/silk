package helper

import (
	"encoding/json"
	"os"
)

// CreateInitialCommitBuffer initiates the commit buffer and writes a new file
func CreateInitialCommitBuffer() {
	// Creates the commit buffer file
	commitBuffer, commitBufferErr := os.Create(RootDirectoryName + "/commit/buffer")
	Check(commitBufferErr)
	defer commitBuffer.Close()

	// TODO: meta data should be transferred on pull/push
	// Corral all the project meta data
	// Should return a JSON string of all the data in .silk/*
	// projectMeta := ""

	// Add all project files to initial buffer
	changes := []FileChange{}

	// Creates the commit buffer data & writes to the file
	commitBufferData, _ := json.MarshalIndent(
		&RootCommitBuffer{
			ProjectName: SilkMetaFile().ProjectName,
			Changes:     changes,
		},
		"",
		"  ",
	)

	_, commitBufferWriteErr := commitBuffer.WriteString(string(commitBufferData) + "\n")
	Check(commitBufferWriteErr)
}
