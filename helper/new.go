package helper

import (
	"encoding/json"
	"os"
	"time"
)

// CreateRootMetaFile creates the root project meta file
func CreateRootMetaFile(projectName string) {
	// Creates the project meta json file
	projectMeta, projectMetaErr := os.Create(RootDirectoryName + "/meta.json")
	Check(projectMetaErr)
	defer projectMeta.Close()

	// Creates the project metadata & writes to the file
	dT := time.Now().String()
	projectMetaData, _ := json.MarshalIndent(
		&ProjectMeta{
			ProjectName: projectName,
			InitDate:    dT,
			Version:     "0.0.0",
		},
		"",
		"  ",
	)

	_, projectMetaWriteErr := projectMeta.WriteString(string(projectMetaData) + "\n")
	Check(projectMetaWriteErr)
}
