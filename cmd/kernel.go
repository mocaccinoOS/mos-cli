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

package cmd

import (
	cmdkernel "github.com/MocaccinoOS/mos-cli/cmd/kernel"
	"github.com/spf13/cobra"
)

var kernelCmd = &cobra.Command{
	Use:   "kernel",
	Short: "Manage system kernels and initrd (experimental).",
	Long:  `Manage kernels and initrd images of your system.`,
}

func init() {
	rootCmd.AddCommand(kernelCmd)

	kernelCmd.AddCommand(
		cmdkernel.NewListcommand(),
		cmdkernel.NewGeninitrdCommand(),
		cmdkernel.NewProfilesCommand(),
	)
}
