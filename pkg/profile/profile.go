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
	"errors"
	"os"
	"path/filepath"

	"github.com/MocaccinoOS/mos-cli/pkg/utils"
)

type Profile struct {
	Path   string
	Active bool
}

func (p Profile) Name() string {
	path := filepath.Base(p.Path)
	var extension = filepath.Ext(path)
	var profileName = path[0 : len(path)-len(extension)]
	return profileName
}

type ProfileHandler struct {
	ActiveDirectory  string
	ProfileDirectory string
}

func (ph ProfileHandler) List() []Profile {
	files, err := utils.ListDir(ph.ProfileDirectory)
	if err != nil {
		return []Profile{}
	}
	profiles := []Profile{}
	for _, f := range files {
		if ph.isProfileActive(f) {
			profiles = append(profiles, Profile{Path: f, Active: true})
		} else {
			profiles = append(profiles, Profile{Path: f, Active: false})
		}
	}

	return profiles
}

func (ph ProfileHandler) isProfileActive(s string) bool {
	if _, err := os.Lstat(filepath.Join(ph.ActiveDirectory, filepath.Base(s))); err == nil {
		return true
	}
	return false
}

func (ph ProfileHandler) Activate(p Profile) error {
	if ph.isProfileActive(p.Path) {
		return nil
	}

	if _, err := os.Stat(ph.ActiveDirectory); err != nil {
		os.MkdirAll(ph.ActiveDirectory, 600)
	}
	return os.Symlink(filepath.Join(ph.ProfileDirectory, filepath.Base(p.Path)), filepath.Join(ph.ActiveDirectory, filepath.Base(p.Path)))
}

func (ph ProfileHandler) Deactivate(p Profile) error {
	if !ph.isProfileActive(p.Path) {
		return nil
	}
	return os.Remove(filepath.Join(ph.ActiveDirectory, filepath.Base(p.Path)))
}

func (ph ProfileHandler) Search(p string) (Profile, error) {
	for _, pp := range ph.List() {

		if pp.Name() == p {
			return pp, nil
		}
	}

	return Profile{}, errors.New("profile not found")
}
