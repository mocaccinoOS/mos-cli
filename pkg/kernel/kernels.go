package kernel

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/MocaccinoOS/mos-cli/pkg/utils"
)

type SearchResult struct {
	Packages []Package
}

type Package struct {
	Name, Category, Version string
}

func (p Package) Equal(pp Package) bool {
	if p.Name == pp.Name && p.Category == pp.Category && p.Version == pp.Version {
		return true
	}
	return false
}

func (p Package) EqualS(s string) bool {
	if s == fmt.Sprintf("%s/%s", p.Category, p.Name) {
		return true
	}
	return false
}

func (p Package) EqualNoV(pp Package) bool {
	if p.Name == pp.Name && p.Category == pp.Category {
		return true
	}
	return false
}

type Kernel struct {
	Type           string
	PackageName    string
	ModulesPackage string
}

func Installed() (searchResult SearchResult, err error) {
	var res []byte
	res, err = utils.RunSH("search", "luet search --installed kernel --output json")
	if err != nil {
		return
	}
	json.Unmarshal(res, &searchResult)
	searchResult = searchResult.FilterByCategory("kernel")
	return
}

func All() (searchResult SearchResult, err error) {
	var res []byte
	res, err = utils.RunSH("search", "luet search kernel --output json")
	if err != nil {
		return
	}
	json.Unmarshal(res, &searchResult)
	searchResult = searchResult.FilterByCategory("kernel")
	return
}

func (s SearchResult) FilterByCategory(cat string) SearchResult {
	new := SearchResult{Packages: []Package{}}

	for _, r := range s.Packages {
		if r.Category == cat {
			new.Packages = append(new.Packages, r)
		}
	}
	return new
}

func (s SearchResult) FilterByName(name string) SearchResult {
	new := SearchResult{Packages: []Package{}}

	for _, r := range s.Packages {
		if !strings.Contains(r.Name, name) {
			new.Packages = append(new.Packages, r)
		}
	}
	return new
}
