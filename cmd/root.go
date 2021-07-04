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
	"fmt"
	"os"
	"strings"

	. "github.com/MocaccinoOS/mos-cli/pkg/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// Build time and commit information.
//
// ⚠️ WARNING: should only be set by "-ldflags".
var (
	Version     string
	BuildTime   string
	BuildCommit string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "mos",
	Short:   "A CLI to handle MocaccinoOS system settings",
	Version: fmt.Sprintf("%s-g%s %s", Version, BuildCommit, BuildTime),
	Long: `mos allows you to enable/disable/list system profiles, and manage configuration files.

To list all the available system profiles:

$ mos profile list

To enable one of them:

$ mos profile enable profile1

To disable: 

$ mos profile disable profile1


`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig()

		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			viper.Set("debug", debug)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mos.yaml)")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug logging.")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	viper.SetEnvPrefix("LUET")
	viper.AutomaticEnv() // read in environment variables that match
	// Create EnvKey Replacer for handle complex structure
	replacer := strings.NewReplacer(".", "__")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetTypeByDefaultValue(true)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	InitAurora()
	ZapLogger()
}
