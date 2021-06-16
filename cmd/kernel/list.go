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
	"fmt"
	"os"

	"github.com/MocaccinoOS/mos-cli/pkg/kernel"
	kernelspecs "github.com/MocaccinoOS/mos-cli/pkg/kernel/specs"

	tablewriter "github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func NewListcommand() *cobra.Command {
	c := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List available kernels on system.",
		Long: `Shows kernels available in your system

$ mos kernel list

`,
		Run: func(cmd *cobra.Command, args []string) {

			jsonOutput, _ := cmd.Flags().GetBool("json")
			bootDir, _ := cmd.Flags().GetString("bootDir")

			// Temporary static configuration. I will move to
			// configuration file soon to permit more easy
			// customization and to support multiple kernel types.
			types := []kernelspecs.KernelType{
				kernelspecs.KernelType{
					Suffix:   "sabayon",
					Type:     "genkernel",
					WithArch: true,
				},
				kernelspecs.KernelType{
					Suffix:   "mocaccino",
					Type:     "vanilla",
					WithArch: true,
				},
			}

			bootFiles, err := kernel.ReadBootDir(bootDir, types)
			if err != nil {
				fmt.Println("Error on read boot directory: " + err.Error())
				os.Exit(1)
			}

			if jsonOutput {
				fmt.Println(bootFiles)
			} else {

				if len(bootFiles.Files) == 0 {
					fmt.Println("No kernel files available.")
					os.Exit(0)
				}

				table := tablewriter.NewWriter(os.Stdout)
				table.SetBorders(tablewriter.Border{
					Left: true, Top: false, Right: true, Bottom: false,
				})
				table.SetCenterSeparator("|")
				table.SetHeader([]string{
					"Type",
					"Suffix",
					"Version",
					"Has Initrd",
					"Has Kernel Image",
					"Has bzImage,Initrd links",
				})

				for _, kf := range bootFiles.Files {

					hasInitrd := false
					hasKernel := false
					hasLinks := false

					row := []string{
						kf.Type.GetType(),
						kf.Type.GetSuffix(),
					}

					version := ""
					if kf.Initrd != nil {
						version = kf.Initrd.GetVersion()
						hasInitrd = true
					}

					if kf.Kernel != nil {
						version = kf.Kernel.GetVersion()
						hasKernel = true
					}

					if hasKernel && hasInitrd &&
						bootFiles.BzImageLink != "" &&
						bootFiles.InitrdLink != "" &&
						kf.Kernel.GetFilename() == bootFiles.BzImageLink &&
						kf.Initrd.GetFilename() == bootFiles.InitrdLink {
						hasLinks = true
					}

					row = append(row, []string{
						version,
						fmt.Sprintf("%v", hasInitrd),
						fmt.Sprintf("%v", hasKernel),
						fmt.Sprintf("%v", hasLinks),
					}...)

					table.Append(row)
				}

				table.Render()
			}
		},
	}

	flags := c.Flags()
	flags.Bool("json", false, "JSON output")
	flags.String("bootdir", "/boot", "Directory where analyze kernel files.")

	return c
}
