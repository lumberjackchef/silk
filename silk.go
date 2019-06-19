package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/lumberjackchef/silk/cmd"
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
				if _, err := os.Stat(cmd.RootDirectoryName); os.IsNotExist(err) {
					if c.NArg() > 0 {
						// Creates the silk directory
						os.Mkdir(cmd.RootDirectoryName, 0766)

						// Creates the project meta json file
						projectMeta, projectMetaErr := os.Create(cmd.RootDirectoryName + "/meta.json")
						cmd.Check(projectMetaErr)
						defer projectMeta.Close()

						// Creates the project metadata & writes to the file
						dT := time.Now().String()
						projectMetaData, _ := json.MarshalIndent(
							&cmd.ProjectMeta{
								ProjectName: fmt.Sprintf(c.Args().Get(0)),
								InitDate:    dT,
								Version:     "0.0.0",
							},
							"",
							"  ",
						)

						_, projectMetaWriteErr := projectMeta.WriteString(string(projectMetaData) + "\n")
						cmd.Check(projectMetaWriteErr)

						// Create a blank components list file
						componentsList, componentsListErr := os.Create(cmd.RootDirectoryName + "/components.json")
						cmd.Check(componentsListErr)
						defer componentsList.Close()

						// Creates the components data & writes to the file
						componentsListData, _ := json.MarshalIndent(
							&cmd.ComponentList{
								ProjectName:   fmt.Sprintf(c.Args().Get(0)),
								ComponentList: []string{},
							},
							"",
							"  ",
						)

						_, componentsListWriteError := componentsList.WriteString(string(componentsListData) + "\n")
						cmd.Check(componentsListWriteError)

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
				cmd.CommandAction(
					func() {
						// Print status
						fmt.Printf("\t%s "+cmd.SilkMetaFile().ProjectName+"\n\n", cNotice("Project:"))

						if cmd.IsComponentOrRoot() == "component" {
							os.Chdir(cmd.SilkComponentRoot())
						} else {
							os.Chdir(cmd.SilkRoot())
						}

						currentWorkingDirectory, _ := os.Getwd()

						// File list
						files := cmd.ComposeFileList(currentWorkingDirectory)

						// Print the file status
						cmd.ListFilesInCommitBuffer(files)
					},
				)
				return nil
			},
		},
		{
			Name:  "clone",
			Usage: "Copies down the project root from remote and sets up all default branches & remotes.",
			Action: func(c *cli.Context) error {
				cmd.CommandAction(func() { fmt.Printf("\t%s\n", cNotice("Coming Soon!")) })
				return nil
			},
		},
		{
			Name:    "component",
			Aliases: []string{"c"},
			Usage:   "If no arguments, lists all components in the current project. If a name is supplied, this will either move to the component, clone from remote & move to the component, or it will create a new component of name [name]",
			Action: func(c *cli.Context) error {
				cmd.CommandAction(func() {
					if c.NArg() > 0 {
						// Parameterized & lower-cased version of the user input string
						componentName := fmt.Sprintf(strings.Join(strings.Split(strings.ToLower(c.Args().Get(0)), " "), "-"))
						componentConfigDirectory := cmd.SilkRoot() + "/" + componentName + "/.silk-component"

						cmd.CreateComponentsListFile(componentName, componentConfigDirectory)
					} else {
						// Lists index of components
						if len(cmd.GetComponentIndex()) > 0 {
							fmt.Println(cNotice("\tComponents in project " + cmd.SilkMetaFile().ProjectName + ":"))
							for _, component := range cmd.GetComponentIndex() {
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
					Name:   "remove",
					Usage:  "remove an existing component",
					Action: cmd.RemoveComponent,
				},
			},
		},
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "Lists or edits the current version of the project",
			Action: func(c *cli.Context) error {
				cmd.CommandAction(
					func() {
						cmd.PrintOrChangeVersion(c)
					},
				)
				return nil
			},
		},
		{
			Name:  "add",
			Usage: "Adds a file or files to the current commit buffer",
			Action: func(c *cli.Context) error {
				// TODO: Add a list of changes to files to a commit buffer
				// Check if the root commit exists
				// diff whole files first to eliminate unchanged files
				// diff changed files line by line
				// add changes to a commit buffer file
				// should be file name, line number, & actual code changes
				cmd.CommandAction(func() { fmt.Printf("\t%s\n", cNotice("Coming Soon!")) })
				return nil
			},
		},
		{
			Name:  "commit",
			Usage: "Tags the current commit buffer and resets all file statuses",
			Action: func(c *cli.Context) error {
				// TODO: add commit buffer to some sort of hash function
				// TODO: return as a commit hash of some kind when committing
				// TODO: add commit message ability
				cmd.CommandAction(func() { fmt.Printf("\t%s\n", cNotice("Coming Soon!")) })
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
