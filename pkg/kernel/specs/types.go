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
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

func (t *KernelType) SetKernelPrefix(p string) { t.KernelPrefix = p }
func (t *KernelType) SetInitrdPrefix(p string) { t.InitrdPrefix = p }
func (t *KernelType) SetSuffix(s string)       { t.Suffix = s }
func (t *KernelType) SetType(s string)         { t.Type = s }

func (t *KernelType) GetKernelPrefix() string { return t.KernelPrefix }
func (t *KernelType) GetInitrdPrefix() string { return t.InitrdPrefix }
func (t *KernelType) GetSuffix() string       { return t.Suffix }
func (t *KernelType) GetType() string         { return t.Type }
func (t *KernelType) GetName() string         { return t.Name }

func (t *KernelType) GetInitrdPrefixSanitized() string {
	initrdprefix := t.InitrdPrefix
	if t.InitrdPrefix == "" {
		initrdprefix = "initramfs"
	}

	return initrdprefix
}

func (t *KernelType) GetKernelPrefixSanitized() string {
	kprefix := t.KernelPrefix
	if t.KernelPrefix == "" {
		kprefix = "kernel"
	}

	return kprefix
}

func (t *KernelType) IsInitrdFile(f string) (bool, error) {
	ans := false
	if f == "" {
		return ans, errors.New("Invalid file path")
	}

	initrdprefix := t.GetInitrdPrefixSanitized()

	if strings.HasPrefix(f, initrdprefix) {
		ans = true
	}

	return ans, nil
}

func (t *KernelType) IsKernelFile(f string) (bool, error) {
	ans := false
	if f == "" {
		return ans, errors.New("Invalid kernel file path")
	}

	kprefix := t.GetKernelPrefixSanitized()

	if strings.HasPrefix(f, kprefix) {
		ans = true
	}

	return ans, nil
}

func (t *KernelType) getKernelRegex() string {
	kprefix := t.GetKernelPrefixSanitized()

	ans := fmt.Sprintf("^%s-", kprefix)

	if t.Type != "" {
		ans += fmt.Sprintf("%s-.*", t.Type)
	} else {
		ans += ".*"
	}

	if t.Suffix != "" {
		ans += fmt.Sprintf("-%s$", t.Suffix)
	}

	return ans
}

func (t *KernelType) getInitrdRegex() string {
	initrdprefix := t.GetInitrdPrefixSanitized()

	ans := fmt.Sprintf("^%s-", initrdprefix)

	if t.Type != "" {
		ans += fmt.Sprintf("%s-.*", t.Type)
	} else {
		ans += ".*"
	}

	if t.Suffix != "" {
		ans += fmt.Sprintf("-%s$", t.Suffix)
	}

	return ans
}

func (t *KernelType) GetRegex() *regexp.Regexp {
	if t.Regex == nil {
		regstrk := t.getKernelRegex()
		regstri := t.getInitrdRegex()

		t.Regex = regexp.MustCompile(
			fmt.Sprintf("%s|%s", regstrk, regstri),
		)
	}

	return t.Regex
}

func KernelTypeFromYaml(data []byte) (*KernelType, error) {
	ans := &KernelType{}
	if err := yaml.Unmarshal(data, ans); err != nil {
		return nil, err
	}

	return ans, nil
}
