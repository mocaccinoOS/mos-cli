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
	"strings"
)

func NewInitrdImage() *InitrdImage {
	return &InitrdImage{}
}

func NewInitrdImageFromFile(t *KernelType, file string) (*InitrdImage, error) {
	ans := NewInitrdImage()
	ans.Filename = file

	iprefix := t.InitrdPrefix
	if t.InitrdPrefix == "" {
		iprefix = "initramfs"
	}

	ans.Prefix = iprefix

	// Skip prefix + '-'
	file = file[len(iprefix)+1:]

	if t.Type != "" {
		ans.KernelType = t.Type
		file = file[len(t.Type)+1:]
	}

	words := strings.Split(file, "-")
	if t.WithArch {
		file = file[len(words[0])+1:]
		ans.Arch = words[0]
	}

	ans.Version = words[1]
	file = file[len(words[1])+1:]

	if t.Suffix != "" {
		ans.Suffix = file
	}

	return ans, nil
}

func (i *InitrdImage) SetPrefix(p string)     { i.Prefix = p }
func (i *InitrdImage) SetSuffix(s string)     { i.Suffix = s }
func (i *InitrdImage) SetKernelType(t string) { i.KernelType = t }
func (i *InitrdImage) SetArch(a string)       { i.Arch = a }
func (i *InitrdImage) SetVersion(v string)    { i.Version = v }
func (i *InitrdImage) SetFilename(f string)   { i.Filename = f }

func (i *InitrdImage) GetPrefix() string     { return i.Prefix }
func (i *InitrdImage) GetSuffix() string     { return i.Suffix }
func (i *InitrdImage) GetKernelType() string { return i.KernelType }
func (i *InitrdImage) GetArch() string       { return i.Arch }
func (i *InitrdImage) GetVersion() string    { return i.Version }
func (i *InitrdImage) GetFilename() string   { return i.Filename }

func (i *InitrdImage) String() string {
	data, _ := json.Marshal(i)
	return string(data)
}

func (i *InitrdImage) EqualTo(in *InitrdImage) bool {
	if i.Prefix != in.Prefix {
		return false
	}

	if i.KernelType != in.KernelType {
		return false
	}

	if i.Arch != in.Arch {
		return false
	}

	if i.Version != in.Version {
		return false
	}

	return true
}

func (i *InitrdImage) GenerateFilename() string {

	iprefix := i.Prefix
	if i.Prefix == "" {
		iprefix = "initramfs"
	}

	ans := iprefix + "-"

	if i.KernelType != "" {
		ans += i.KernelType
	}

	if i.Arch != "" {
		ans += "-" + i.Arch
	}

	if i.Version != "" {
		ans += "-" + i.Version
	}

	if i.Suffix != "" {
		ans += "-" + i.Suffix
	}

	return ans
}
