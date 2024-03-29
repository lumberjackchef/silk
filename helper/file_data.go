package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
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
	cWarning := color.New(color.FgYellow).SprintFunc()

	// Open, check, & defer closing of the meta data file
	jsonFile, err := os.Open(SilkRoot() + "/.silk/meta.json")
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to open project meta file")
		fmt.Print("\n")
	}
	defer jsonFile.Close()

	// Get the []byte version of the json data
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to read byte value of project meta file")
		fmt.Print("\n")
	}

	// Transform the []byte data into usable struct data
	err = json.Unmarshal(byteValue, &fileData)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to unmarshal project meta byte value")
		fmt.Print("\n")
	}

	return fileData
}

// This is currently not in use

// ComponentMetaFile provides component metadata in an easy to consume format
func ComponentMetaFile() ComponentMeta {
	var fileData ComponentMeta
	cWarning := color.New(color.FgYellow).SprintFunc()

	// Open, check, & defer closing of the meta data file
	jsonFile, err := os.Open(SilkComponentRoot() + "/.silk-component/meta.json")
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to open component meta file")
		fmt.Print("\n")
	}
	defer jsonFile.Close()

	// Get the []byte version of the json data
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to read byte value of component meta file")
		fmt.Print("\n")
	}

	// Transform the []byte data into usable struct data
	err = json.Unmarshal(byteValue, &fileData)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to unmarshal component meta byte value")
		fmt.Print("\n")
	}

	return fileData
}

// CommitBuffer provides commit buffer data in an easy to consume format
func CommitBuffer() RootCommitBuffer {
	// TODO: update to be useable for root or a component
	var bufferData RootCommitBuffer
	cWarning := color.New(color.FgYellow).SprintFunc()

	// Open, check, & defer closing of the meta data file
	jsonFile, err := os.Open(SilkRoot() + "/" + RootDirectoryName + "/commit/buffer")
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to open project buffer file")
		fmt.Print("\n")
	}
	defer jsonFile.Close()

	// Get the []byte version of the json data
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to read byte value of project buffer file")
		fmt.Print("\n")
	}

	// Transform the []byte data into usable struct data
	err = json.Unmarshal(byteValue, &bufferData)
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to unmarshal project buffer byte value")
		fmt.Print("\n")
	}

	return bufferData
}

// LatestCommit provides latest commit data in an easy to consume format
func LatestCommit() RootCommitBuffer {
	var bufferData RootCommitBuffer
	cWarning := color.New(color.FgYellow).SprintFunc()
	latestCommit := RootDirectoryName + "/commit/latest"

	if _, err := os.Stat(latestCommit); !os.IsNotExist(err) {
		// Open, check, & defer closing of the meta data file
		jsonFile, err := os.Open(SilkRoot() + "/" + latestCommit)
		if err != nil {
			fmt.Println(cWarning("\n\tError") + ": unable to open latest project commit file")
			fmt.Print("\n")
		}
		defer jsonFile.Close()

		// Get the []byte version of the json data
		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			fmt.Println(cWarning("\n\tError") + ": unable to read byte value of latest project commit file")
			fmt.Print("\n")
		}

		// Transform the []byte data into usable struct data
		err = json.Unmarshal(byteValue, &bufferData)
		if err != nil {
			fmt.Println(cWarning("\n\tError") + ": unable to unmarshal latest project commit byte value")
			fmt.Print("\n")
		}
	}

	return bufferData
}
