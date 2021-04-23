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

package utils

import (
	"fmt"
	"strings"
)

func Ask(question string) bool {
	var input string

	fmt.Printf("%s? [y/N]: \n", question)
	_, err := fmt.Scanln(&input)
	if err != nil {
		return false
	}
	input = strings.ToLower(input)

	if input == "y" || input == "yes" {
		return true
	}
	return false
}

const (
	AcceptQuestion = iota
	IgnoreQuestion
	DiscardQuestion
)

func Accept(question string) int {
	var input string

	fmt.Printf("%s? [Accept/Discard/Ignore]: \n", question)
	_, err := fmt.Scanln(&input)
	if err != nil {
		return IgnoreQuestion
	}
	input = strings.ToLower(input)

	if input == "a" || input == "accept" {
		return AcceptQuestion
	}

	if input == "d" || input == "discard" {
		return DiscardQuestion
	}

	return IgnoreQuestion
}
