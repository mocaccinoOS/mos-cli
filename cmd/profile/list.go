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

package profile

// https://github.com/rancher-sandbox/cOS-toolkit/blob/master/packages/cos-features/cos-feature.sh
import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewListcommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
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
			handler := profileHandler(cmd)
			fmt.Println("Listing available and active system profiles")
			for _, p := range handler.List() {
				if p.Active {
					fmt.Println("- " + p.Name() + " (active)")
				} else {
					fmt.Println("- " + p.Name())
				}
			}
		},
	}

	profileHandlerFlags(c)
	return c
}
