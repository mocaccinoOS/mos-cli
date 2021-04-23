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
	"github.com/spf13/cobra"
)

func NewCheckCommand() *cobra.Command {
	c := &cobra.Command{Use: "check",
		Short: "Display a summary of available changes to review in the system",
		Long:  `Checks for unmerged configuration files in the system and displays a visual summary`,
		Run: func(cmd *cobra.Command, args []string) {

			path, _ := cmd.Flags().GetString("path")
			res := config.Scan(path)

			if len(res.Files()) == 0 {
				fmt.Println("All good!")
				return
			}

			fmt.Printf("Files with unmerged config files: %d\n", len(res.Files()))
			for _, f := range res.Files() {
				changefiles := res[f]
				fmt.Printf("- %s (%d unmerged config files)", f, len(changefiles))
			}
		}}

	c.Flags().StringP("path", "p", "/etc", "Path to scan for unmanaged config files")

	return c
}
