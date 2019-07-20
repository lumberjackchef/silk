/*
Package helper provides global helper funcs for Silk core and all related packages
*/
package helper

import (
	"fmt"

	"github.com/fatih/color"
)

// TODO: Move all error handling to an errors.go file/package

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

// UniqueNonEmptyElementsOf ...
func UniqueNonEmptyElementsOf(s []string) []string {
	unique := make(map[string]bool, len(s))
	us := make([]string, len(unique))
	for _, elem := range s {
		if len(elem) != 0 {
			if !unique[elem] {
				us = append(us, elem)
				unique[elem] = true
			}
		}
	}

	return us
}
