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

package kernelspecs

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func NewKernelFiles(t *KernelType) *KernelFiles {
	return &KernelFiles{
		Type: t,
	}
}

func NewBootFiles(dir string) *BootFiles {
	return &BootFiles{
		Dir:   dir,
		Files: []*KernelFiles{},
	}
}

func (b *BootFiles) GetDir() string { return b.Dir }

func (b *BootFiles) String() string {
	data, _ := json.Marshal(b)
	return string(data)
}

func (b *BootFiles) AddKernelImage(ki *KernelImage, t *KernelType) error {
	assigned := false

	// Check if exists a kernel image equal to the supplied file
	for _, k := range b.Files {
		if k.Kernel != nil && k.Kernel.EqualTo(ki) {
			return errors.New(fmt.Sprintf("File %s already present", ki.GetFilename()))
		}
	}

	// Check if exists a initrd image related to the input file
	for idx, k := range b.Files {
		if k.Kernel == nil && k.Initrd != nil &&
			k.Initrd.GetSuffix() == ki.GetSuffix() &&
			k.Initrd.GetKernelType() == ki.GetType() &&
			k.Initrd.GetVersion() == ki.GetVersion() &&
			k.Initrd.GetArch() == ki.GetArch() {

			b.Files[idx].Kernel = ki
			assigned = true
			break
		}
	}

	if !assigned {
		kf := NewKernelFiles(t)
		kf.Kernel = ki
		b.Files = append(b.Files, kf)
	}

	return nil
}

func (b *BootFiles) BzImageLinkExistingKernel() bool {
	ans := false

	if b.BzImageLink != "" {
		for _, f := range b.Files {
			if f.Kernel != nil && f.Kernel.GetFilename() == b.BzImageLink {
				ans = true
				break
			}
		}
	}

	return ans
}

func (b *BootFiles) RetrieveBzImageSelectedKernel() *KernelFiles {
	var kernelImage *KernelImage = nil

	if b.BzImageLink == "" {
		return nil
	}

	for _, f := range b.Files {

		if f.Kernel == nil {
			continue
		}

		kPrefix := f.Kernel.GetPrefix()

		if f.Kernel.GetType() != "" {
			kPrefix += "-" + f.Kernel.GetType()
		}

		if strings.HasPrefix(b.BzImageLink, kPrefix) {

			if f.Kernel.GetSuffix() != "" && !strings.HasSuffix(b.BzImageLink, f.Kernel.GetSuffix()) {
				continue
			}

			if f.Type.GetRegex().MatchString(b.BzImageLink) {
				// POST: matched kernel
				kimage, err := NewKernelImageFromFile(f.Type, b.BzImageLink)
				if err == nil && kimage != nil {
					kernelImage = kimage
					break
				}
			}

		}

	}

	if kernelImage != nil {
		for _, f := range b.Files {
			if f.Kernel != nil && f.Kernel.GetPrefix() == kernelImage.GetPrefix() &&
				f.Kernel.GetArch() == kernelImage.GetArch() &&
				f.Kernel.GetSuffix() == kernelImage.GetSuffix() &&
				f.Kernel.GetType() == kernelImage.GetType() {
				return f
			}
		}
	}

	return nil
}

func (b *BootFiles) GetFile(version, ktype string) (*KernelFiles, error) {
	var ans *KernelFiles = nil

	if version == "" {
		return nil, errors.New("Invalid version")
	}

	for idx, f := range b.Files {
		if f.Kernel != nil && f.Kernel.GetVersion() == version {

			if ktype != "" && f.Kernel.GetType() != ktype {
				continue
			}

			ans = b.Files[idx]
		}
	}

	if ans == nil {
		return ans, errors.New("No kernel found with the option selected.")
	}

	return ans, nil
}

func (b *BootFiles) AddInitrdImage(i *InitrdImage, t *KernelType) error {
	assigned := false

	// Check if exists a initrd image equal to the supplied file
	for _, f := range b.Files {
		if f.Initrd != nil && f.Initrd.EqualTo(i) {
			return errors.New(fmt.Sprintf("File %s already present", i.GetFilename()))
		}
	}

	// Check if exists a kernel image related to the input file
	for idx, k := range b.Files {
		if k.Initrd == nil && k.Kernel != nil &&
			k.Kernel.GetSuffix() == i.GetSuffix() &&
			k.Kernel.GetType() == i.GetKernelType() &&
			k.Kernel.GetVersion() == i.GetVersion() &&
			k.Kernel.GetArch() == i.GetArch() {

			b.Files[idx].Initrd = i
			assigned = true
			break
		}
	}

	if !assigned {
		kf := NewKernelFiles(t)
		kf.Initrd = i
		b.Files = append(b.Files, kf)
	}

	return nil
}
