package common

import (
	"regexp"
	"strings"
)

// Cmd ...
func Cmd(input string) []string {

	var texts []string
	if strings.HasPrefix(input, "/") {
		// 字串 - 執行指令
		cmdLine := strings.TrimLeft(input, "/")

		re := regexp.MustCompile(`(^\w+)\s*(\w*)`)
		matches := re.FindStringSubmatch(cmdLine)

		if len(matches) > 1 {

		}
		return texts
	}
	return texts
}
