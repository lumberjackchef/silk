package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func main() {
	// Colors setup
	cNotice := color.New(color.FgGreen).SprintFunc()
	cWarning := color.New(color.FgYellow).SprintFunc()

	// Application setup
	app := cli.NewApp()
	app.Name = "silk"
	app.Usage = "A modern version control paradigm for service oriented architectures."

	// Allows for bash completion of commands & subcommands
	app.EnableBashCompletion = true

	app.Commands = []cli.Command{
		{
			Name:    "new",
			Aliases: []string{"n"},
			Usage:   "Create a new silk project",
			Action: func(c *cli.Context) error {
				// Project tracking folder. This checks if the folder exists, creates it if not.
				// TODO: implement checking similar to SilkRoot() to ensure we're not within the bounds of another project
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
						fmt.Println("\tNew project " + cNotice(fmt.Sprintf(c.Args().Get(0))) + " created!")
					} else {
						fmt.Printf("\t%s No project name specified!\n", cWarning("Warning:"))
					}
				} else {
					fmt.Printf("\t%s this is an existing silk project!\n", cWarning("Warning:"))
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
						fmt.Printf("\t%s "+SilkMetaFile().ProjectName+"\n\n", cNotice("Project:"))

						if IsComponentOrRoot() == "component" {
							os.Chdir(SilkComponentRoot())
						} else {
							os.Chdir(SilkRoot())
						}

						currentWorkingDirectory, _ := os.Getwd()

						// File list
						files := ComposeFileList(currentWorkingDirectory)

						// Print the file status
						ListFilesInCommitBuffer(files)
					},
				)
				return nil
			},
		},
		{
			Name:  "clone",
			Usage: "Copies down the project root from remote and sets up all default branches & remotes.",
			Action: func(c *cli.Context) error {
				commandAction(func() { fmt.Printf("\t%s\n", cNotice("Coming Soon!")) })
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
						componentName := fmt.Sprintf(strings.Join(strings.Split(strings.ToLower(c.Args().Get(0)), " "), "-"))
						componentDirectory := SilkRoot() + "/" + componentName
						componentConfigDirectory := componentDirectory + "/.silk-component"

						CreateComponentsListFile(componentName, componentConfigDirectory)
					} else {
						// Lists index of components
						if len(GetComponentIndex()) > 0 {
							fmt.Println(cNotice("\tComponents in project " + SilkMetaFile().ProjectName + ":"))
							for _, component := range GetComponentIndex() {
								fmt.Println("\t\t" + component)
							}
						} else {
							fmt.Printf("\t%s There are no components in the current project.\n", cWarning("Warning:"))
						}
					}
				})
				return nil
			},
			Subcommands: []cli.Command{
				{
					Name:  "remove",
					Usage: "remove an existing component",
					Action: func(c *cli.Context) error {
						commandAction(func() {
							if c.NArg() > 0 {
								// Parameterized & lower-cased version of the user input string
								componentName := fmt.Sprintf(strings.Join(strings.Split(strings.ToLower(c.Args().Get(0)), " "), "-"))
								RemoveComponent(componentName)
							} else {
								fmt.Printf("\t%s No component name specified.\n", cWarning("Error:"))
							}
						})
						return nil
					},
				},
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
						metaFile, metaFileErr := os.Open(SilkRoot() + "/.silk/meta.json")
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
							metaFileWriteErr := ioutil.WriteFile(SilkRoot()+"/.silk/meta.json", []byte(string(metaDataJSON)+"\n"), 0766)
							check(metaFileWriteErr)

							// Confirmation message
							fmt.Println("\tVersion successfully updated to " + cNotice(metaData.Version) + "!")
						} else {
							// If the user just wants to check the version and not change it
							fmt.Println("\t" + cNotice(metaData.Version))
						}
					},
				)
				return nil
			},
		},
		{
			Name:  "add",
			Usage: "Adds a file or files to the current commit buffer",
			Action: func(c *cli.Context) error {
				// TODO: Add a list of files to a commit buffer
				commandAction(func() { fmt.Printf("\t%s\n", cNotice("Coming Soon!")) })
				return nil
			},
		},
		{
			Name:  "commit",
			Usage: "Tags the current commit buffer and resets all file statuses",
			Action: func(c *cli.Context) error {

				// TODO: add list of files in commit buffer to some sort of hash function
				// TODO: return as a commit hash of some kind when committing
				// TODO: add commit message ability
				commandAction(func() { fmt.Printf("\t%s\n", cNotice("Coming Soon!")) })
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
