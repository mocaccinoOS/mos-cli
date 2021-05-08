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

package kernelswitcher

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/MocaccinoOS/mos-cli/pkg/kernel"
	"github.com/spf13/cobra"
)

func NewSwitchcommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "switch",
		Short: "List available MocaccinoOS system profiles",
		Long: `Shows currently installed profiles in the system

$ mos profile list

Profiles can be installed with luet, to show all the available profiles, run:

$ luet search system-profile

To install one of them:

$ luet install system-profile/default-systemd

To enable:

$ mos profile enable default-systemd`,
		Run: func(cmd *cobra.Command, args []string) {

			argkernel := args[0]

			fmt.Println("Switching to kernel", argkernel)

			allKernelsPackages, err := kernel.All()
			if err != nil {
				log.Fatal(err)
			}
			installed, err := kernel.Installed()
			if err != nil {
				log.Fatal(err)
			}

			for _, i := range installed.Packages {
				if i.EqualS(argkernel) {
					log.Fatal("Kernel already installed")
				}
			}

			found := false
			for _, i := range allKernelsPackages.Packages {
				if i.EqualS(argkernel) {
					found = true
				}
			}
			if !found {
				log.Fatal("Provided kernel not found")
			}

			binary, lookErr := exec.LookPath("luet")
			if lookErr != nil {
				panic(lookErr)
			}

			cmdargs := []string{"luet", "replace", "--nodeps", "--for", argkernel, "--for", fmt.Sprintf("%smodules", strings.ReplaceAll(argkernel, "full", ""))}
			for _, i := range installed.Packages {
				cmdargs = append(cmdargs, fmt.Sprintf("%s/%s", i.Category, i.Name))
			}

			env := os.Environ()

			execErr := syscall.Exec(binary, cmdargs, env)
			if execErr != nil {
				panic(execErr)
			}
		},
	}

	return c
}
