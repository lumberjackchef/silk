/*
Package helper provides global helper funcs for Silk core and all related packages
*/
package helper

import (
	"fmt"

	"github.com/fatih/color"
)

// Check provides basic error checking & logging
// TODO: implement better logging/error handling. Panic is not the only way to handle an error
//       need to implement recovers as well
// TODO: Move all error handling to an errors.go file/package?
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// CommandAction checks if this is a silk project before running a command
func CommandAction(f func()) string {
	cWarning := color.New(color.FgYellow).SprintFunc()

	if IsComponentOrRoot() == "false" {
		fmt.Printf("\t%s this is not a silk project! To create a new silk project, run `$ silk new`\n", cWarning("Warning:"))
	} else {
		f()
	}
	return ""
}
