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

	return cli.Command{
		Name:  "commit",
		Usage: "Tags the current commit buffer and resets all file statuses",
		Action: func(c *cli.Context) error {
			commitFile, err := os.Create(helper.RootDirectoryName + "/commit/latest")
			helper.Check(err)
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

			_, err = commitFile.WriteString(string(commitFileData) + "\n")
			helper.Check(err)

			hash := sha1.New()
			if _, err := io.Copy(hash, commitFile); err != nil {
				log.Fatal(err)
			}
			sum := hash.Sum(nil)
			str := hex.EncodeToString(sum)

			hashedCommitFile, err := os.Create(helper.RootDirectoryName + "/commit/" + str)
			helper.Check(err)
			defer hashedCommitFile.Close()

			// TODO: compress these files (brotli?)
			_, err = hashedCommitFile.Write(commitFileData)
			helper.Check(err)

			// TODO: remove committed files from buffer somehow
			fmt.Println("\n\t" + cNotice("Success!"))
			fmt.Print("\n")
			// TODO: find a way to identify merge commits clearly, AWAITING branching commands creation
			return nil
		},
	}
}
