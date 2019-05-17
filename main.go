package main

import (
  "encoding/json"
  "fmt"
  "log"
  "os"

  "github.com/urfave/cli"
)

type project_meta struct {
  ProjectName  string  `json:"project_name"`
}

func main() {
  app := cli.NewApp()
  app.Name = "silk"
  app.Usage = "A modern version control paradigm for service oriented architectures."
  app.Version = "0.0.0"

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

            f_config, err := os.Create(".silk/config")
            check(err)
            defer f_config.Close()

            project_meta_data, _ := json.Marshal(&project_meta{
              ProjectName: fmt.Sprintf(c.Args().Get(0)),
            })
            _, err2 := f_config.WriteString(string(project_meta_data))
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
