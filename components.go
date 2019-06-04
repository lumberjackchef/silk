package main

import (
  "encoding/json"
  // "fmt"
  "io/ioutil"
  "os"
)

type ComponentList struct {
  ProjectName   string    `json:"project_name"`
  ComponentList []string  `json:"component_list"`
}

type ComponentMeta struct {
  ProjectName   string  `json:"project_name"`
  ComponentName string  `json:"component_name"`
  InitDate      string  `json:"init_date"`
  Version       string  `json:"version"`
  Description   string  `json:"description"`
}

// Project metadata helper
func AddToSilkComponentList(componentName string) error {
  var componentFileData ComponentList

  // Open, check, & defer closing of the component list file
  // TODO: Add SilkRoot() here
  componentJsonFile, componentJsonFileErr := os.Open(".silk/components.json")
  check(componentJsonFileErr)
  defer componentJsonFile.Close()

  // Get the []byte version of the json data
  componentByteValue, componentByteValueErr := ioutil.ReadAll(componentJsonFile)
  check(componentByteValueErr)

  // Transform the []byte data into usable struct data
  componentJsonDataErr := json.Unmarshal(componentByteValue, &componentFileData)
  check(componentJsonDataErr)

  // Append the component name
  componentFileData.ComponentList = append(componentFileData.ComponentList, componentName)
  componentFileJsonData, componentFileJsonDataErr := json.MarshalIndent(componentFileData, "", "  ")
  check(componentFileJsonDataErr)

  // Write the component addition to the file
  // TODO: Add SilkRoot() here
  componentFileJsonDataWriteErr := ioutil.WriteFile(".silk/components.json", []byte(string(componentFileJsonData) + "\n"), 0766)
  check(componentFileJsonDataWriteErr)

  return nil
}
