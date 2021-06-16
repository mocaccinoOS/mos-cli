// Copyright © 2021 Daniele Rondina, geaaru@sabayonlinux.org
//
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

package kernel

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	kernelspecs "github.com/MocaccinoOS/mos-cli/pkg/kernel/specs"
	//"path/filepath"
)

func ReadBootDir(bootdir string, supportedTypes []kernelspecs.KernelType) (*kernelspecs.BootFiles, error) {
	if bootdir == "" {
		bootdir = "/boot"
	}

	files, err := ioutil.ReadDir(bootdir)
	if err != nil {
		return nil, err
	}

	ans := kernelspecs.NewBootFiles(bootdir)

	for _, t := range supportedTypes {
		r := t.GetRegex()
		if r == nil {
			return nil, errors.New(
				fmt.Sprintf("Error on create regex for kernel type %s: %s",
					t.GetType(), err.Error()),
			)
		}
	}

	for _, file := range files {
		if file.IsDir() {
			// Ignoring directory
			continue
		}

		// Retrieve bzImage link
		if file.Name() == "bzImage" && (file.Mode()&os.ModeSymlink != 0) {
			linkedFile, err := os.Readlink(file.Name())
			if err == nil {
				ans.BzImageLink = linkedFile
			}
		}

		// Retrive Initrd link
		if file.Name() == "Initrd" && (file.Mode()&os.ModeSymlink != 0) {
			linkedFile, err := os.Readlink(file.Name())
			if err == nil {
				ans.InitrdLink = linkedFile
			}
		}

		for _, t := range supportedTypes {
			if t.GetRegex().MatchString(file.Name()) {

				isInirtd, err := t.IsInitrdFile(file.Name())
				if err != nil {
					return nil, errors.New(
						fmt.Sprintf("Error on check if the file %s is an initrd file: %s",
							file.Name(), err.Error(),
						))
				}

				if isInirtd {
					// Initrd image
					iimage, err := kernelspecs.NewInitrdImageFromFile(&t, file.Name())
					if err != nil {
						return nil, errors.New(
							fmt.Sprintf("Error on parse file %s: %s",
								file.Name(), err.Error(),
							))
					}

					err = ans.AddInitrdImage(iimage, &t)
					if err != nil {
						return nil, err
					}

				} else {
					// Kernel image
					kimage, err := kernelspecs.NewKernelImageFromFile(&t, file.Name())
					if err != nil {
						return nil, errors.New(
							fmt.Sprintf("Error on parse file %s: %s",
								file.Name(), err.Error(),
							))
					}

					err = ans.AddKernelImage(kimage, &t)
					if err != nil {
						return nil, err
					}
				}

				//fmt.Println("Read file ", file.Name())
				goto nextFile
			}
		}

	nextFile:
	}

	return ans, nil
}
