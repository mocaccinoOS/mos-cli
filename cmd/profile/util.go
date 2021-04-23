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

import (
	profile "github.com/MocaccinoOS/mos-cli/pkg/profile"
	"github.com/spf13/cobra"
)

func profileHandler(cmd *cobra.Command) profile.ProfileHandler {
	activeDirectory, _ := cmd.Flags().GetString("active-directory")
	profileDirectory, _ := cmd.Flags().GetString("profile-directory")

	return profile.ProfileHandler{ActiveDirectory: activeDirectory, ProfileDirectory: profileDirectory}
}

func profileHandlerFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("active-directory", "a", "/etc/mocaccino/profiles/active", "Path to active directory where to symlink")
	cmd.Flags().StringP("profile-directory", "p", "/etc/mocaccino/profiles/available", "Path to available profiles")
}
