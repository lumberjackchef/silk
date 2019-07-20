package cmd

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/lumberjackchef/silk/helper"
	"github.com/urfave/cli"
)

// CommitFile gives a structure to a commit file
type CommitFile struct {
	Project string              `json:"project_name"`
	User    string              `json:"user"`
	Date    string              `json:"date"`
	Message string              `json:"message"`
	Changes []helper.FileChange `json:"changes"`
	// Tree		[]helper.Tree				`json:"tree"` // TODO: need to create
}

// Commit adds commit buffer data to a commit file
func Commit() cli.Command {
	cNotice := color.New(color.FgGreen).SprintFunc()
	cWarning := color.New(color.FgYellow).SprintFunc()

	return cli.Command{
		Name:  "commit",
		Usage: "Tags the current commit buffer and resets all file statuses",
		Action: func(c *cli.Context) error {
			var commitFile *os.File
			latestCommittFile := helper.SilkRoot() + "/" + helper.RootDirectoryName + "/commit/latest"

			// TODO: check if there are changes in the commit buffer

			if _, err := os.Stat(latestCommittFile); os.IsNotExist(err) {
				commitFile, err = os.Create(latestCommittFile)
				if err != nil {
					fmt.Println(cWarning("\n\tError") + ": unable to create latest commit file")
					fmt.Print("\n")
				}
			} else {
				commitFile, err = os.OpenFile(latestCommittFile, os.O_APPEND|os.O_RDWR, os.ModeAppend)
				if err != nil {
					fmt.Println(cWarning("\n\tError") + ": unable to open latest commit file")
					fmt.Print("\n")
				}
			}

			defer commitFile.Close()

			commitFileData, _ := json.MarshalIndent(
				&CommitFile{
					Project: helper.SilkMetaFile().ProjectName,
					// TODO: add user creation, deletion, & management capabilities
					// TODO: update to call current user
					User: "root",
					Date: time.Now().String(),
					// TODO: add commit message ability
					Message: "latest commit",
					Changes: helper.CommitBuffer().Changes,
				},
				"",
				"	",
			)

			_, err := commitFile.WriteString(string(commitFileData) + "\n")
			if err != nil {
				fmt.Println(cWarning("\n\tError") + ": unable to write to commit file")
				fmt.Print("\n")
			}

			hash := sha1.New()
			if _, err := io.Copy(hash, commitFile); err != nil {
				log.Fatal(err)
				if err != nil {
					fmt.Println(cWarning("\n\tError") + ": unable to copy commit file contents to hash")
					fmt.Print("\n")
				}
			}
			sum := hash.Sum(nil)
			str := hex.EncodeToString(sum)

			hashedCommitFile, err := os.Create(helper.RootDirectoryName + "/commit/" + str)
			if err != nil {
				fmt.Println(cWarning("\n\tError") + ": unable to create hashed commit file")
				fmt.Print("\n")
			}
			defer hashedCommitFile.Close()

			// TODO: compress these files (brotli?)
			_, err = hashedCommitFile.Write(commitFileData)
			if err != nil {
				fmt.Println(cWarning("\n\tError") + ": unable to write to hashed commit file")
				fmt.Print("\n")
			}

			// TODO: remove committed files from buffer somehow
			helper.ClearCommitBuffer()

			fmt.Println("\n\t" + cNotice("Success!"))
			fmt.Print("\n")
			// TODO: find a way to identify merge commits clearly, AWAITING branching commands creation
			return nil
		},
	}
}
