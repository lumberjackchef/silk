package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "time"

  "github.com/urfave/cli"
)

type ProjectMeta struct {
  ProjectName string  `json:"project_name"`
  InitDate    string  `json:"init_date"`
  Version     string  `json:"version"`
}

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
            os.Mkdir(path, 0766)

            f_meta, err := os.Create(".silk/meta.json")
            check(err)
            defer f_meta.Close()

            dt := time.Now().String()

            project_meta_data, _ := json.MarshalIndent(&ProjectMeta{
              ProjectName:  fmt.Sprintf(c.Args().Get(0)),
              InitDate:     dt,
              Version:      "0.0.0",
            }, "", "  ")
            _, err2 := f_meta.WriteString(string(project_meta_data) + "\n")
            check(err2)

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
        commandAction(func() { fmt.Println("Coming Soon!") })
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
            var meta_data ProjectMeta

            // Open, check, & defer closing of the meta data file
            meta_file, meta_file_err := os.Open(".silk/meta.json")
            check(meta_file_err)
            defer meta_file.Close()

            // Get the []byte version of the json data
            byte_value, byte_value_err := ioutil.ReadAll(meta_file)
            check(byte_value_err)

            // Transform the []byte data into usable struct data
            meta_data_err := json.Unmarshal(byte_value, &meta_data)
            check(meta_data_err)

            if c.NArg() > 0 {
              // Change the version & transform back to actual json
              meta_data.Version = fmt.Sprintf(c.Args().Get(0))
              meta_data_json, meta_data_json_err := json.MarshalIndent(meta_data, "", " ")
              check(meta_data_json_err)

              // Write the version change to the file
              meta_file_write_err := ioutil.WriteFile(".silk/meta.json", []byte(string(meta_data_json) + "\n"), 0766)
              check(meta_file_write_err)

              // Confirmation message
              fmt.Println("Version successfull updated to " + meta_data.Version + "!")
            } else {
              // If the user just wants to check the version and not change it
              fmt.Println(meta_data.Version)
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

var commandAction = func(f func()) string {
  path := ".silk"
  if _, err := os.Stat(path); os.IsNotExist(err) {
    fmt.Println("Warning: this is not a silk project! To create a new silk project, run `$ silk new`")
  } else {
    f()
  }
  return ""
}

var check = func(e error) {
  if e != nil {
    panic(e)
  }
}
