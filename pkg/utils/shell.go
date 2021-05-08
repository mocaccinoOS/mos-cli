package utils

import (
	"os/exec"
	"strings"
)

func RunSH(stepName, bashFragment string) ([]byte, error) {
	cmd := exec.Command("sh", "-s")
	cmd.Stdin = strings.NewReader(bashWrap(bashFragment))
	//	log.Printf("Running in background: %v", stepName)

	return cmd.CombinedOutput()
}

func bashWrap(cmd string) string {
	return `
set -o errexit
set -o nounset
set -o pipefail

` + cmd + `
`
}
