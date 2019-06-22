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
