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
	"os/exec"
	"path/filepath"
	"strings"

	kernelspecs "github.com/MocaccinoOS/mos-cli/pkg/kernel/specs"
	. "github.com/MocaccinoOS/mos-cli/pkg/logger"
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

		DebugC("Analyzing file", file.Name(), "...")

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

				DebugC("File", file.Name(), "match type", t.GetName())

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

func GrubMkconfig(grubCfgFile string, dryRun bool) error {
	if grubCfgFile == "" {
		return errors.New("Invalid grub config file path")
	}

	grubDir := filepath.Dir(grubCfgFile)
	_, err := os.Stat(grubDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(grubDir, 0755)
			if err != nil {
				return errors.New(fmt.Sprintf(
					"Error on create directory %s: %s",
					grubDir, err.Error()))
			}
		} else {
			return errors.New(fmt.Sprintf(
				"Error on stat directory %s: %s",
				grubDir, err.Error()))
		}

	}

	//grub-mkconfig -o ${MOCACCINO_TARGET}/boot/grub/grub.cfg
	grubBinary := "grub-mkconfig"
	args := []string{
		"-o", grubCfgFile,
	}

	if dryRun {
		fmt.Println("[dry-run mode] command: " + grubBinary + " " + strings.Join(args, " "))
		return nil
	}

	fmt.Println(fmt.Sprintf("Creating grub config file %s...", grubCfgFile))

	grubCommand := exec.Command(grubBinary, args...)
	grubCommand.Stdout = os.Stdout
	grubCommand.Stderr = os.Stderr

	err = grubCommand.Start()
	if err != nil {
		return errors.New(
			fmt.Sprintf("Error on start %s command: %s", grubBinary, err.Error()))
	}

	err = grubCommand.Wait()
	if err != nil {
		return errors.New(
			fmt.Sprintf("Error on waiting %s command: %s", grubBinary, err.Error()))
	}

	if grubCommand.ProcessState.ExitCode() != 0 {
		return errors.New(
			fmt.Sprintf("%s command exiting with %d",
				grubBinary, grubCommand.ProcessState.ExitCode()))
	}

	return nil
}
