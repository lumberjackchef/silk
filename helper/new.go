package helper

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

// CreateRootMetaFile creates the root project meta file
func CreateRootMetaFile(projectName string) {
	cWarning := color.New(color.FgYellow).SprintFunc()

	// Creates the project meta json file
	projectMeta, err := os.Create(RootDirectoryName + "/meta.json")
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to create root meta file")
		fmt.Print("\n")
	}
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

	_, err = projectMeta.WriteString(string(projectMetaData) + "\n")
	if err != nil {
		fmt.Println(cWarning("\n\tError") + ": unable to write project meta")
		fmt.Print("\n")
	}
}
