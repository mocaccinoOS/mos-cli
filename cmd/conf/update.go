/*
Copyright © 2021 Ettore Di Giacinto <mudler@sabayon.org>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package conf

import (
	"fmt"

	config "github.com/MocaccinoOS/mos-cli/pkg/configfile"
	"github.com/MocaccinoOS/mos-cli/pkg/utils"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println("ERROR:", err)
	}
}
func evaluateDiff(f string, changeset config.ConfigChange, res config.Configs, interactive, removeall bool) {
	diffs, err := changeset.Diff(f)
	if err != nil {
		checkErr(err)
		return
	}

	drop := func() {
		if removeall {
			checkErr(res.CleanChanges(f))
		} else {
			checkErr(changeset.Remove())
		}
	}

	accept := func() {
		mergeRes, err := changeset.Merge(f)
		checkErr(err)
		checkErr(mergeRes.Apply())
	}

	if len(diffs) > 0 {
		if interactive {
			fmt.Print("\033[H\033[2J")
			dmp := diffmatchpatch.New()
			fmt.Printf("Diff for file: %s (changeset %s)\n", f, changeset.Path)

			fmt.Println("-----------------------------------------------------")

			fmt.Println(dmp.DiffPrettyText(diffs))

			fmt.Println("-----------------------------------------------------")

			r := utils.Accept("Do you want to accept the following changes")
			switch r {
			case utils.AcceptQuestion:
				accept()
				drop()
			case utils.DiscardQuestion:
				drop()
			}
		} else {
			fmt.Printf("Merging configuration for file: %s (changeset %s)\n", f, changeset.Path)
			accept()
			drop()
		}
	} else {
		drop()
	}
}

func NewUpdateCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "update",
		Short: "Updates configuration file in the system",
		Long: `update allows you to merge changes in configuration files after upgrades. Offers an interactive 
mode which guides the user to apply the changes.

To merge all configs:
$ mos update --interactive=false

Interactively:

$ mos update

Checking all diffs 1-by-1:

$ mos update --all
`,
		Run: func(cmd *cobra.Command, args []string) {
			path, _ := cmd.Flags().GetString("path")
			interactive, _ := cmd.Flags().GetBool("interactive")
			all, _ := cmd.Flags().GetBool("all")

			res := config.Scan(path)

			fmt.Printf("Unmerged configuration files: %d\n", len(res.Files()))
			for _, f := range res.Files() {
				if !all {
					fmt.Printf("Latest changeset for %s (%d)\n", f, res.LatestFor(f).Version)
					changeset := res.LatestFor(f)
					evaluateDiff(f, changeset, res, interactive, true)
				} else {
					for _, c := range res[f] {
						evaluateDiff(f, c, res, interactive, false)
					}
				}
			}
		}}

	c.Flags().StringP("path", "p", "/etc", "Path to scan for unmanaged config files")
	c.Flags().BoolP("interactive", "i", true, "Interactive. Prompt to accept changes")
	c.Flags().BoolP("all", "a", false, "Review ALL found changes")

	return c
}
