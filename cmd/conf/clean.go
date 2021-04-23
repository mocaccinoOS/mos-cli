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
	"github.com/spf13/cobra"
)

func NewCleanCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "clean",
		Short: "Cleans up unmerged config files",
		Long: `Removes unmerged config files from the system. To run without prompts, run with --interactive=false.
		
For example:

$ mos clean --interactive=false

Cleans up all the unmerged config diff files in /etc.
`,
		Run: func(cmd *cobra.Command, args []string) {
			path, _ := cmd.Flags().GetString("path")
			interactive, _ := cmd.Flags().GetBool("interactive")
			res := config.Scan(path)

			fmt.Printf("Unmerged configuration files: %d", len(res.Files()))
			for _, f := range res.Files() {
				changeset := res[f]
				fmt.Printf("Found %d changes for %s", len(changeset), f)
				if interactive {
					if utils.Ask("Do you want to clean config merges for " + f) {
						checkErr(res.CleanChanges(f))
					}
				} else {
					checkErr(res.CleanChanges(f))
				}
			}
		}}

	c.Flags().StringP("path", "p", "/etc", "Path to scan for unmanaged config files")
	c.Flags().BoolP("interactive", "i", true, "Interactive. Prompt to accept changes")

	return c
}
