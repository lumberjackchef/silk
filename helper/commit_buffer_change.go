package helper

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
)

// AddFileToCommitBuffer ...
// TODO: update to work for file or directory
func AddFileToCommitBuffer(fileName string) {
	var commitBuffer RootCommitBuffer
	var changes []FileChange
	oldChanges := CommitBuffer().Changes

	bufferFile, err := os.Open(SilkRoot() + "/" + RootDirectoryName + "/commit/buffer")
	Check(err)
	defer bufferFile.Close()

	byteValue, err := ioutil.ReadAll(bufferFile)
	Check(err)

	err = json.Unmarshal(byteValue, &commitBuffer)
	Check(err)

	changedFile, err := os.Open(SilkRoot() + "/" + fileName)
	Check(err)
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
	Check(err)

	err = ioutil.WriteFile(SilkRoot()+"/.silk/commit/buffer", []byte(string(commitBufferJSON)+"\n"), 0766)
	Check(err)
}
