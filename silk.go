package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "path/filepath"
  "os"
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
      Name: "new",
      Aliases: []string{"n"},
      Usage: "Create a new silk project",
      Action: func(c *cli.Context) error {
        // Project tracking folder. This checks if the folder exists, creates it if not.
        path := ".silk"
        if _, err := os.Stat(path); os.IsNotExist(err) {
          if c.NArg() > 0 {
            // Creates the silk directory
            os.Mkdir(path, 0766)

            // Creates the project meta json file
            fMeta, err := os.Create(".silk/meta.json")
            check(err)
            defer fMeta.Close()

            // Creates the project metadata & writes to the file
            dT := time.Now().String()
            projectMetaData, _ := json.MarshalIndent(&ProjectMeta{
              ProjectName:  fmt.Sprintf(c.Args().Get(0)),
              InitDate:     dT,
              Version:      "0.0.0",
            }, "", "  ")
            _, err2 := fMeta.WriteString(string(projectMetaData) + "\n")
            check(err2)

            // Confirmation message
            fmt.Println("New project created!")
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
      Name: "status",
      Aliases: []string{"s"},
      Usage: "Get the status of the current project and/or component.",
      Action: func(c *cli.Context) error {
        commandAction(
          func() {
            // Print status
            fmt.Println("Project: " + metaFile().ProjectName + "\n")

            // File list, needs to be updated in some way to indicate status
            var files []string
            err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
              // Need to add some sort of .silkignore file here to exclude certain files/types && always ignore the .silk directory files
              if !info.IsDir() && !strings.HasPrefix(path, ".") {
                files = append(files, path)
              }
              return nil
            })
            check(err)

            for _, file := range files{
              fmt.Println("\t" + file)
            }
          },
        )
        return nil
      },
    },
    {
      Name: "clone",
      Usage: "Copies down the project root from remote and sets up all default branches & remotes.",
      Action: func(c *cli.Context) error {
        commandAction(func() { fmt.Println("Coming Soon!") })
        return nil
      },
    },
    {
      Name: "component",
      Aliases: []string{"c"},
      Usage: "If no arguments, lists all components in the current project. If a name is supplied, this will either move to the component, clone from remote & move to the component, or it will create a new component of name [name]",
      Action: func(c *cli.Context) error {
        commandAction(func() { fmt.Println("Coming Soon!") })
        return nil
      },
    },
    {
      Name: "version",
      Aliases: []string{"v"},
      Usage: "Lists or edits the current version of the project",
      Action: func(c *cli.Context) error {
        commandAction(
          func() {
            var metaData ProjectMeta

            // Open, check, & defer closing of the meta data file
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
              metaDataJson, metaDataJsonErr := json.MarshalIndent(metaData, "", " ")
              check(metaDataJsonErr)

              // Write the version change to the file
              metaFileWriteErr := ioutil.WriteFile(".silk/meta.json", []byte(string(metaDataJson) + "\n"), 0766)
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
