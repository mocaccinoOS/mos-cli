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

package profile

import (
	"io/ioutil"
	"path"
	"regexp"

	kernelspecs "github.com/MocaccinoOS/mos-cli/pkg/kernel/specs"
)

func GetDefaultKernelProfiles() []kernelspecs.KernelType {
	return []kernelspecs.KernelType{
		kernelspecs.KernelType{
			Name:     "Sabayon",
			Suffix:   "sabayon",
			Type:     "genkernel",
			WithArch: true,
		},
		kernelspecs.KernelType{
			Name:     "Mocaccino",
			Suffix:   "mocaccino",
			Type:     "vanilla",
			WithArch: true,
		},
	}
}

func LoadKernelProfiles(dir string) ([]kernelspecs.KernelType, error) {
	ans := []kernelspecs.KernelType{}

	var regexRepo = regexp.MustCompile(`.yml$|.yaml$`)
	var err error

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return ans, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !regexRepo.MatchString(file.Name()) {
			continue
		}

		content, err := ioutil.ReadFile(path.Join(dir, file.Name()))
		if err != nil {
			// TODO: integrate warning logger
			continue
		}

		ktype, err := kernelspecs.KernelTypeFromYaml(content)
		if err != nil {
			continue
		}

		if ktype.GetType() != "" {
			ans = append(ans, *ktype)
		}
	}

	return ans, nil
}
