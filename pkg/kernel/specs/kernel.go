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

func NewKernelImage() *KernelImage {
	return &KernelImage{}
}

func NewKernelImageFromFile(t *KernelType, file string) (*KernelImage, error) {
	ans := NewKernelImage()
	ans.Filename = file

	kprefix := t.KernelPrefix
	if t.KernelPrefix == "" {
		kprefix = "kernel"
	}

	ans.Prefix = kprefix

	// Skip prefix + '-'
	file = file[len(kprefix)+1:]

	if t.Type != "" {
		ans.Type = t.Type
		file = file[len(t.Type)+1:]
	}

	words := strings.Split(file, "-")
	if t.WithArch {
		file = file[len(words[0])+1:]
		ans.Arch = words[0]
	}

	ans.Version = words[1]

	if t.Suffix != "" {
		file = file[len(words[1])+1:]
		ans.Suffix = file
	}

	return ans, nil
}

func (k *KernelImage) SetPrefix(p string)   { k.Prefix = p }
func (k *KernelImage) SetSuffix(s string)   { k.Suffix = s }
func (k *KernelImage) SetVersion(v string)  { k.Version = v }
func (k *KernelImage) SetArch(v string)     { k.Arch = v }
func (k *KernelImage) SetType(t string)     { k.Type = t }
func (k *KernelImage) SetFilename(f string) { k.Filename = f }

func (k *KernelImage) GetPrefix() string   { return k.Prefix }
func (k *KernelImage) GetSuffix() string   { return k.Suffix }
func (k *KernelImage) GetVersion() string  { return k.Version }
func (k *KernelImage) GetArch() string     { return k.Arch }
func (k *KernelImage) GetType() string     { return k.Type }
func (k *KernelImage) GetFilename() string { return k.Filename }

func (k *KernelImage) String() string {
	data, _ := json.Marshal(k)
	return string(data)
}

func (k *KernelImage) EqualTo(i *KernelImage) bool {
	if k.Type != i.Type {
		return false
	}

	if k.Prefix != i.Prefix {
		return false
	}

	if k.Suffix != i.Suffix {
		return false
	}

	if k.Version != i.Version {
		return false
	}

	if k.Arch != i.Arch {
		return false
	}

	return true
}

func (k *KernelImage) GenerateFilename() string {

	kprefix := k.Prefix
	if k.Prefix == "" {
		kprefix = "kernel"
	}

	ans := kprefix + "-"

	if k.Type != "" {
		ans += k.Type
	}

	if k.Arch != "" {
		ans += "-" + k.Arch
	}

	if k.Version != "" {
		ans += "-" + k.Version
	}

	if k.Suffix != "" {
		ans += "-" + k.Suffix
	}

	return ans
}
