package helper

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
)

// AddFileToCommitBuffer ...
// TODO: update to work for file or directory
func AddFileToCommitBuffer(fileName string) {
	var commitBuffer RootCommitBuffer
	var changes []FileChange
	cWarning := color.New(color.FgYellow).SprintFunc()
	oldChanges := CommitBuffer().Changes

	bufferFile, err := os.Open(SilkRoot() + "/" + RootDirectoryName + "/commit/buffer")
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to open buffer file")
		fmt.Print("\n")
	}
	defer bufferFile.Close()

	byteValue, err := ioutil.ReadAll(bufferFile)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to read buffer file byte value")
		fmt.Print("\n")
	}

	err = json.Unmarshal(byteValue, &commitBuffer)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to unmarshal commit buffer")
		fmt.Print("\n")
	}

	changedFile, err := os.Open(SilkRoot() + "/" + fileName)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to open file")
		fmt.Print("\n")
	}
	defer changedFile.Close()

	scanner := bufio.NewScanner(changedFile)
	line := 0
	for scanner.Scan() {
		line = line + 1

		fileChange := FileChange{
			FileName:   fileName,
			LineNumber: line,
			Text:       scanner.Text(),
		}

		changes = append(changes, fileChange)
	}
	// TODO: remove FileChange{}s already committed
	changes = append(changes, oldChanges...)
	commitBuffer.Changes = changes

	commitBufferJSON, err := json.MarshalIndent(commitBuffer, " ", "")
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to marshal json to commit buffer")
		fmt.Print("\n")
	}

	err = ioutil.WriteFile(SilkRoot()+"/.silk/commit/buffer", []byte(string(commitBufferJSON)+"\n"), 0766)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to write to commit buffer file")
		fmt.Print("\n")
	}
}

// ClearCommitBuffer removes all entries from the commit buffer
func ClearCommitBuffer() {
	var commitBuffer RootCommitBuffer
	cWarning := color.New(color.FgYellow).SprintFunc()

	bufferFile, err := os.Open(SilkRoot() + "/" + RootDirectoryName + "/commit/buffer")
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to open commit buffer file")
		fmt.Print("\n")
	}
	defer bufferFile.Close()

	// Remove all bufferFile.Changes()
	byteValue, err := ioutil.ReadAll(bufferFile)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to read buffer file byte values")
		fmt.Print("\n")
	}

	err = json.Unmarshal(byteValue, &commitBuffer)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to create latest commit file")
		fmt.Print("\n")
	}

	commitBuffer.Changes = []FileChange{}

	commitBufferJSON, err := json.MarshalIndent(commitBuffer, " ", "")
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to marshal json to commit buffer")
		fmt.Print("\n")
	}

	err = ioutil.WriteFile(SilkRoot()+"/.silk/commit/buffer", []byte(string(commitBufferJSON)+"\n"), 0766)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to write to commit buffer file")
		fmt.Print("\n")
	}
}
