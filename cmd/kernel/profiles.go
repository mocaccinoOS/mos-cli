// Copyright © 2021 Daniele Rondina <geaaru@sabayonlinux.org>
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, see <http://www.gnu.org/licenses/>.

package cmdkernel

import (
	"encoding/json"
	"fmt"
	"os"

	kernelspecs "github.com/MocaccinoOS/mos-cli/pkg/kernel/specs"
	"github.com/MocaccinoOS/mos-cli/pkg/profile"

	tablewriter "github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func NewProfilesCommand() *cobra.Command {
	c := &cobra.Command{
		Use:     "profiles",
		Aliases: []string{"p"},
		Short:   "List available kernels profiles.",
		Long: `Shows kernels available in your system

$ mos kernel profiles

`,
		Run: func(cmd *cobra.Command, args []string) {

			jsonOutput, _ := cmd.Flags().GetBool("json")
			kernelProfilesDir, _ := cmd.Flags().GetString("kernel-profiles-dir")

			types := []kernelspecs.KernelType{}
			if kernelProfilesDir != "" {
				types, _ = profile.LoadKernelProfiles(kernelProfilesDir)
			} else {
				types = profile.GetDefaultKernelProfiles()
			}

			if jsonOutput {
				data, err := json.Marshal(types)
				if err != nil {
					fmt.Println(fmt.Errorf("Error on convert data to json: %s", err.Error()))
					os.Exit(1)
				}
				fmt.Println(string(data))

			} else {

				if len(types) == 0 {
					fmt.Println("No kernel profiles availables. I will use default profiles.")
					os.Exit(0)
				}

				table := tablewriter.NewWriter(os.Stdout)
				table.SetBorders(tablewriter.Border{
					Left: true, Top: false, Right: true, Bottom: false,
				})
				table.SetCenterSeparator("|")
				table.SetHeader([]string{
					"Name",
					"Kernel Prefix",
					"Initrd Prefix",
					"Suffix",
					"Type",
					"With Arch",
				})

				for _, kt := range types {

					table.Append([]string{
						kt.GetName(),
						kt.GetKernelPrefixSanitized(),
						kt.GetInitrdPrefixSanitized(),
						kt.GetSuffix(),
						kt.GetType(),
						fmt.Sprintf("%v", kt.WithArch),
					})
				}

				table.Render()
			}
		},
	}

	flags := c.Flags()
	flags.Bool("json", false, "JSON output")
	flags.String("kernel-profiles-dir", "/etc/mocaccino/kernels-profiles/",
		"Specify the directory where read the kernel types profiles supported.")

	return c
}
