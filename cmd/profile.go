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

package cmd

import (
	"github.com/MocaccinoOS/mos-cli/cmd/profile"
	"github.com/spf13/cobra"
)

// profileCmd represents the profile command
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Handle MocaccinoOS system profiles",
	Long:  `Disable, Enable and list available system profiles`,
}

func init() {
	rootCmd.AddCommand(profileCmd)

	profileCmd.AddCommand(
		profile.NewLoadcommand(),
		profile.NewListcommand(),
		profile.NewUnloadcommand(),
	)
}
