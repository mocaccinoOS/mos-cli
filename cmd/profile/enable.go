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
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
)

func NewLoadcommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "enable",
		Short: "Enables a MocaccinoOS system profile",
		Long: `Turns on a system profile by calling

$ mos profile enable profile1 profile2 profile3...

Profiles can be listed with:

$ mos profile list
`,
		Run: func(cmd *cobra.Command, args []string) {
			handler := profileHandler(cmd)
			var reserr error
			for _, a := range args {
				profile, err := handler.Search(a)
				if err != nil {
					reserr = multierror.Append(reserr, err)
					continue
				}
				fmt.Println("Activating", profile.Name())
				if err := handler.Activate(profile); err != nil {
					fmt.Println("Failed activating", profile.Name(), err)
					reserr = multierror.Append(reserr, err)
				}
			}
			if reserr != nil {
				fmt.Println("Failed applying profile:", reserr)
				os.Exit(1)
			} else {
				fmt.Println("Profiles loaded correctly")
			}
		},
	}

	profileHandlerFlags(c)
	return c
}
