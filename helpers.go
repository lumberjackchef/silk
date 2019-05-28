package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "os"
)

// Checks if this is a silk project before running a command
func commandAction(f func()) string {
  path := ".silk"
  if _, err := os.Stat(path); os.IsNotExist(err) {
    fmt.Println("Warning: this is not a silk project! To create a new silk project, run `$ silk new`")
  } else {
    f()
  }
  return ""
}

// Error checking & logging
func check(e error) {
  if e != nil {
    panic(e)
  }
}

// Project metadata helper
func metaFile() ProjectMeta {
  var fileData ProjectMeta

  // Open, check, & defer closing of the meta data file
  jsonFile, jsonFileErr := os.Open(".silk/meta.json")
  check(jsonFileErr)
  defer jsonFile.Close()

  // Get the []byte version of the json data
  byteValue, byteValueErr := ioutil.ReadAll(jsonFile)
  check(byteValueErr)

  // Transform the []byte data into usable struct data
  jsonDataErr := json.Unmarshal(byteValue, &fileData)
  check(jsonDataErr)

  return fileData
}
