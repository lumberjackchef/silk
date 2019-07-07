package helper

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// ProjectMeta ...
type ProjectMeta struct {
	ProjectName string `json:"project_name"`
	InitDate    string `json:"init_date"`
	Version     string `json:"version"`
	Description string `json:"description"`
	ProjectURL  string `json:"url"`
}

// RootCommitBuffer ...
type RootCommitBuffer struct {
	ProjectName string       `json:"project_name"`
	Changes     []FileChange `json:"changes"`
}

// FileChange structure for a single file change
type FileChange struct {
	FileName   string `json:"file_name"`
	LineNumber int    `json:"line_number"`
	Text       string `json:"text"`
}

// ComponentList .silk/components.json file structure
type ComponentList struct {
	ProjectName   string   `json:"project_name"`
	ComponentList []string `json:"component_list"`
}

// ComponentMeta .silk/../{component}/.silk-component/meta.json file structure
type ComponentMeta struct {
	ProjectName   string `json:"project_name"`
	ComponentName string `json:"component_name"`
	InitDate      string `json:"init_date"`
	Version       string `json:"version"`
	Description   string `json:"description"`
}

// SilkMetaFile provides project metadata in an easy to consume format
func SilkMetaFile() ProjectMeta {
	var fileData ProjectMeta

	// Open, check, & defer closing of the meta data file
	jsonFile, jsonFileErr := os.Open(SilkRoot() + "/.silk/meta.json")
	Check(jsonFileErr)
	defer jsonFile.Close()

	// Get the []byte version of the json data
	byteValue, byteValueErr := ioutil.ReadAll(jsonFile)
	Check(byteValueErr)

	// Transform the []byte data into usable struct data
	jsonDataErr := json.Unmarshal(byteValue, &fileData)
	Check(jsonDataErr)

	return fileData
}

// This is currently not in use

// ComponentMetaFile provides component metadata in an easy to consume format
func ComponentMetaFile() ComponentMeta {
	var fileData ComponentMeta

	// Open, check, & defer closing of the meta data file
	jsonFile, jsonFileErr := os.Open(SilkComponentRoot() + "/.silk-component/meta.json")
	Check(jsonFileErr)
	defer jsonFile.Close()

	// Get the []byte version of the json data
	byteValue, byteValueErr := ioutil.ReadAll(jsonFile)
	Check(byteValueErr)

	// Transform the []byte data into usable struct data
	jsonDataErr := json.Unmarshal(byteValue, &fileData)
	Check(jsonDataErr)

	return fileData
}

// CommitBuffer provides commit buffer data in an easy to consume format
func CommitBuffer() RootCommitBuffer {
	// TODO: update to be useable for root or a component
	var bufferData RootCommitBuffer

	// Open, check, & defer closing of the meta data file
	jsonFile, jsonFileErr := os.Open(SilkRoot() + "/" + RootDirectoryName + "/commit/buffer")
	Check(jsonFileErr)
	defer jsonFile.Close()

	// Get the []byte version of the json data
	byteValue, byteValueErr := ioutil.ReadAll(jsonFile)
	Check(byteValueErr)

	// Transform the []byte data into usable struct data
	jsonDataErr := json.Unmarshal(byteValue, &bufferData)
	Check(jsonDataErr)

	return bufferData
}

// LatestCommit provides latest commit data in an easy to consume format
func LatestCommit() RootCommitBuffer {
	latestCommit := RootDirectoryName + "/commit/latest"
	if _, err := os.Stat(latestCommit); !os.IsNotExist(err) {
		var bufferData RootCommitBuffer

		// Open, check, & defer closing of the meta data file
		jsonFile, jsonFileErr := os.Open(SilkRoot() + "/" + latestCommit)
		Check(jsonFileErr)
		defer jsonFile.Close()

		// Get the []byte version of the json data
		byteValue, byteValueErr := ioutil.ReadAll(jsonFile)
		Check(byteValueErr)

		// Transform the []byte data into usable struct data
		jsonDataErr := json.Unmarshal(byteValue, &bufferData)
		Check(jsonDataErr)

		return bufferData
	}
	return RootCommitBuffer{}
}
