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

package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// Configs is a map of files -> and related config change
type Configs map[string][]ConfigChange

type ConfigChange struct {
	Path    string
	Version int
}

func (c ConfigChange) Diff(frompath string) ([]diffmatchpatch.Diff, error) {
	dat, err := ioutil.ReadFile(frompath)
	if err != nil {
		return nil, errors.Wrapf(err, "while reading file '%s'", frompath)
	}

	dmp := diffmatchpatch.New()

	str, err := c.Content()
	if err != nil {
		return nil, errors.Wrap(err, "while reading diff content")
	}

	diffs := dmp.DiffMain(string(dat), str, false)

	//fmt.Println(dmp.DiffPrettyText(diffs))
	return diffs, nil
}

func (c ConfigChange) Content() (string, error) {
	dat, err := ioutil.ReadFile(c.Path)
	if err != nil {
		return "", errors.Wrap(err, "while reading file")
	}
	return string(dat), nil
}

func (c ConfigChange) Remove() error {
	return os.RemoveAll(c.Path)
}

func (c ConfigChange) Merge(s string) (Merge, error) {
	diffs, err := c.Diff(s)
	if err != nil {
		return Merge{}, errors.Wrap(err, "while computing diffs")
	}
	dmp := diffmatchpatch.New()

	datByte, err := ioutil.ReadFile(s)
	if err != nil {
		return Merge{}, errors.Wrap(err, "while reading file")
	}

	patch := dmp.PatchMake(diffs)

	res, applies := dmp.PatchApply(patch, string(datByte))

	return Merge{Content: res, Applies: applies, Path: s}, nil
}

func (c Configs) LatestFor(s string) ConfigChange {
	latestIteration := &c[s][0]

	for _, change := range c[s] {
		if change.Version > latestIteration.Version {
			latestIteration = &change
		}
	}
	return *latestIteration
}

type Merge struct {
	Content string
	Applies []bool
	Path    string
}

func (m Merge) Apply() error {
	info, err := os.Stat(m.Path)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(m.Path, []byte(m.Content), info.Mode().Perm())
}

func (c Configs) Merge(s string) (Merge, error) {
	return c.LatestFor(s).Merge(s)
}

func (c Configs) CleanChanges(s string) error {
	changes, ok := c[s]
	if !ok {
		return errors.New("changes not found")
	}
	var err error
	for _, v := range changes {
		err = multierror.Append(v.Remove())
	}
	return err
}

func (c Configs) Files() []string {
	res := []string{}
	for f := range c {
		res = append(res, f)
	}
	return res
}

func getDiffFileData(s string) map[string]string {
	return getParams(`\.\_cfg(?P<Number>\d+)\_(?P<File>\S+).*`, s)
}

func getParams(regEx, url string) (paramsMap map[string]string) {

	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(url)

	paramsMap = make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}

// Scan does:
// 1: walk etc
// 2: find filenames with pattern https://github.com/mudler/luet/blob/a83be204e8bd3498a736fa9fe401fb00bb91a893/pkg/compiler/artifact.go#L457
//    fmt.Sprintf("._cfg%04d_%s", i, filepath.Base(path))))
// 3: get latest - compare with current with diff. If no diff automerge.
//    offer option to automerge
//    failback to ask. discard/accept change
//fmt.Println(res)
func Scan(path string) Configs {
	res := Configs{}
	err := filepath.Walk(path, func(currentpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if strings.HasPrefix(info.Name(), "._cfg") {
			data := getDiffFileData(info.Name())
			nr, _ := strconv.Atoi(data["Number"])
			realLocation := filepath.Join(filepath.Dir(currentpath), data["File"])
			_, ok := res[realLocation]
			if !ok {
				res[realLocation] = []ConfigChange{{Version: nr, Path: currentpath}}
			} else {
				res[realLocation] = append(res[realLocation], ConfigChange{Version: nr, Path: currentpath})
			}
		}
		return nil
	})
	if err != nil {
		return res
	}
	return res
}
