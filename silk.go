package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "silk"
	app.Usage = "A modern version control paradigm for service oriented architectures."

	app.Commands = []cli.Command{
		{
			Name:    "new",
			Aliases: []string{"n"},
			Usage:   "Create a new silk project",
			Action: func(c *cli.Context) error {
				// Project tracking folder. This checks if the folder exists, creates it if not.
				// TODO: Add SilkRoot() here
				if _, err := os.Stat(rootDirectoryName); os.IsNotExist(err) {
					if c.NArg() > 0 {
						// Creates the silk directory
						os.Mkdir(rootDirectoryName, 0766)

						// Creates the project meta json file
						projectMeta, projectMetaErr := os.Create(rootDirectoryName + "/meta.json")
						check(projectMetaErr)
						defer projectMeta.Close()

						// Creates the project metadata & writes to the file
						dT := time.Now().String()
						projectMetaData, _ := json.MarshalIndent(
							&ProjectMeta{
								ProjectName: fmt.Sprintf(c.Args().Get(0)),
								InitDate:    dT,
								Version:     "0.0.0",
							},
							"",
							"  ",
						)

						_, projectMetaWriteErr := projectMeta.WriteString(string(projectMetaData) + "\n")
						check(projectMetaWriteErr)

						// Create a blank components list file
						componentsList, componentsListErr := os.Create(rootDirectoryName + "/components.json")
						check(componentsListErr)
						defer componentsList.Close()

						// Creates the components data & writes to the file
						componentsListData, _ := json.MarshalIndent(
							&ComponentList{
								ProjectName:   fmt.Sprintf(c.Args().Get(0)),
								ComponentList: []string{},
							},
							"",
							"  ",
						)

						_, componentsListWriteError := componentsList.WriteString(string(componentsListData) + "\n")
						check(componentsListWriteError)

						// Confirmation message
						fmt.Println("New project " + fmt.Sprintf(c.Args().Get(0)) + " created!")
					} else {
						fmt.Println("No project name specified!")
					}
				} else {
					fmt.Println("Warning: this is an existing silk project!")
				}
				return nil
			},
		},
		{
			Name:    "status",
			Aliases: []string{"s"},
			Usage:   "Get the status of the current project and/or component.",
			Action: func(c *cli.Context) error {
				commandAction(
					func() {
						// Print status
						fmt.Println("Project: " + SilkMetaFile().ProjectName + "\n")

						// File list
						// TODO: needs to be updated in some way to indicate status
						// TODO: Add check for if within component or root
						// TODO: Add SilkRoot() here
						// TODO: Add ComponentRoot() here
						var files []string
						err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
							// TODO: Need to add some sort of .silkignore file here to exclude certain files/types && always ignore the .silk directory files
							if !info.IsDir() && !strings.HasPrefix(path, ".") {
								files = append(files, path)
							}
							return nil
						})
						check(err)

						for _, file := range files {
							fmt.Println("\t" + file)
						}

						// TODO: add list of returned files to some sort of hash function & return as a commit hash of some kind
					},
				)
				return nil
			},
		},
		{
			Name:  "clone",
			Usage: "Copies down the project root from remote and sets up all default branches & remotes.",
			Action: func(c *cli.Context) error {
				commandAction(func() { fmt.Println("Coming Soon!") })
				// fmt.Println(SilkRoot())
				return nil
			},
		},
		{
			Name:    "component",
			Aliases: []string{"c"},
			Usage:   "If no arguments, lists all components in the current project. If a name is supplied, this will either move to the component, clone from remote & move to the component, or it will create a new component of name [name]",
			Action: func(c *cli.Context) error {
				commandAction(func() {
					if c.NArg() > 0 {
						// Parameterized & lower-cased version of the user input string
						// TODO: add SilkRoot() here
						componentDirectory := fmt.Sprintf(strings.Join(strings.Split(strings.ToLower(c.Args().Get(0)), " "), "-"))
						var componentConfigDirectory = componentDirectory + "/.silk-component"

						// Component tracking directory. This checks if the directory exists, creates it if not.
						_, componentConfigErr := os.Stat(componentConfigDirectory)
						if os.IsNotExist(componentConfigErr) {
							// creates the '{component}/.silk-component directory as well as the {component} directory if one or both don't exist
							os.MkdirAll(componentConfigDirectory, 0766)

							// Creates the project meta json file
							componentMeta, componentMetaErr := os.Create(componentConfigDirectory + "/meta.json")
							check(componentMetaErr)
							defer componentMeta.Close()

							// Creates the project metadata & writes to the file
							dT := time.Now().String()
							componentMetaData, _ := json.MarshalIndent(&ComponentMeta{
								ProjectName:   SilkMetaFile().ProjectName,
								ComponentName: componentDirectory,
								InitDate:      dT,
								Version:       "0.0.0",
							}, "", "  ")
							_, componentMetaWriteErr := componentMeta.WriteString(string(componentMetaData) + "\n")
							check(componentMetaWriteErr)

							// Adds component to component list file
							AddToSilkComponentList(componentDirectory)

							// Confirmation message
							fmt.Println("New component " + componentDirectory + " created!")
						}
						// cd to component directory
						// NOTE: awaiting addition of SilkRoot() above
						os.Chdir(componentDirectory)

						// TODO: checkout this component's latest working branch branch, `master` for new components
					} else {
						// TODO: list an index of components
						fmt.Println(GetComponentIndex())
					}
				})
				return nil

				// TODO add method for removing a component
			},
		},
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "Lists or edits the current version of the project",
			Action: func(c *cli.Context) error {
				commandAction(
					func() {
						var metaData ProjectMeta

						// Open, check, & defer closing of the meta data file
						// TODO: Add SilkRoot() here
						metaFile, metaFileErr := os.Open(".silk/meta.json")
						check(metaFileErr)
						defer metaFile.Close()

						// Get the []byte version of the json data
						byteValue, byteValueErr := ioutil.ReadAll(metaFile)
						check(byteValueErr)

						// Transform the []byte data into usable struct data
						metaDataErr := json.Unmarshal(byteValue, &metaData)
						check(metaDataErr)

						if c.NArg() > 0 {
							// Change the version & transform back to actual json
							metaData.Version = fmt.Sprintf(c.Args().Get(0))
							metaDataJSON, metaDataJSONErr := json.MarshalIndent(metaData, "", " ")
							check(metaDataJSONErr)

							// Write the version change to the file
							// TODO: Add SilkRoot() here
							metaFileWriteErr := ioutil.WriteFile(".silk/meta.json", []byte(string(metaDataJSON)+"\n"), 0766)
							check(metaFileWriteErr)

							// Confirmation message
							fmt.Println("Version successfull updated to " + metaData.Version + "!")
						} else {
							// If the user just wants to check the version and not change it
							fmt.Println(metaData.Version)
						}
					},
				)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// Checks if this is a silk project before running a command
func commandAction(f func()) string {
	// TODO: Add SilkRoot() here
	if _, err := os.Stat(rootDirectoryName); os.IsNotExist(err) {
		fmt.Println("Warning: this is not a silk project! To create a new silk project, run `$ silk new`")
	} else {
		f()
	}
	return ""
}
