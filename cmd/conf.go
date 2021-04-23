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
	"github.com/MocaccinoOS/mos-cli/cmd/conf"
	"github.com/spf13/cobra"
)

// confCmd represents the conf command
var confCmd = &cobra.Command{
	Use:   "config-update",
	Short: "Automatically merge system-wide configuration files",
	Long: `After running upgrades, protected files might be created to avoid override custom user-defined
configuration files (for example, in /etc).

config-update allows you to check if there are configuration files that needs to be reviewed,
merges them automatically or interactively.

Compatible with etc-update and dispatch-conf`,
}

func init() {
	rootCmd.AddCommand(confCmd)

	confCmd.AddCommand(
		conf.NewCheckCommand(),
		conf.NewUpdateCommand(),
		conf.NewCleanCommand(),
	)
}
