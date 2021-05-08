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

	"github.com/MocaccinoOS/mos-cli/pkg/kernel"
	"github.com/spf13/cobra"
)

func NewListcommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "List available MocaccinoOS kernels",
		Long: `Shows kernels available in mocaccino os repositories

$ mos kernel-switcher list

To switch to one of them

$ mos kernel-switcher switch kernel/mocaccino-full
`,
		Run: func(cmd *cobra.Command, args []string) {

			allKernelsPackages, err := kernel.All()
			if err != nil {
				log.Fatal(err)
			}
			installed, err := kernel.Installed()
			if err != nil {
				log.Fatal(err)
			}
			kernels := allKernelsPackages.FilterByName("firmware").FilterByName("initramfs").FilterByName("minimal").FilterByName("modules").FilterByName("sources")

			for _, j := range kernels.Packages {
				install := ""
				for _, k := range installed.Packages {
					if k.EqualNoV(j) {
						install = "installed"
					}
				}
				fmt.Println(fmt.Sprintf("%s/%s (%s) %s", j.Category, j.Name, j.Version, install))
			}
		},
	}

	return c
}
