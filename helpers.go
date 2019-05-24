package silk

import (
  "fmt"
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
